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

type RequestAddUserToList struct {
	UserID   uint `json:"userId" validate:"required" example:"3" format:"uint"`
	FriendID uint `json:"friendId" validate:"required" example:"5" format:"uint"`
}

type RequestUserUpdatePassword struct {
	UserID      uint   `json:"userId" validate:"required" example:"3" format:"uint"`
	Username    string `json:"username" validate:"required" example:"johndoe" format:"string"`
	OldPassword string `json:"oldPassword" validate:"required" format:"string"`
	NewPassword string `json:"newPassword" validate:"required" format:"string"`
}

type RequestUserUpdateUsername struct {
	UserID      uint   `json:"userId" validate:"required" example:"3" format:"uint"`
	Username    string `json:"username" validate:"required" example:"johndoe" format:"string"`
	Password    string `json:"password" validate:"required" format:"string"`
	NewUsername string `json:"newUsername" validate:"required" example:"johndoe" format:"string"`
}
