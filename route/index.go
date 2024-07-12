package route

import (
	"github.com/gin-gonic/gin"
	"polarstar/api"
	"polarstar/api/futures"
	"polarstar/middleware"
)

// 初始路由
func InitRoutes(r *gin.Engine) {
	r.GET("/", api.Hello)
	apiGroup := r.Group("/api")
	spotGroup := apiGroup.Group("/spot")
	sporUser := spotGroup.Group("/user")
	futuresGroup := apiGroup.Group("/um")
	futuresUser := futuresGroup.Group("/user")
	deliveryGroup := apiGroup.Group("/cm")
	spotGroup.Use(middleware.IPAuthMiddleware())
	{

	}
	sporUser.Use(middleware.IPAuthMiddleware())
	{
		sporUser.GET("/userInfo")
	}
	futuresGroup.Use(middleware.IPAuthMiddleware())
	{
		futuresGroup.GET("/getHoldOrderList", futures.GetHoldOrderList)
		futuresGroup.POST("/getHoldOrderDetail", futures.GetHoldOrderDetail)
		futuresGroup.GET("/getKline", futures.GetUMKline)
		futuresGroup.POST("/changeLever", futures.ChangeLever)
		futuresGroup.POST("/umInOrder", futures.ContractInOrder)
		futuresGroup.POST("/cancelHoldOrder", futures.CancelHoldOrderList)
		futuresGroup.POST("/cancelHoldOrderDetail", futures.CancelHoldOrderDetail)
		futuresGroup.GET("/getUMOrder", futures.GetUMOrder)
		futuresGroup.GET("/getUMOrderDetail", futures.GetUMOrderDetail)
		futuresGroup.GET("/historyKline", futures.Tet)
	}
	futuresUser.Use(middleware.IPAuthMiddleware())
	{
		futuresUser.GET("/userInfo", futures.GetUMUserInfo)
	}
	deliveryGroup.Use(middleware.IPAuthMiddleware())
	{

	}

}
