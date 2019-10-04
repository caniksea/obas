package domain

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Register struct {
	Email string `json:"email"`
}
type ForgetPassword struct {
	Email string `json:"email"`
}
type LoginToken struct {
	Email string `json:"email"`
	Token string `json:"token"`
}
