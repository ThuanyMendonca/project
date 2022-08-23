package api

import (
	"github.com/ThuanyMendonca/project/api/balance"
	"github.com/ThuanyMendonca/project/api/login"
	"github.com/ThuanyMendonca/project/api/transaction"
	"github.com/ThuanyMendonca/project/api/user"
	"github.com/ThuanyMendonca/project/middleware"
	"github.com/gin-gonic/gin"
)

func Router(e *gin.Engine) {
	// authNotRequired := e.Group("/")
	authRequired := e.Group("/api/v1")
	authRequired.Use(middleware.Auth())

	authNotRequired := e.Group("/api/v1")

	user.Router(authNotRequired.Group("/user"))
	transaction.Router(authRequired.Group("/transaction"))
	balance.Router(authRequired.Group("/balance"))
	login.Router(authNotRequired.Group("/login"))
}
