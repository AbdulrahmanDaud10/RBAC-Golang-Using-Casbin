package app

import (
	"fmt"
	"log"

	"github.com/AbdulrahmanDaud10/RBAC-Casbin-Golang/pkg/repository"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetUpRoutes(db *gorm.DB) {
	httpRouter := gin.Default()

	// Initialize  casbin adapter
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize casbin adapter: %v", err))
	}

	// Load model configuration file and policy store adapter
	enforcer, err := casbin.NewEnforcer("cmd/config/rbac_model.conf", adapter)
	if err != nil {
		panic(fmt.Sprintf("failed to create casbin enforcer: %v", err))
	}

	//add policy
	if hasPolicy := enforcer.HasPolicy("doctor", "report", "read"); !hasPolicy {
		enforcer.AddPolicy("doctor", "report", "read")
	}
	if hasPolicy := enforcer.HasPolicy("doctor", "report", "write"); !hasPolicy {
		enforcer.AddPolicy("doctor", "report", "write")
	}
	if hasPolicy := enforcer.HasPolicy("patient", "report", "read"); !hasPolicy {
		enforcer.AddPolicy("patient", "report", "read")
	}

	userRepository := repository.NewUserRepository(db)

	if err := userRepository.Migrate(); err != nil {
		log.Fatal("User migrate err", err)
	}

	userController := NewUserController(userRepository)

	apiRoutes := httpRouter.Group("/api")

	{
		apiRoutes.POST("/register", userController.AddUser(enforcer))
		apiRoutes.POST("/signin", userController.SignInUser)
	}

	userProtectedRoutes := apiRoutes.Group("/users", AuthorizeJWT())
	{
		userProtectedRoutes.GET("/", Authorize("report", "read", enforcer), userController.GetAllUser)
		userProtectedRoutes.GET("/:user", Authorize("report", "read", enforcer), userController.GetUser)
		userProtectedRoutes.PUT("/:user", Authorize("report", "write", enforcer), userController.UpdateUser)
		userProtectedRoutes.DELETE("/:user", Authorize("report", "write", enforcer), userController.DeleteUser)
	}

	httpRouter.Run()
}
