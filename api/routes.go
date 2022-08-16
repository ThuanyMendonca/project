package api

import (
	"github.com/ThuanyMendonca/project/api/user"
	"github.com/gin-gonic/gin"
)

func Router(e *gin.Engine) {
	authNotRequired := e.Group("/api/v1")
	user.Router(authNotRequired.Group("/user"))
}
