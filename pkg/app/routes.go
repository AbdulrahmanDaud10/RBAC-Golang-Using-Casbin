package app

import (
	"log"

	"github.com/AbdulrahmanDaud10/RBAC-Casbin-Golang/pkg/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetUpRoutes(db *gorm.DB) {
	httpRouter := gin.Default()

	userRepository := repository.NewUserRepository(db)

	if err := userRepository.Migrate(); err != nil {
		log.Fatal("User migrate err", err)
	}

	httpRouter.Run()
}
