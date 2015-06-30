package controllers

import (
	"adexchange/engine"
	"adexchange/lib"
	m "adexchange/models"
	"github.com/astaxie/beego"
	"time"
)

type RequestController struct {
	BaseController
}

//Request Ad
func (this *RequestController) RequestAd() {
	adRequest := m.AdRequest{}
	adResponse := new(m.AdResponse)
	beego.Debug("Enter Request ad")
	if err := this.ParseForm(&adRequest); err != nil {

		adResponse.StatusCode = lib.ERROR_PARSE_REQUEST
	} else {
		adRequest.RequestTime = time.Now().Unix()
		tmp := engine.InvokeDemand(&adRequest)

		if tmp == nil {
			adResponse.StatusCode = lib.ERROR_UNKNON_ERROR
		} else {
			adResponse = tmp
		}
		adRequest.StatusCode = adResponse.StatusCode
		SendLog(adRequest, 1)
		//if err != nil {
		//	beego.Debug("Enter sss ad")
		//	if e, ok := err.(*lib.SysError); ok {
		//		adResponse.StatusCode = e.ErrorCode
		//	} else {
		//		adResponse.StatusCode = lib.ERROR_UNKNON_ERROR
		//	}
		//	beego.Debug("Enter ssaass ad")
		//}
	}

	this.Data["json"] = adResponse.GenerateCommonResponse()
	this.ServeJson()

}
