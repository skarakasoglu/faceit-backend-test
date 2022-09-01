package user

import "time"

type User struct {
	Id        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Nickname  string    `json:"nickname"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserResponse struct {
	User
}

type UpdateUserResponse struct {
	User
}

type DeleteUserResponse struct {
	Id string `json:"id"`
}

type GetUsersManyResponse struct {
	Users []User `json:"users"`
}
