package domain

import (
	"time"
)

type User struct {
	Id         string    `dynamodbav:"id"`
	Name       string    `dynamodbav:"name"`
	Email      string    `dynamodbav:"email"`
	Password   string    `dynamodbav:"password"`
	Created_at time.Time `dynamodbav:"created_at"`
	Updated_at time.Time `dynamodbav:"updated_at"`
}
