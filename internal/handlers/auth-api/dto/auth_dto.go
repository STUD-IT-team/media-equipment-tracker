package dto

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,min=4,max=50" example:"test@mail.ru"`
	Password string `json:"password" binding:"required,min=4" example:"12345678"`
}

type LoginUserResponse struct {
	AccessToken string `json:"access_token" example:"7249ede8-5083-4bd6-ad09-0b5fa3c5f2de"`
}

type RegisterUserRequest struct {
	FullName string `json:"full_name" binding:"required,min=4,max=50" example:"Ivan Ivanov Ivanovich"`
	Email    string `json:"email" binding:"required,min=4,max=50" example:"test@mail.ru"`
	Password string `json:"password" binding:"required,min=4" example:"12345678"`
}
