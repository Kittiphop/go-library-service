package handler

import (
	"go-library-service/cmd/api/constant"
	"go-library-service/cmd/api/entity"
	"go-library-service/cmd/api/middleware"
	errmap "go-library-service/internal/error_map"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// CreateBook creates a new book
// @Summary Create a new book
// @Description Create a new book
// @Tags management books
// @Accept  json
// @Produce  json
// @Param   book  body      entity.BookCreateRequest  true  "Create book"
// @Success 201 {object} entity.ResponseData
// @Failure 400 {object} entity.ResponseError
// @Failure 500 {object} entity.ResponseError
// @Security BearerAuth
// @Router /management/books [post]
func (h *Handler) CreateBook(c *gin.Context) {
	var req entity.BookCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(errors.Wrap(err, "[Handler.CreateBook]: invalid request"))
		c.JSON(http.StatusBadRequest, entity.ResponseError{Error: "invalid request", Code: http.StatusBadRequest})
		return
	}

	if err := h.deps.Validator.Struct(req); err != nil {
		log.Error(errors.Wrap(err, "[Handler.CreateBook]: invalid request"))
		c.JSON(http.StatusBadRequest, entity.ResponseError{Error: err.Error(), Code: http.StatusBadRequest})
		return
	}

	err := h.deps.Service.CreateBook(req)
	if err != nil {
		log.Error(errors.Wrap(err, "[Handler.CreateBook]: unable to create book"))
		c.JSON(http.StatusInternalServerError, entity.ResponseError{Error: "Unable to create book", Code: http.StatusInternalServerError})
		return
	}

	c.AbortWithStatus(http.StatusCreated)
}

// ListBook lists all books
// @Summary List all books
// @Description Get a list of books with pagination and optional filters
// @Tags books
// @Accept  json
// @Produce  json
// @Param   page      query     int     false  "Page number"
// @Param   size	  query     int     false  "Number of items per page"
// @Param   search    query     string  false  "Search query"
// @Success 200 {object} entity.ResponseData{data=[]entity.BookResponse}
// @Failure 400 {object} entity.ResponseError
// @Failure 500 {object} entity.ResponseError
// @Security BearerAuth
// @Router /books [get]
func (h *Handler) ListBook(c *gin.Context) {
	var req entity.ListBookRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Error(errors.Wrap(err, "[Handler.ListBook]: unable to bind query"))
		c.JSON(http.StatusBadRequest, entity.ResponseError{Error: "Unable to bind query", Code: http.StatusBadRequest})
		return
	}

	if err := h.deps.Validator.Struct(req); err != nil {
		log.Error(errors.Wrap(err, "[Handler.ListBook]: invalid request"))
		c.JSON(http.StatusBadRequest, entity.ResponseError{Error: err.Error(), Code: http.StatusBadRequest})
		return
	}

	books, err := h.deps.Service.ListBook(req)
	if err != nil {
		log.Error(errors.Wrap(err, "[Handler.ListBook]: unable to list books"))
		c.JSON(http.StatusInternalServerError, entity.ResponseError{Error: "Unable to list books", Code: http.StatusInternalServerError})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, entity.ResponseData{Data: books})
}

// GetBookByID gets a book by id
// @Summary Get a book by ID
// @Description Get a book by ID
// @Tags books
// @Accept  json
// @Produce  json
// @Param   id   path      int  true  "Book ID"
// @Success 200 {object} entity.ResponseData{data=entity.BookResponse}
// @Failure 400 {object} entity.ResponseError
// @Failure 404 {object} entity.ResponseError
// @Failure 500 {object} entity.ResponseError
// @Security BearerAuth
// @Router /books/{id} [get]
func (h *Handler) GetBookByID(c *gin.Context) {
	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error(errors.Wrap(err, "[Handler.GetBookByID]: unable to convert book id"))
		c.JSON(http.StatusBadRequest, entity.ResponseError{Error: "invalid request", Code: http.StatusBadRequest})
		return
	}

	book, err := h.deps.Service.GetBookByID(uint(bookID))
	if err != nil {
		log.Error(errors.Wrap(err, "[Handler.GetBookByID]: unable to get book"))
		c.JSON(http.StatusInternalServerError, entity.ResponseError{Error: "Unable to get book", Code: http.StatusInternalServerError})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, entity.ResponseData{Data: book})
}

