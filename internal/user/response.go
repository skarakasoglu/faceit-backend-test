package user

import "time"

// User represents user api model
// @Description user model
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

// CreateUserResponse create user endpoint response model containing the user information
// @Description create user endpoint response model containing the user information
type CreateUserResponse struct {
	User
}

// UpdateUserResponse update user endpoint response model containing updated user information
// @Description update user endpoint response model containing updated user information
type UpdateUserResponse struct {
	User
}

// DeleteUserResponse delete user response model containing the id of deleted user
// @Description delete user response model containing the id of deleted user
type DeleteUserResponse struct {
	Id string `json:"id"`
}

// GetUsersManyResponse get users response model that contains the users returned
// @Description get users response model that contains the users returned
type GetUsersManyResponse struct {
	Users []User `json:"users"`
}
