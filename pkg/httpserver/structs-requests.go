package httpserver

type RequestUserRegister struct {
	Username string `json:"username" validate:"required" example:"johndoe" format:"string"`
	Email    string `json:"email" validate:"required" example:"johndoe@example.com" format:"string"`
	Password string `json:"password" validate:"required" format:"string"`
}

type RequestUserAuthenticate struct {
	Username string `json:"username" validate:"required" example:"johndoe" format:"string"`
	Password string `json:"password" validate:"required" format:"string"`
}
