#!/bin/bash

# Check if no parameters are passed
if [ -z "$1" ]; then
  echo "Usage: $0 <module_name> [package_dir]"
  exit 1
fi

MODULE_NAME=$1
PACKAGE_DIR=$2

if [ -n "$PACKAGE_DIR" ]; then
  MODULE_DIR="modules/$PACKAGE_DIR/$MODULE_NAME"
else
  MODULE_DIR="modules/$MODULE_NAME"
fi

# NEW_MODULE_NAME="$(tr '[:lower:]' '[:upper:]' <<< ${MODULE_NAME:0:1})${MODULE_NAME:1}"
NEW_MODULE_NAME=""
IFS='.' read -ra PARTS <<< "$MODULE_NAME"
for PART in "${PARTS[@]}"; do
  NEW_MODULE_NAME+="$(tr '[:lower:]' '[:upper:]' <<< ${PART:0:1})${PART:1}"
done

PACKAGE_NAME="${MODULE_NAME//./}"

# Create module directory if it doesn't exist
mkdir -p "$MODULE_DIR"

# Create file controller
cat <<EOL > "$MODULE_DIR/${MODULE_NAME}.controllers.go"
package $PACKAGE_NAME

import (
	"context"

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
    return &Controller{ctx, db, store}
}

// Get${NEW_MODULE_NAME} godoc
// @Summary      Get ${NEW_MODULE_NAME}
// @Description  Returns a hello world message
// @Tags         ${NEW_MODULE_NAME}
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /api/${NEW_MODULE_NAME} [get]
func (${NEW_MODULE_NAME}Controller *Controller) Get${NEW_MODULE_NAME}(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}
EOL

# Create file routes
cat <<EOL > "$MODULE_DIR/${MODULE_NAME}.routes.go"
package $PACKAGE_NAME

import "github.com/gin-gonic/gin"

// Router Router
type Router struct {
	${NEW_MODULE_NAME}Controller Controller
}

// NewRouter NewRouter
func NewRouter(${NEW_MODULE_NAME}Controller Controller) Router {
	return Router{${NEW_MODULE_NAME}Controller: ${NEW_MODULE_NAME}Controller}
}

// RegisterRoutes RegisterRoutes
func (router *Router) RegisterRoutes(rGroup *gin.RouterGroup) {
	${NEW_MODULE_NAME}Router := rGroup.Group("/${NEW_MODULE_NAME}")
	{
		${NEW_MODULE_NAME}Router.GET("", router.${NEW_MODULE_NAME}Controller.Get${NEW_MODULE_NAME})
	}
}
EOL

# Create file types
cat <<EOL > "$MODULE_DIR/${MODULE_NAME}.types.go"
package $PACKAGE_NAME

EOL

# Create file helpers
cat <<EOL > "$MODULE_DIR/${MODULE_NAME}.helpers.go"
package $PACKAGE_NAME

EOL

echo "Module '$MODULE_NAME' created successfully in $MODULE_DIR"