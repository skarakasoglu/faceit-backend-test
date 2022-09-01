package user

type CreateUserRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Nickname  string `json:"nickname" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Country   string `json:"country" binding:"required"`
}

type UpdateUserRequest struct {
	Id        string `json:"-"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Nickname  string `json:"nickname" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Country   string `json:"country" binding:"required"`
}

type DeleteUserByIdRequest struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type GetUsersManyRequest struct {
	Page    int  `uri:"page"`
	PerPage int  `uri:"perPage"`
	Filter  User `uri:"filter"`
}
