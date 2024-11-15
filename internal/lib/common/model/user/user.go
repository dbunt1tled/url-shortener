package user

import "go_first/internal/lib/common/model"

type User struct {
	ID          int64           `json:"id" jsonapi:"primary,url"`
	FirstName   string          `json:"firstName" jsonapi:"attr,firstName"`
	SecondName  string          `json:"secondName" jsonapi:"attr,secondName"`
	Email       string          `json:"email" jsonapi:"attr,email"`
	PhoneNumber string          `json:"phoneNumber" jsonapi:"attr,phoneNumber"`
	Password    string          `json:"password" jsonapi:"attr,password"`
	Status      int             `json:"status" jsonapi:"attr,status"`
	CreatedAt   model.Timestamp `json:"created_at" jsonapi:"attr,createdAt"`
	UpdatedAt   model.Timestamp `json:"updated_at" jsonapi:"attr,updatedAt"`
}

const (
	StatusActive   = 1
	StatusInactive = 2
)
