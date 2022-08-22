package user

import (
	"github.com/ThuanyMendonca/project/dependency"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.RouterGroup) {
	b := NewUserBusiness(dependency.UserRepository)
	c := NewUserController(b)

	r.POST("/create", c.Post)
	r.PUT("/inactive/:id", c.Inactive)

}
