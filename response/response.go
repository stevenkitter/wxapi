package response

const (
	//ResponseCodeOK 成功
	ResponseCodeOK = 200
	//ResponseCodeOKChanged 成功 数据变了
	ResponseCodeOKChanged = 211
	//ResponseCodeFail 失败
	ResponseCodeFail = 400
	//ResponseCodeTokenInValid token失效
	ResponseCodeTokenInValid = 606
)

//Response 接口统一返回数据格式
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//MenuResponse 接口统一返回数据格式 点菜
type MenuResponse struct {
	Code    int         `json:"code"`
	Version int         `json:"version"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