// UpdateBook updates a book
// @Summary Update a book
// @Description Update a book by ID
// @Tags management books
// @Accept  json
// @Produce  json
// @Param   id   path      int  true  "Book ID"
// @Param   update  body      entity.BookUpdateRequest  true  "Update book"
// @Success 200 {object} entity.ResponseData
// @Failure 400 {object} entity.ResponseError
// @Failure 404 {object} entity.ResponseError
// @Failure 500 {object} entity.ResponseError
// @Security BearerAuth
// @Router /management/books/{id} [put]
func (h *Handler) UpdateBook(c *gin.Context) {
	var req entity.BookUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(errors.Wrap(err, "[Handler.UpdateBook]: invalid request"))
		c.JSON(http.StatusBadRequest, entity.ResponseError{Error: "invalid request", Code: http.StatusBadRequest})
		return
	}

	if err := h.deps.Validator.Struct(req); err != nil {
		log.Error(errors.Wrap(err, "[Handler.UpdateBook]: invalid request"))
		c.JSON(http.StatusBadRequest, entity.ResponseError{Error: err.Error(), Code: http.StatusBadRequest})
		return
	}

	if err := h.deps.Service.UpdateBook(req); err != nil {
		if errors.Is(err, errmap.ErrmapNotFound) {
			c.JSON(http.StatusNotFound, entity.ResponseError{Error: "book not found", Code: http.StatusNotFound})
			return
		}

		if errors.Is(err, errmap.ErrmapInvalidStock) {
			c.JSON(http.StatusConflict, entity.ResponseError{Error: "invalid stock", Code: http.StatusConflict})
			return
		}

		log.Error(errors.Wrap(err, "[Handler.UpdateBook]: unable to update book"))
		c.JSON(http.StatusInternalServerError, entity.ResponseError{Error: "unable to update book", Code: http.StatusInternalServerError})
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

// ListLatestBooks lists latest books
// @Summary List latest books
// @Description Get a list of latest books
// @Tags books
// @Accept  json
// @Produce  json
// @Success 200 {object} entity.ResponseData{data=[]entity.BookResponse}
// @Failure 400 {object} entity.ResponseError
// @Failure 500 {object} entity.ResponseError
// @Security BearerAuth
// @Router /books/latest [get]
func (h *Handler) ListLatestBooks(c *gin.Context) {
	books, err := h.deps.Service.ListLatestBooks()
	if err != nil {
		log.Error(errors.Wrap(err, "[Handler.ListLatestBooks]: unable to list latest books"))
		c.JSON(http.StatusInternalServerError, entity.ResponseError{Error: "unable to list latest books", Code: http.StatusInternalServerError})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, entity.ResponseData{Data: books})
}

// BorrowBook borrows a book
// @Summary Borrow a book
// @Description Borrow a book by ID
// @Tags books
// @Accept   json
// @Produce  json
// @Param   id   path      int  true  "Book ID"
// @Param   borrow  body      entity.BorrowBookRequest  true  "Borrow book"
// @Success 200 {object} entity.ResponseData
// @Failure 400 {object} entity.ResponseError
// @Failure 404 {object} entity.ResponseError
// @Failure 409 {object} entity.ResponseError
// @Failure 500 {object} entity.ResponseError
// @Security BearerAuth
// @Router /books/:id/borrow [post]
func (h *Handler) BorrowBook(c *gin.Context) {
	userID := h.getJWTInfo(c)

	var req entity.BorrowBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(errors.Wrap(err, "[Handler.BorrowBook]: invalid request"))
		c.JSON(http.StatusBadRequest, entity.ResponseError{Error: "invalid request", Code: http.StatusBadRequest})
		return
	}

	if err := h.deps.Validator.Struct(req); err != nil {
		log.Error(errors.Wrap(err, "[Handler.BorrowBook]: invalid request"))
		c.JSON(http.StatusBadRequest, entity.ResponseError{Error: err.Error(), Code: http.StatusBadRequest})
		return
	}

	req.UserID = userID

	book, err := h.deps.Service.BorrowBook(req)
	if err != nil {
		if errors.Is(err, errmap.ErrmapNotFound) {
			c.JSON(http.StatusNotFound, entity.ResponseError{Error: "book not found", Code: http.StatusNotFound})
			return
		}
		if errors.Is(err, errmap.ErrmapInvalidStock) {
			c.JSON(http.StatusNotFound, entity.ResponseError{Error: "out of stock", Code: http.StatusConflict})
			return
		}
		log.Error(errors.Wrap(err, "[Handler.BorrowBook]: unable to borrow book"))
		c.JSON(http.StatusInternalServerError, entity.ResponseError{Error: "unable to borrow book", Code: http.StatusInternalServerError})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, entity.ResponseData{Data: book})
}

