package mysql

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"go_first/internal/lib/common/model"
	"go_first/internal/lib/common/model/user"
	argon "go_first/internal/lib/common/password"
	"time"
)

type UserFilter struct {
	ID          int64
	Email       string
	PhoneNumber string
	Status      int
}

func (db *Mysql) CreateUser(
	firstName string,
	secondName string,
	email string,
	phoneNumber string,
	password string,
	status int,
) (*user.User, error) {
	dateTime := time.Now()
	smt, err := db.db.Prepare(`
		INSERT INTO users (
			first_name,
			second_name,
			email,
			phone_number,
			password,
			status,
		    updated_at,
		    created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return nil, errors.Wrap(err, "user prepare error")
	}
	hashPassword, err := argon.Hash(password)
	if err != nil {
		return nil, errors.Wrap(err, "user hash password error")
	}
	res, err := smt.Exec(
		firstName,
		secondName,
		email,
		phoneNumber,
		hashPassword,
		status,
		dateTime.Format(time.DateTime),
		dateTime.Format(time.DateTime),
	)
	if err != nil {
		var dbErr *mysql.MySQLError
		if errors.As(err, &dbErr) {
			if dbErr.Number == 1062 {
				return nil, errors.New("user already exists")
			}
		}
		return nil, errors.Wrap(err, "create user error")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "last insert user id error")
	}
	return &user.User{
		ID:          id,
		FirstName:   firstName,
		SecondName:  secondName,
		Email:       email,
		PhoneNumber: phoneNumber,
		Password:    hashPassword,
		Status:      status,
		UpdatedAt:   model.Timestamp{Time: dateTime},
		CreatedAt:   model.Timestamp{Time: dateTime},
	}, nil
}

func (db *Mysql) GetUser(filter UserFilter) (*user.User, error) {
	sqlQuery := ""
	if filter.PhoneNumber != "" {
		sqlQuery += " phone_number = '" + filter.PhoneNumber + "' OR "
	}

	if filter.Email != "" {
		sqlQuery += " AND email = '" + filter.Email + "' "
	}

	if filter.Status > 0 {
		sqlQuery += " AND status = " + fmt.Sprint(filter.Status) + " "
	}

	sqlQuery = fmt.Sprintf("SELECT * FROM users WHERE 1=1 AND (%s)LIMIT 1;", sqlQuery)

	smt, err := db.db.Prepare(sqlQuery)
	if err != nil {
		return nil, errors.Wrap(err, "get url prepare error")
	}
	res, err := smt.Query()
	if err != nil {
		return nil, errors.Wrap(err, "get url error")
	}

	if res.Next() {
		return CastUser(res)
	}
	return nil, errors.New("user not found")
}

func (db *Mysql) GetIdentity(login string) (*user.User, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM users WHERE status=1 AND (phone_number = ? OR email = ?) LIMIT 1;")

	smt, err := db.db.Prepare(sqlQuery)
	if err != nil {
		return nil, errors.Wrap(err, "get url prepare error")
	}
	res, err := smt.Query(login, login)
	if err != nil {
		return nil, errors.Wrap(err, "get url error")
	}

	if res.Next() {
		return CastUser(res)
	}
	return nil, errors.New("user not found")
}

func CastUser(res *sql.Rows) (*user.User, error) {
	var id int64
	var firstName string
	var secondName string
	var email string
	var phoneNumber string
	var password string
	var status int
	var updatedAt time.Time
	var createdAt time.Time
	err := res.Scan(&id, &firstName, &secondName, &email, &phoneNumber, &password, &status, &updatedAt, &createdAt)
	if err != nil {
		return nil, errors.Wrap(err, "scan error")
	}
	return &user.User{
		ID:          id,
		FirstName:   firstName,
		SecondName:  secondName,
		Email:       email,
		PhoneNumber: phoneNumber,
		Password:    password,
		Status:      status,
		UpdatedAt:   model.Timestamp{Time: updatedAt},
		CreatedAt:   model.Timestamp{Time: createdAt},
	}, nil
}
