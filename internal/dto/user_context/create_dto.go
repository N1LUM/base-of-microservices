package user_context

type CreateUserInput struct {
	Name     string `json:"name" binding:"required,min=2"`
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}
