//	基础接口
package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"learn/models"
	"time"
)

//  基础接口模型
type baseController struct {
	beego.Controller
}

//	操作成功
func (this *baseController) SendSuccessJSON(data interface{}) {
	this.Data["json"] = models.GetSuccessResponse(data)
	this.ServeJson()
}

//	操作失败
func (this *baseController) SendFailedJSON(status int, data interface{}) {
	this.Data["json"] = models.GetFailedResponse(status, data)
	this.ServeJson()
}

//	服务器内部发生错误
func (this *baseController) SendErrorJSON(err error) {
	this.SendFailedJSON(3, err.Error())
}

//  返回appid 和 sessid， 并计算签名保存在session中
//	@Return appid 时间戳
//	@Return sessid sessionID
func (this *baseController) SetSignature() (appid, sessid string) {
	appid = fmt.Sprintf("%d", int(time.Now().Unix()/3600))
	sessid = this.Ctx.Input.Cookie("beegosessionID")

	if this.GetSession("token") == nil {
		return
	}
	//	计算签名并保存在session中
	signature := models.GetSignature(this.GetSession("token").(string), appid, sessid)
	if this.GetSession("signature") != nil {
		this.DelSession("signature")
	}
	this.SetSession("signature", signature)
	return
}
