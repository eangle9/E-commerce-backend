package controller

import (
	"Eccomerce-website/internal/core/common/router"
	"Eccomerce-website/internal/core/entity"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/port/service"
	"Eccomerce-website/internal/infra/middleware"
	"strconv"

	// "Eccomerce-website/internal/infra/middleware"
	"net/http"

	"github.com/gin-gonic/gin"

	// "github.com/golodash/galidator"
	pagination "github.com/webstradev/gin-pagination"
)

// var (
// 	g = galidator.New().CustomMessages(galidator.Messages{
// 		"required": "$field is required",
// 	})
// 	customizer1  = g.Validator(request.SignUpRequest{})
// 	customizer2  = g.Validator(request.LoginRequest{})
// 	customizer3  = g.Validator(request.RefreshRequest{})
// 	customizer4  = g.Validator(request.ProductCategoryRequest{})
// 	customizer5  = g.Validator(request.ColorRequest{})
// 	customizer6  = g.Validator(request.ProductRequest{})
// 	customizer8  = g.Validator(request.CartRequest{})
// 	customizer9  = g.Validator(request.SizeRequest{})
// 	customizer10 = g.Validator(request.ReviewRequest{})
// )

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
	pagination.New("page", "page_size", "1", "10", 1, 150)
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
// @Summary		    Insert user
// @Description	    Add a new user
// @Tags			user
// @ID				register-user
// @Accept			json
// @Produce		    json
// @Param			user	body		request.SignUpRequest	true	"User data"
// @Success		    201		{object}	response.Response
// @Router			/user/register [post]
func (u UserController) registerHandler(c *gin.Context) {
	var signUpRequest request.SignUpRequest

	if err := c.ShouldBindJSON(&signUpRequest); err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "failed to decode json request body").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		return
	}

	// validate := c.MustGet("validator").(*validator.Validate)
	// if err := validate.Struct(request); err != nil {
	// 	errorResponse := response.Response{
	// 		Status:       http.StatusBadRequest,
	// 		ErrorType:    errorcode.ValidationError,
	// 		ErrorMessage: customizer1.DecryptErrors(err),
	// 	}
	// 	// c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
	// 	c.Set("error", errorResponse)
	// 	return
	// }

	// signUpRequest := request
	resp, err := u.userService.SignUp(signUpRequest)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}

// LoginHandler   godoc
// @Summary		  Login user
// @Description	  User login by email and password
// @Tags		  user
// @ID			  login-user
// @Accept		  json
// @Produce		  json
// @Param		  user	    body		request.LoginRequest	true	"Login data"
// @Success		  200		{object}	response.Response
// @Router		  /user/login [post]
func (u UserController) LoginHandler(c *gin.Context) {
	var request request.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "failed to bind json request body").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		return
	}

	// validate := c.MustGet("validator").(*validator.Validate)
	// if err := validate.Struct(request); err != nil {
	// 	errorResponse := response.Response{
	// 		Status:       http.StatusBadRequest,
	// 		ErrorType:    errorcode.ValidationError,
	// 		ErrorMessage: customizer2.DecryptErrors(err),
	// 	}
	// 	// c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
	// 	c.Set("error", errorResponse)
	// 	return
	// }

	resp, err := u.userService.LoginUser(request)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(resp.StatusCode, resp)
}

// listUserHandler godoc
// @Summary		   List users
// @Description	   Retrieves a list of users. Requires authentication with JWT token.
// @Tags		   user
// @ID			   list-users
// @Produce		   json
// @Security	   JWT
// @Success		   200	{object}	response.Response
// @Router		   /user/list [get]
func (u UserController) listUserHandler(c *gin.Context) {
	var paginationQuery request.PaginationQuery
	if err := c.ShouldBindQuery(&paginationQuery); err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "failed to bind pagination query to struct").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		return
	}

	resp, err := u.userService.GetUsers(paginationQuery)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}

// getUserHandler    godoc
// @Summary		     Get user
// @Description	     Get a single user by id
// @Tags			 user
// @ID				 get-user-by-id
// @Produce		     json
// @Security		 JWT
// @Param			 id	path		int	true	"User ID"
// @Success		     200	        {object}	response.Response
// @Router			 /user/{id} [get]
func (u UserController) getUserHandler(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "invalid id.Please enter a valid integer id").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		return
	}

	resp, err := u.userService.GetUser(id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)

}

// updateUserHandler   godoc
// @Summary		       Update user
// @Description	       update user by id
// @Tags			   user
// @ID				   update-user-by-id
// @Accept			   json
// @Produce		       json
// @Security		   JWT
// @Param			   id		path		int					true	"UserID"
// @Param			   user	    body		request.UpdateUser	true	"Update user data"
// @Success		       200		{object}	response.Response
// @Router			   /user/update/{id} [put]
func (u UserController) updateUserHandler(c *gin.Context) {
	var user request.UpdateUser

	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "invalid id.Please enter a valid integer id").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "failed to decode json request body").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		return
	}

	resp, err := u.userService.UpdateUser(id, user)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// deleteUserHandler  godoc
// @Summary		      Delete user
// @Description	      delete user by id
// @Tags			  user
// @ID				  delete-user-by-id
// @Produce		      json
// @Security		  JWT
// @Param			  id	path		int	true	"UserID"
// @Success	     	  200	{object}	response.Response
// @Router			  /user/delete/{id} [delete]
func (u UserController) deleteUserHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "invalid id.Please enter a valid integer id").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		return
	}

	resp, err := u.userService.DeleteUser(id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}

// refreshTokenHandler  godoc
// @Summary		        Refresh token
// @Description	        refresh the expired access token
// @Tags			    user
// @ID				    refresh-access-token
// @Accept			    json
// @Produce		        json
// @Security		    JWT
// @Param			    token	 body		request.RefreshRequest	true	"Refresh token"
// @Success		        200		{object}	response.Response
// @Router			    /user/token [post]
func (u UserController) refreshTokenHandler(c *gin.Context) {
	var rfToken request.RefreshRequest

	if err := c.ShouldBindJSON(&rfToken); err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "failed to decode json request body").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		return
	}

	// validate := c.MustGet("validator").(*validator.Validate)
	// if err := validate.Struct(rfToken); err != nil {
	// 	errorResponse := response.Response{
	// 		Status:       http.StatusBadRequest,
	// 		ErrorType:    errorcode.ValidationError,
	// 		ErrorMessage: customizer3.DecryptErrors(err),
	// 	}
	// 	// c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
	// 	c.Set("error", errorResponse)
	// 	return
	// }

	resp, err := u.userService.RefreshToken(rfToken)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}
