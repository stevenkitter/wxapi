package response

//UserLoginResponse 登陆使用的
type UserLoginResponse struct {
	UserPhone string `json:"user_phone"`
	Token     string `json:"token"`
}
