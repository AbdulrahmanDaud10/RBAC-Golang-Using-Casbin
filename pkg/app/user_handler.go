package app

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/AbdulrahmanDaud10/RBAC-Casbin-Golang/pkg/api"
	"github.com/AbdulrahmanDaud10/RBAC-Casbin-Golang/pkg/repository"
	"github.com/casbin/casbin/v2"

	"github.com/gin-gonic/gin"
)

// UserController : represent the user's controller contract
type UserService interface {
	AddUser(enforcer *casbin.Enforcer) gin.HandlerFunc
	GetUser(*gin.Context)
	GetAllUser(*gin.Context)
	SignInUser(*gin.Context)
	UpdateUser(*gin.Context)
	DeleteUser(*gin.Context)
}

type userService struct {
	userRepo repository.UserRepository
}

// NewUserController -> returns new user controller
func NewUserController(repo repository.UserRepository) UserService {
	return userService{
		userRepo: repo,
	}
}

func (h userService) GetAllUser(ctx *gin.Context) {
	fmt.Println(ctx.Get("userID"))
	user, err := h.userRepo.GetAllUser()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	ctx.JSON(http.StatusOK, user)

}

func (h userService) GetUser(ctx *gin.Context) {
	id := ctx.Param("user")
	intID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.userRepo.GetUser(intID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	ctx.JSON(http.StatusOK, user)

}

func (h userService) SignInUser(ctx *gin.Context) {
	var user api.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	dbUser, err := h.userRepo.GetByEmail(user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "No Such User Found"})
		return

	}
	if isTrue := api.ComparePassword(dbUser.Password, user.Password); isTrue {
		fmt.Println("user before", dbUser.ID)
		token := api.GenerateToken(dbUser.ID)
		ctx.JSON(http.StatusOK, gin.H{"msg": "Successfully SignIN", "token": token})
		return
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Password not matched"})

}

func (h userService) AddUser(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user api.User
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		api.HashPassword(&user.Password)
		user, err := h.userRepo.AddUser(user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return

		}
		enforcer.AddGroupingPolicy(fmt.Sprint(user.ID), user.Role)
		user.Password = ""
		ctx.JSON(http.StatusOK, user)

	}
}

func (h userService) UpdateUser(ctx *gin.Context) {
	var user api.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := ctx.Param("user")
	intID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	user.ID = uint(intID)
	user, err = h.userRepo.UpdateUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	ctx.JSON(http.StatusOK, user)

}

func (h userService) DeleteUser(ctx *gin.Context) {
	var user api.User
	id := ctx.Param("user")
	intID, _ := strconv.Atoi(id)
	user.ID = uint(intID)
	user, err := h.userRepo.DeleteUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	ctx.JSON(http.StatusOK, user)

}
