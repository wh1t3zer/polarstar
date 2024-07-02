package route

import (
	"github.com/gin-gonic/gin"
	"polarstar/api/futures"
	"polarstar/middleware"
)

// 初始路由
func InitRoutes(r *gin.Engine) {
	apiGroup := r.Group("/api")
	spotGroup := apiGroup.Group("/spot")
	futuresGroup := apiGroup.Group("/um")
	deliveryGroup := apiGroup.Group("/cm")
	spotGroup.Use(middleware.IPAuthMiddleware())
	{

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
	}
	deliveryGroup.Use(middleware.IPAuthMiddleware())
	{

	}

}
