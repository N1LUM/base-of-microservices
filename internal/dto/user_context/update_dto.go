package user_context

type UpdateUserInput struct {
	Name     *string `json:"name" binding:"omitempty,min=2"`
	Username *string `json:"username" binding:"omitempty,min=3"`
	Password *string `json:"password" binding:"omitempty,min=6"`
}
