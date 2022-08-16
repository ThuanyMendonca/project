package user

import (
	"github.com/ThuanyMendonca/project/repository"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.RouterGroup) {
	b := NewUserBusiness(repository.UserRepository{})
	c := NewUserController(b)

	r.POST("/create", c.Post)

}
