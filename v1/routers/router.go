package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/toolbox"
	"ums/v1/controllers"
)

func init() {
	toolbox.AddHealthCheck("database", &controllers.DatabaseCheck{})
	//与设备通信路由
	beego.Router("/UMS/UserRegister.do", &controllers.RegisterController{})
	beego.Router("/UMS/UserAuth.do", &controllers.UserAuthController{})
	beego.Router("/UMS/UserInfoUpdate.do", &controllers.UpdateController{})
	beego.Router("/UMS/UserDeauth.do", &controllers.DeauthController{})
	//beego.Router("/", &controllers.MainController{})
}
