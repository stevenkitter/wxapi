package response

//ResOK 一般的成功返回
func ResOK(msg string) Response {
	return Response{Code: ResponseCodeOK, Message: msg, Data: nil}
}

//ResFAIL 一般的失败返回
func ResFAIL(msg string) Response {
	return Response{Code: ResponseCodeFail, Message: msg, Data: nil}
}

//TokenInValid token
func TokenInValid(msg string) Response {
	return Response{Code: ResponseCodeTokenInValid, Message: msg, Data: nil}
}
