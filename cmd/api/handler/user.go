package handler

import (
	"net/http"
	"strconv"

	"go-library-service/cmd/api/constant"
	"go-library-service/cmd/api/entity"
	"go-library-service/cmd/api/middleware"
	errmap "go-library-service/internal/error_map"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// Login handles user authentication
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   credentials  body      entity.UserLoginRequest  true  "Login credentials"
// @Success 200 {object} entity.ResponseData{data=entity.LoginResponse}
// @Failure 400 {object} entity.ResponseError
// @Failure 401 {object} entity.ResponseError
// @Failure 500 {object} entity.ResponseError
// @Router /login [post]
func (h *Handler) Login(c *gin.Context) {
	var req entity.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, entity.ResponseError{Error: "Invalid request", Code: http.StatusBadRequest})
		return
	}

	user, err := h.deps.Service.LoginUser(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, errmap.ErrmapNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, entity.ResponseError{Error: "user not found", Code: http.StatusNotFound})
			return
		}
		if errors.Is(err, errmap.ErrmapInvalidPassword) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, entity.ResponseError{Error: "invalid password", Code: http.StatusUnauthorized})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, entity.ResponseError{Error: "unable to login", Code: http.StatusInternalServerError})
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, entity.ResponseError{Error: "Failed to generate token", Code: http.StatusInternalServerError})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, entity.ResponseData{
		Data: entity.LoginResponse{
			Token: token,
		},
	})
}

// CreateUser handles user creation
// @Summary Create a new user
// @Description Create a new user with the input payload
// @Tags users
// @Accept  json
// @Produce  json
// @Param   user  body      entity.UserCreateRequest  true  "Create user"
// @Success 201 {object} entity.ResponseData{data=entity.UserResponse}
// @Failure 400 {object} entity.ResponseError
// @Failure 409 {object} entity.ResponseError
// @Failure 500 {object} entity.ResponseError
// @Router /register [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var req entity.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, entity.ResponseError{Error: "invalid request", Code: http.StatusBadRequest})
		return
	}

	if err := h.deps.Validator.Struct(req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, entity.ResponseError{Error: err.Error(), Code: http.StatusBadRequest})
		return
	}


	userID, err := h.deps.Service.CreateUser(req)
	if err != nil {
		if errors.Is(err, errmap.ErrmapConflict) {
			c.AbortWithStatusJSON(http.StatusConflict, entity.ResponseError{Error: "username already exists", Code: http.StatusConflict})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, entity.ResponseError{Error: "unable to create user", Code: http.StatusInternalServerError})
		return
	}

	if userID == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, entity.ResponseError{Error: "unable to create user", Code: http.StatusInternalServerError})
		return
	}

	token, err := middleware.GenerateToken(*userID, constant.UserTypeUser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, entity.ResponseError{Error: "failed to generate token", Code: http.StatusInternalServerError})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, entity.ResponseData{Data: entity.LoginResponse{Token: token}})
}

// GetUser handles getting a user by ID
// @Summary Get a user by ID
// @Description Get a user by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} entity.ResponseData{data=entity.UserResponse}
// @Failure 404 {object} entity.ResponseError
// @Failure 500 {object} entity.ResponseError
// @Router /users/info [get]
func (h *Handler) GetUser(c *gin.Context) {
	userID := h.getJWTInfo(c)

	user, err := h.deps.Service.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, errmap.ErrmapNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, entity.ResponseError{Error: "user not found", Code: http.StatusNotFound})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, entity.ResponseError{Error: "unable to get user", Code: http.StatusInternalServerError})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, entity.ResponseData{Data: user})
}

// UpdateUser handles updating a user
// @Summary Update a user
// @Description Update a user with the input payload
// @Tags users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   user  body      entity.UserUpdateRequest  true  "Update user"
// @Success 200 {object} entity.ResponseData{data=entity.UserResponse}
// @Failure 400 {object} entity.ResponseError
// @Failure 404 {object} entity.ResponseError
// @Failure 409 {object} entity.ResponseError
// @Failure 500 {object} entity.ResponseError
// @Router /users/info [put]
func (h *Handler) UpdateUser(c *gin.Context) {
	userID := h.getJWTInfo(c)

	var user entity.UserUpdateRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, entity.ResponseError{Error: "invalid request", Code: http.StatusBadRequest})
		return
	}

	if err := h.deps.Validator.Struct(user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, entity.ResponseError{Error: err.Error(), Code: http.StatusBadRequest})
		return
	}

	user.ID = uint(userID)

	if err := h.deps.Service.UpdateUser(user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, entity.ResponseError{Error: "unable to update user", Code: http.StatusInternalServerError})
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

// DeleteUser handles deleting a user
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags management users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   id   path      int  true  "User ID"
// @Success 200 {object} entity.ResponseData
// @Failure 404 {object} entity.ResponseError
// @Failure 500 {object} entity.ResponseError
// @Router /management/users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, entity.ResponseError{Error: "invalid user id", Code: http.StatusBadRequest})
		return
	}

	if err := h.deps.Service.DeleteUser(uint(userID)); err != nil {
		if errors.Is(err, errmap.ErrmapNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, entity.ResponseError{Error: "user not found", Code: http.StatusNotFound})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, entity.ResponseError{Error: "unable to delete user", Code: http.StatusInternalServerError})
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

// RegisterUserRoutes registers user routes
func RegisterUserRoutes(router *gin.RouterGroup, handler *Handler) {

	userRoutes := router.Group("/users")
	userRoutes.Use(middleware.AuthMiddleware())
	userRoutes.Use(middleware.RoleMiddleware(constant.UserTypeUser, constant.UserTypeStaff))
	{
		userRoutes.GET("/info", handler.GetUser)
		userRoutes.PUT("/info", handler.UpdateUser)
	}

	managementUserRoutes := router.Group("/management/users")
	managementUserRoutes.Use(middleware.AuthMiddleware())
	managementUserRoutes.Use(middleware.RoleMiddleware(constant.UserTypeStaff))
	{
		managementUserRoutes.DELETE("/:id", handler.DeleteUser)
	}
}
