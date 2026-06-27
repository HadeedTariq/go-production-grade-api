package authDto

type SignupRequest struct {
	Name       string `json:"name" validate:"required,min=3,max=100"`
	Username   string `json:"username" validate:"required,min=3,max=50"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8"`
	Profession string `json:"profession" validate:"required"`
}
