//	回复消息体
package models

//	JSON格式的消息返回模型
type JSONResponse struct {
	Status int
	Data   interface{}
}

//	操作成功
func GetSuccessResponse(data interface{}) *JSONResponse {
	res := &JSONResponse{Status: 0, Data: data}
	return res
}

//	操作失败
func GetFailedResponse(status int, data interface{}) *JSONResponse {
	res := &JSONResponse{Status: status, Data: data}
	return res
}