// ReturnBook returns a book
// @Summary Return a borrowed book
// @Description Return a borrowed book by ID
// @Tags books
// @Accept  json
// @Produce  json
// @Param   id   path      int  true  "Borrow ID"
// @Param   return  body      entity.ReturnBookRequest  true  "Return book"
// @Success 200 {object} entity.ResponseData
// @Failure 400 {object} entity.ResponseError
// @Failure 404 {object} entity.ResponseError
// @Failure 500 {object} entity.ResponseError
// @Security BearerAuth
// @Router /books/:id/return [post]
func (h *Handler) ReturnBook(c *gin.Context) {
	userID := h.getJWTInfo(c)

	var req entity.ReturnBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(errors.Wrap(err, "[Handler.BorrowBook]: invalid request"))
		c.JSON(http.StatusBadRequest, entity.ResponseError{Error: "invalid request", Code: http.StatusBadRequest})
		return
	}

	if err := h.deps.Validator.Struct(req); err != nil {
		log.Error(errors.Wrap(err, "[Handler.BorrowBook]: invalid request"))
		c.JSON(http.StatusBadRequest, entity.ResponseError{Error: err.Error(), Code: http.StatusBadRequest})
		return
	}

	req.UserID = userID

	if err := h.deps.Service.ReturnBook(req); err != nil {
		if errors.Is(err, errmap.ErrmapNotFound) {
			c.JSON(http.StatusNotFound, entity.ResponseError{Error: "borrow history not found", Code: http.StatusNotFound})
			return
		}

		if errors.Is(err, errmap.ErrmapConflict) {
			c.JSON(http.StatusConflict, entity.ResponseError{Error: "borrow history is conflict", Code: http.StatusConflict})
			return
		}

		log.Error(errors.Wrap(err, "[Handler.ReturnBook]: unable to return book"))
		c.JSON(http.StatusInternalServerError, entity.ResponseError{Error: "unable to return book", Code: http.StatusInternalServerError})
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

// GetBookBorrowHistory gets a book borrow history
// @Summary Get borrow history for a book
// @Description Get the borrow history for book
// @Tags management books
// @Accept  json
// @Produce  json
// @Param   id   path      int  true  "Book ID"
// @Success 200 {object} entity.ResponseData{data=[]entity.BorrowHistoryResponse}
// @Failure 400 {object} entity.ResponseError
// @Failure 404 {object} entity.ResponseError
// @Failure 500 {object} entity.ResponseError
// @Security BearerAuth
// @Router /management/books/{id}/history [get]
func (h *Handler) GetBookBorrowHistory(c *gin.Context) {
	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error(errors.Wrap(err, "[Handler.GetBookBorrowHistory]: unable to convert book id"))
		c.JSON(http.StatusBadRequest, entity.ResponseError{Error: "invalid request", Code: http.StatusBadRequest})
		return
	}

	histories, err := h.deps.Service.GetBookBorrowHistory(uint(bookID))
	if err != nil {
		log.Error(errors.Wrap(err, "[Handler.GetBookBorrowHistory]: unable to get book borrow history"))
		c.JSON(http.StatusInternalServerError, entity.ResponseError{Error: "unable to get book borrow history", Code: http.StatusInternalServerError})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, entity.ResponseData{Data: histories})
}

func RegisterBookRoutes(router *gin.RouterGroup, handler *Handler) {
	bookRoutes := router.Group("/books")

	{
		bookRoutes.Use(middleware.AuthMiddleware())
		bookRoutes.Use(middleware.RoleMiddleware(constant.UserTypeUser, constant.UserTypeStaff))

		bookRoutes.GET("", handler.ListBook)
		bookRoutes.GET("/:id", handler.GetBookByID)
		bookRoutes.POST("/:id/borrow", handler.BorrowBook)
		bookRoutes.POST("/:id/return", handler.ReturnBook)
	}

	managementBookRoutes := router.Group("/management/books")
	{
		managementBookRoutes.Use(middleware.AuthMiddleware())
		managementBookRoutes.Use(middleware.RoleMiddleware(constant.UserTypeStaff))
		
		managementBookRoutes.POST("", handler.CreateBook)
		managementBookRoutes.PUT("/:id", handler.UpdateBook)
		managementBookRoutes.GET("/:id/history", handler.GetBookBorrowHistory)
	}
	
}