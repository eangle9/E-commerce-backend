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
	customizer3 = g.Validator(request.RefreshRequest{})
	customizer4 = g.Validator(request.ProductCategoryRequest{})
	customizer5 = g.Validator(request.ColorRequest{})
	customizer6 = g.Validator(request.ProductRequest{})
	customizer7 = g.Validator(request.ProductItemRequest{})
	customizer8 = g.Validator(request.CartRequest{})
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
	api.GET("/:id", protectedMiddleware(), u.getUserHandler)
	api.PUT("/update/:id", protectedMiddleware(), u.updateUserHandler)
	api.DELETE("/delete/:id", protectedMiddleware(), u.deleteUserHandler)
	api.POST("/token", protectedMiddleware(), u.refreshTokenHandler)
}

// registerHandler  godoc
//
//	@Summary		Insert user
//	@Description	Add a new user
//	@Tags			user
//	@ID				register-user
//	@Accept			json
//	@Produce		json
//	@Param			user	body		request.SignUpRequest	true	"User data"
//	@Success		201		{object}	response.Response
//	@Router			/user/register [post]
func (u UserController) registerHandler(c *gin.Context) {
	var request request.SignUpRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "failed to decode json request body",
		}
		// c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		c.Set("error", errorResponse)
		return
	}

	validate := c.MustGet("validator").(*validator.Validate)
	if err := validate.Struct(request); err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.ValidationError,
			ErrorMessage: customizer1.DecryptErrors(err),
		}
		// c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		c.Set("error", errorResponse)
		return
	}

	signUpRequest := request
	resp := u.userService.SignUp(signUpRequest)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}

// LoginHandler   godoc
//
//	@Summary		Login user
//	@Description	User login by email and password
//	@Tags			user
//	@ID				login-user
//	@Accept			json
//	@Produce		json
//	@Param			user	body		request.LoginRequest	true	"Login data"
//	@Success		200		{object}	response.Response
//	@Router			/user/login [post]
func (u UserController) LoginHandler(c *gin.Context) {
	var request request.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "failed to decode json request body",
		}
		// c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		c.Set("error", errorResponse)
		return
	}

	validate := c.MustGet("validator").(*validator.Validate)
	if err := validate.Struct(request); err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.ValidationError,
			ErrorMessage: customizer2.DecryptErrors(err),
		}
		// c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		c.Set("error", errorResponse)
		return
	}

	resp := u.userService.LoginUser(request)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}
	c.JSON(resp.Status, resp)
}

// listUserHandler godoc
//
//	@Summary		List users
//	@Description	Retrieves a list of users. Requires authentication with JWT token.
//	@Tags			user
//	@ID				list-users
//	@Produce		json
//	@Security		JWT
//	@Success		200	{object}	response.Response
//	@Router			/user/list [get]
func (u UserController) listUserHandler(c *gin.Context) {
	// id, _ := c.Get("userId")
	// role, _ := c.Get("role")
	// fmt.Printf("id: %d and role: %s \n", id, role)
	resp := u.userService.GetUsers()
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}
	c.JSON(resp.Status, resp)
}

// getUserHandler    godoc
//
//	@Summary		Get user
//	@Description	Get a single user by id
//	@Tags			user
//	@ID				get-user-by-id
//	@Produce		json
//	@Security		JWT
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	response.Response
//	@Router			/user/{id} [get]
func (u UserController) getUserHandler(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "invalid id.Please enter a valid integer id",
		}
		// c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		c.Set("error", errorResponse)
		return
	}

	resp := u.userService.GetUser(id)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}
	c.JSON(resp.Status, resp)

}

// updateUserHandler   godoc
//
//	@Summary		Update user
//	@Description	update user by id
//	@Tags			user
//	@ID				update-user-by-id
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Param			id		path		int					true	"UserID"
//	@Param			user	body		utils.UpdateUser	true	"Update user data"
//	@Success		200		{object}	response.Response
//	@Router			/user/update/{id} [put]
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
		// c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		c.Set("error", errorResponse)
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "failed to decode json request body",
		}
		// c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		c.Set("error", errorResponse)
		return
	}

	resp := u.userService.UpdateUser(id, user)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// deleteUserHandler  godoc
//
//	@Summary		Delete user
//	@Description	delete user by id
//	@Tags			user
//	@ID				delete-user-by-id
//	@Produce		json
//	@Security		JWT
//	@Param			id	path		int	true	"UserID"
//	@Success		200	{object}	response.Response
//	@Router			/user/delete/{id} [delete]
func (u UserController) deleteUserHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "invalid id.Please enter a valid integer id",
		}
		// c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		c.Set("error", errorResponse)
		return
	}

	resp := u.userService.DeleteUser(id)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}
	c.JSON(resp.Status, resp)
}

// refreshTokenHandler  godoc
//
//	@Summary		Refresh token
//	@Description	refresh the expired access token
//	@Tags			user
//	@ID				refresh-access-token
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Param			token	body		request.RefreshRequest	true	"Refresh token"
//	@Success		200		{object}	response.Response
//	@Router			/user/token [post]
func (u UserController) refreshTokenHandler(c *gin.Context) {
	var rfToken request.RefreshRequest

	if err := c.ShouldBindJSON(&rfToken); err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "failed to decode json request body",
		}
		// c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		c.Set("error", errorResponse)
		return
	}

	validate := c.MustGet("validator").(*validator.Validate)
	if err := validate.Struct(rfToken); err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.ValidationError,
			ErrorMessage: customizer3.DecryptErrors(err),
		}
		// c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		c.Set("error", errorResponse)
		return
	}

	resp := u.userService.RefreshToken(rfToken)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}
	c.JSON(resp.Status, resp)
}
