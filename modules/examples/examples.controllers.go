package examples

import (
	"context"
	"net/http"

	db "example.com/m/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Controller Controller
type Controller struct {
	ctx   context.Context
	db    *db.Queries
	store *db.Store
}

// NewController NewController
func NewController(ctx context.Context, db *db.Queries, store *db.Store) *Controller {
	return &Controller{ctx: ctx, db: db, store: store}
}

// GetExample godoc
// @Summary      Get example
// @Description  Returns a hello world message
// @Tags         example
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /api/example [get]
func (exampleController *Controller) GetExample(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}
