package controller

import (
	"Eccomerce-website/internal/core/common/router"
	"Eccomerce-website/internal/core/common/utils"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/service"
	"Eccomerce-website/internal/infra/middleware"
	"strconv"

	// "Eccomerce-website/internal/infra/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golodash/galidator"
)

var (
	g = galidator.New().CustomMessages(galidator.Messages{
		"required": "$field is required",
	})
	customizer1 = g.Validator(request.SignUpRequest{})
	customizer2 = g.Validator(request.LoginRequest{})
)

type UserController struct {
	engine      *router.Router
	userService service.UserService
}

func NewUserController(engine *router.Router, userService service.UserService) *UserController {
	return &UserController{
		engine:      engine,
		userService: userService,
	}
}

func (u *UserController) InitRouter() {
	protectedMiddleware := middleware.ProtectedMiddleware
	r := u.engine
	api := r.Group("/user")

	api.POST("/register", u.registerHandler)
	api.POST("/login", u.LoginHandler)
	api.GET("/list", protectedMiddleware(), u.listUserHandler)
	api.GET("/:id", u.getUserHandler)
	api.PUT("/update/:id", u.updateUserHandler)
	api.DELETE("/delete/:id", u.deleteUserHandler)
}

func (u UserController) registerHandler(c *gin.Context) {
	var request request.SignUpRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "failed to decode json request body",
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		return
	}

	validate := c.MustGet("validator").(*validator.Validate)
	if err := validate.Struct(request); err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.ValidationError,
			ErrorMessage: customizer1.DecryptErrors(err),
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		return
	}

	signUpRequest := request
	resp := u.userService.SignUp(signUpRequest)
	c.JSON(resp.Status, resp)
}

func (u UserController) LoginHandler(c *gin.Context) {
	var request request.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "failed to decode json request body",
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		return
	}

	validate := c.MustGet("validator").(*validator.Validate)
	if err := validate.Struct(request); err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.ValidationError,
			ErrorMessage: customizer2.DecryptErrors(err),
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		return
	}

	resp := u.userService.LoginUser(request)
	c.JSON(resp.Status, resp)
}

func (u UserController) listUserHandler(c *gin.Context) {
	resp := u.userService.GetUsers()
	c.JSON(resp.Status, resp)
}

func (u UserController) getUserHandler(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "invalid id.Please enter a valid integer id",
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		return
	}

	resp := u.userService.GetUser(id)
	c.JSON(resp.Status, resp)

}

func (u UserController) updateUserHandler(c *gin.Context) {
	var user utils.UpdateUser

	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "invalid id.Please enter a valid integer id",
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "failed to decode json request body",
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		return
	}

	resp := u.userService.UpdateUser(id, user)
	c.JSON(http.StatusOK, resp)
}

func (u UserController) deleteUserHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "invalid id.Please enter a valid integer id",
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		return
	}

	resp := u.userService.DeleteUser(id)
	c.JSON(resp.Status, resp)
}
