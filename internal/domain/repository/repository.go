package repository

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/dbunt1tled/url-shortener/storage/mysql"
	"github.com/pkg/errors"
)

type Filter struct {
	Field    string
	Operator string
	Value    any
}

type Sort struct {
	Field     string
	Direction string
}

type BaseRepository[T any] struct {
	db        *mysql.Mysql
	tableName string
}

func NewBaseRepository[T any](db *mysql.Mysql, tableName string) *BaseRepository[T] {
	return &BaseRepository[T]{
		db:        db,
		tableName: tableName,
	}
}

func (r *BaseRepository[T]) Create(ctx context.Context, entity any) (*T, error) {
	var (
		cols []string
		vals []any
		id   int64
	)

	switch v := entity.(type) {
	case map[string]any:
		for k, val := range v {
			cols = append(cols, k)
			vals = append(vals, val)
		}
	default:
		t := reflect.TypeOf(v)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		val := reflect.ValueOf(v)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tag := field.Tag.Get("db")
			if tag == "" {
				continue
			}
			cols = append(cols, tag)
			vals = append(vals, val.Field(i).Interface())
		}
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		r.tableName,
		strings.Join(cols, ", "),
		strings.Repeat("?, ", len(cols)-1)+"?",
	)

	result, err := r.db.GetDB().ExecContext(ctx, query, vals...)
	if err != nil {
		return nil, errors.Wrap(err, "repository create error")
	}

	id, err = result.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "repository last insert id error")
	}
	return r.FindByID(ctx, id)

}

func (r *BaseRepository[T]) FindByID(ctx context.Context, id int64) (*T, error) {
	var entity T
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ? LIMIT 1", r.tableName)
	err := r.db.GetDB().GetContext(ctx, &entity, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("find by id failed: %w", err)
	}
	return &entity, nil
}

func (r *BaseRepository[T]) List(ctx context.Context, filters []Filter, sort *Sort, limit, offset int) ([]*T, error) {
	var (
		whereParts []string
		args       []any
	)

	for _, f := range filters {
		whereParts = append(whereParts, fmt.Sprintf("%s %s ?", f.Field, f.Operator))
		args = append(args, f.Value)
	}

	query := fmt.Sprintf("SELECT * FROM %s", r.tableName)
	if len(whereParts) > 0 {
		query += " WHERE " + strings.Join(whereParts, " AND ")
	}

	if sort != nil && sort.Field != "" {
		dir := strings.ToUpper(sort.Direction)
		if dir != "ASC" && dir != "DESC" {
			dir = "ASC"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", sort.Field, dir)
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	var results []*T
	if err := r.db.GetDB().SelectContext(ctx, &results, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*T{}, nil
		}
		return nil, errors.Wrap(err, "repository list error")
	}

	return results, nil
}

func (r *BaseRepository[T]) Update(ctx context.Context, id int64, data map[string]any) (*T, error) {
	if len(data) == 0 {
		return r.FindByID(ctx, id)
	}

	setParts := make([]string, 0, len(data))
	args := make([]any, 0, len(data)+1)

	for col, val := range data {
		setParts = append(setParts, fmt.Sprintf("%s = ?", col))
		args = append(args, val)
	}

	args = append(args, id)

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = ?",
		r.tableName,
		strings.Join(setParts, ", "),
	)

	_, err := r.db.GetDB().ExecContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "repository update failed")
	}

	return r.FindByID(ctx, id)
}

func (r *BaseRepository[T]) One(ctx context.Context, filters []Filter, sort *Sort) (*T, error) {
	result, err := r.List(ctx, filters, sort, 1, 0)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}
	return result[0], nil
}

func (r *BaseRepository[T]) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", r.tableName)
	result, err := r.db.GetDB().ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
