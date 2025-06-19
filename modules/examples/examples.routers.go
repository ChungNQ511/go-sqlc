package examples

import "github.com/gin-gonic/gin"

// Router Router
type Router struct {
	exampleController Controller
}

// NewRouter NewRouter
func NewRouter(exampleController Controller) Router {
	return Router{exampleController: exampleController}
}

// RegisterRoutes RegisterRoutes
func (router *Router) RegisterRoutes(rGroup *gin.RouterGroup) {
	exampleRouter := rGroup.Group("/examples")
	{
		exampleRouter.GET("", router.exampleController.GetExample)
	}
}
