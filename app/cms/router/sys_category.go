package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/cms/models"
	"go-admin/app/cms/service/dto"
	"go-admin/common/middleware"

	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/common/actions"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSysCategoryRouter)
}

// 需认证的路由代码
func registerSysCategoryRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	r := v1.Group("/syscategory").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		model := &models.SysCategory{}
		r.GET("", actions.PermissionAction(), actions.IndexAction(model, new(dto.SysCategorySearch), func() interface{} {
			list := make([]models.SysCategory, 0)
			return &list
		}))
		r.GET("/:id", actions.PermissionAction(), actions.ViewAction(new(dto.SysCategoryById), nil))
		r.POST("", actions.CreateAction(new(dto.SysCategoryControl)))
		r.PUT("/:id", actions.PermissionAction(), actions.UpdateAction(new(dto.SysCategoryControl)))
		r.DELETE("", actions.PermissionAction(), actions.DeleteAction(new(dto.SysCategoryById)))
	}
}
