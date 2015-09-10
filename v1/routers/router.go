package routers

import (
	"github.com/astaxie/beego"
	"ums/v1/controllers"
)

func init() {
	//与设备通信路由
	beego.Router("/UMS/UserRegister.do", &controllers.RegisterController{})
	beego.Router("/UMS/UserAuth.do", &controllers.UserAuthController{})
	beego.Router("/UMS/UserInfoUpdate.do", &controllers.UpdateController{})
	beego.Router("/UMS/UserDeauth.do", &controllers.DeauthController{})
	//beego.Router("/", &controllers.MainController{})
}
