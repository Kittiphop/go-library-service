package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"time"

	"go-library-service/cmd/api/constant"
	"go-library-service/cmd/api/entity"
	"go-library-service/cmd/api/handler"
	"go-library-service/cmd/api/handler/mock"
	"go-library-service/cmd/api/middleware"
	errmap "go-library-service/internal/error_map"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Book Handler", func() {
	var (
		h           *handler.Handler
		serviceMock *mock.MockService
		r           *gin.Engine
		ctrl       *gomock.Controller
		testBook    *entity.BookResponse
		testToken   string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		serviceMock = mock.NewMockService(ctrl)
		h = handler.NewHandler(&handler.Dependencies{
			Service:   serviceMock,
			Validator: validator.New(),
		}, &handler.Config{})

		gin.SetMode(gin.TestMode)
		r = gin.Default()
		handler.RegisterBookRoutes(r.Group("/api"), h)

		testBook = &entity.BookResponse{
			ID:     1,
			Title:  "Test Book",
			Author: "Test Author",
			Price:  19.99,
			Stock:  5,
		}

		var err error
		testToken, err = middleware.GenerateToken(uint(1), "ADMIN")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("CreateBook", func() {
		It("should create book successfully", func() {
			reqBody := entity.BookCreateRequest{
				Title:  "Test Book",
				Author: "Test Author",
				Price:  19.99,
				Stock:  5,
			}
			jsonValue, _ := json.Marshal(reqBody)
			req, _ := http.NewRequest(http.MethodPost, "/api/management/books", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+testToken)

			serviceMock.EXPECT().
				CreateBook(gomock.Any()).
				Return(nil)

			w := httptest.NewRecorder()
			c := gin.CreateTestContextOnly(w, r)
			c.Request = req

			h.CreateBook(c)

			Expect(w.Code).To(Equal(http.StatusCreated))
		})

		It("should return error for invalid request", func() {
			invalidPayload := `{"title": 123, "author": "Test Author", "price": 19.99, "stock": 5}`
			req, _ := http.NewRequest(http.MethodPost, "/api/management/books", bytes.NewBuffer([]byte(invalidPayload)))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+testToken)

			w := httptest.NewRecorder()
			c := gin.CreateTestContextOnly(w, r)
			c.Request = req

			h.CreateBook(c)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return error for validation failure", func() {
			reqBody := entity.BookCreateRequest{
				Title:  "", // Empty title should fail validation
				Author: "Test Author",
				Price:  19.99,
				Stock:  5,
			}
			jsonValue, _ := json.Marshal(reqBody)
			req, _ := http.NewRequest(http.MethodPost, "/api/management/books", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+testToken)

			w := httptest.NewRecorder()
			c := gin.CreateTestContextOnly(w, r)
			c.Request = req

			h.CreateBook(c)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return error when service fails", func() {
			reqBody := entity.BookCreateRequest{
				Title:  "Test Book",
				Author: "Test Author",
				Price:  19.99,
				Stock:  5,
			}
			jsonValue, _ := json.Marshal(reqBody)
			req, _ := http.NewRequest(http.MethodPost, "/api/management/books", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+testToken)

			serviceMock.EXPECT().
				CreateBook(gomock.Any()).
				Return(errmap.ErrmapNotFound)

			w := httptest.NewRecorder()
			c := gin.CreateTestContextOnly(w, r)
			c.Request = req

			h.CreateBook(c)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	Context("ListBook", func() {
        It("should list books successfully", func() {
            req, _ := http.NewRequest(http.MethodGet, "/api/books?page=1&size=10", nil)
            req.Header.Set("Authorization", "Bearer "+testToken)

            expectedBooks := []entity.BookResponse{*testBook}
            serviceMock.EXPECT().
                ListBook(gomock.Any()).
                Return(expectedBooks, nil)

            w := httptest.NewRecorder()
            c := gin.CreateTestContextOnly(w, r)
            c.Request = req

            h.ListBook(c)

            Expect(w.Code).To(Equal(http.StatusOK))
            var response map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &response)
            Expect(err).NotTo(HaveOccurred())
            Expect(response["data"]).ToNot(BeNil())
        })

        It("should return error for invalid query params", func() {
            req, _ := http.NewRequest(http.MethodGet, "/api/books?page=0&size=0", nil) // Invalid page and size
            req.Header.Set("Authorization", "Bearer "+testToken)

            w := httptest.NewRecorder()
            c := gin.CreateTestContextOnly(w, r)
            c.Request = req

            h.ListBook(c)

            Expect(w.Code).To(Equal(http.StatusBadRequest))
        })

        It("should return error when service fails", func() {
            req, _ := http.NewRequest(http.MethodGet, "/api/books?page=1&size=10", nil)
            req.Header.Set("Authorization", "Bearer "+testToken)

            serviceMock.EXPECT().
                ListBook(gomock.Any()).
                Return(nil, errors.New("internal server error"))

            w := httptest.NewRecorder()
            c := gin.CreateTestContextOnly(w, r)
            c.Request = req

            h.ListBook(c)

            Expect(w.Code).To(Equal(http.StatusInternalServerError))
        })
    })

    Context("GetBookByID", func() {
        It("should get book by ID successfully", func() {
            req, _ := http.NewRequest(http.MethodGet, "/api/books/1", nil)
            req.Header.Set("Authorization", "Bearer "+testToken)

            serviceMock.EXPECT().
                GetBookByID(uint(1)).
                Return(testBook, nil)

            w := httptest.NewRecorder()
            c := gin.CreateTestContextOnly(w, r)
            c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
            c.Request = req

            h.GetBookByID(c)

            Expect(w.Code).To(Equal(http.StatusOK))
            var response map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &response)
            Expect(err).NotTo(HaveOccurred())
            Expect(response["data"]).ToNot(BeNil())
        })

        It("should return error for invalid book ID", func() {
            req, _ := http.NewRequest(http.MethodGet, "/api/books/invalid", nil)
            req.Header.Set("Authorization", "Bearer "+testToken)

            w := httptest.NewRecorder()
            c := gin.CreateTestContextOnly(w, r)
            c.Params = append(c.Params, gin.Param{Key: "id", Value: "invalid"})
            c.Request = req

            h.GetBookByID(c)

            Expect(w.Code).To(Equal(http.StatusBadRequest))
        })

        It("should return not found when book doesn't exist", func() {
            req, _ := http.NewRequest(http.MethodGet, "/api/books/999", nil)
            req.Header.Set("Authorization", "Bearer "+testToken)

            serviceMock.EXPECT().
                GetBookByID(uint(999)).
                Return(nil, errmap.ErrmapNotFound)

            w := httptest.NewRecorder()
            c := gin.CreateTestContextOnly(w, r)
            c.Params = append(c.Params, gin.Param{Key: "id", Value: "999"})
            c.Request = req

            h.GetBookByID(c)

            Expect(w.Code).To(Equal(http.StatusInternalServerError)) // Note: The handler returns 500 for not found
        })
    })

    Context("UpdateBook", func() {
        It("should update book successfully", func() {
            reqBody := entity.BookUpdateRequest{
                ID:     1,
                Title:  "Updated Book",
                Author: "Updated Author",
                Price:  29.99,
                Stock:  10,
            }
            jsonValue, _ := json.Marshal(reqBody)
            req, _ := http.NewRequest(http.MethodPut, "/api/management/books/1", bytes.NewBuffer(jsonValue))
            req.Header.Set("Content-Type", "application/json")
            req.Header.Set("Authorization", "Bearer "+testToken)

            serviceMock.EXPECT().
                UpdateBook(gomock.Any()).
                Return(nil)

            w := httptest.NewRecorder()
            c := gin.CreateTestContextOnly(w, r)
            c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
            c.Request = req

            h.UpdateBook(c)

            Expect(w.Code).To(Equal(http.StatusOK))
        })

        It("should return error for invalid request", func() {
            invalidPayload := `{"title": 123}`
            req, _ := http.NewRequest(http.MethodPut, "/api/management/books/1", bytes.NewBuffer([]byte(invalidPayload)))
            req.Header.Set("Content-Type", "application/json")
            req.Header.Set("Authorization", "Bearer "+testToken)

            w := httptest.NewRecorder()
            c := gin.CreateTestContextOnly(w, r)
            c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
            c.Request = req

            h.UpdateBook(c)

            Expect(w.Code).To(Equal(http.StatusBadRequest))
        })

        It("should return not found when book doesn't exist", func() {
            reqBody := entity.BookUpdateRequest{
                ID:     999,
                Title:  "Book",
                Author: "Author",
                Price:  9.99,
                Stock:  1,
            }
            jsonValue, _ := json.Marshal(reqBody)
            req, _ := http.NewRequest(http.MethodPut, "/api/management/books/999", bytes.NewBuffer(jsonValue))
            req.Header.Set("Content-Type", "application/json")
            req.Header.Set("Authorization", "Bearer "+testToken)

            serviceMock.EXPECT().
                UpdateBook(gomock.Any()).
                Return(errmap.ErrmapNotFound)

            w := httptest.NewRecorder()
            c := gin.CreateTestContextOnly(w, r)
            c.Params = append(c.Params, gin.Param{Key: "id", Value: "999"})
            c.Request = req

            h.UpdateBook(c)

            Expect(w.Code).To(Equal(http.StatusNotFound))
        })

        It("should return bad request for invalid stock", func() {
            reqBody := entity.BookUpdateRequest{
                ID:     1,
                Title:  "Test Book",
                Author: "Test Author",
                Price:  19.99,
                Stock:  0,
            }
			
            jsonValue, _ := json.Marshal(reqBody)
            req, _ := http.NewRequest(http.MethodPut, "/api/management/books/1", bytes.NewBuffer(jsonValue))
            req.Header.Set("Content-Type", "application/json")
            req.Header.Set("Authorization", "Bearer "+testToken)

            w := httptest.NewRecorder()
            c := gin.CreateTestContextOnly(w, r)
            c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
            c.Request = req

            h.UpdateBook(c)

            Expect(w.Code).To(Equal(http.StatusBadRequest))
        })
    })

	Context("ListLatestBooks", func() {
        It("should list latest books successfully", func() {
            req, _ := http.NewRequest(http.MethodGet, "/api/books/latest", nil)
            req.Header.Set("Authorization", "Bearer "+testToken)

            expectedBooks := []entity.BookResponse{*testBook}
            serviceMock.EXPECT().
                ListLatestBooks().
                Return(expectedBooks, nil)

            w := httptest.NewRecorder()
            c := gin.CreateTestContextOnly(w, r)
            c.Request = req

            h.ListLatestBooks(c)

            Expect(w.Code).To(Equal(http.StatusOK))
            var response map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &response)
            Expect(err).NotTo(HaveOccurred())
        })

        It("should return error when service fails", func() {
            req, _ := http.NewRequest(http.MethodGet, "/api/books/latest", nil)
            req.Header.Set("Authorization", "Bearer "+testToken)

            serviceMock.EXPECT().
                ListLatestBooks().
                Return(nil, errors.New("internal server error"))

            w := httptest.NewRecorder()
            c := gin.CreateTestContextOnly(w, r)
            c.Request = req

            h.ListLatestBooks(c)

            Expect(w.Code).To(Equal(http.StatusInternalServerError))
        })
    })

    Context("BorrowBook", func() {
        It("should borrow book successfully", func() {
            reqBody := entity.BorrowBookRequest{
                UserID: 1,
                BookID:  1,
            }
            jsonValue, _ := json.Marshal(reqBody)
            req, _ := http.NewRequest(http.MethodPost, "/api/books/1/borrow", bytes.NewBuffer(jsonValue))
            req.Header.Set("Content-Type", "application/json")
            req.Header.Set("Authorization", "Bearer "+testToken)

			borrowedAt := time.Now()
			expectedBorrowHistory := entity.BorrowHistory{
				ID:        1,
				BookID:     reqBody.BookID,
				UserID:     reqBody.UserID,
				BorrowedAt: &borrowedAt,
				Status:     constant.BorrowStatusBorrowed,
				CreatedAt:  &borrowedAt,
				UpdatedAt:  &borrowedAt,
			}

            serviceMock.EXPECT().
                BorrowBook(gomock.Any()).
                Return(&expectedBorrowHistory, nil)

            w := httptest.NewRecorder()
            c := gin.CreateTestContextOnly(w, r)
            c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
            c.Request = req
			c.Set("userID", uint(1))


            h.BorrowBook(c)

            Expect(w.Code).To(Equal(http.StatusOK))
        })

        It("should return error for invalid request", func() {
            invalidPayload := `{"bookId": "test"}`
            req, _ := http.NewRequest(http.MethodPost, "/api/books/1/borrow", bytes.NewBuffer([]byte(invalidPayload)))
            req.Header.Set("Content-Type", "application/json")
            req.Header.Set("Authorization", "Bearer "+testToken)

            w := httptest.NewRecorder()
            c := gin.CreateTestContextOnly(w, r)
            c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
            c.Request = req
			c.Set("userID", uint(1))

            h.BorrowBook(c)

            Expect(w.Code).To(Equal(http.StatusBadRequest))
        })

        It("should return not found when book doesn't exist", func() {
            reqBody := entity.BorrowBookRequest{
                UserID: 1,
                BookID:  999,
            }
            jsonValue, _ := json.Marshal(reqBody)
            req, _ := http.NewRequest(http.MethodPost, "/api/books/999/borrow", bytes.NewBuffer(jsonValue))
            req.Header.Set("Content-Type", "application/json")
            req.Header.Set("Authorization", "Bearer "+testToken)

            serviceMock.EXPECT().
                BorrowBook(gomock.Any()).
                Return(nil, errmap.ErrmapNotFound)

            w := httptest.NewRecorder()
            c := gin.CreateTestContextOnly(w, r)
            c.Params = append(c.Params, gin.Param{Key: "id", Value: "999"})
            c.Request = req
			c.Set("userID", uint(1))

            h.BorrowBook(c)

            Expect(w.Code).To(Equal(http.StatusNotFound))
        })

        It("should return conflict when book is not available", func() {
            reqBody := entity.BorrowBookRequest{
                UserID: 1,
                BookID:  1,
            }
            jsonValue, _ := json.Marshal(reqBody)
            req, _ := http.NewRequest(http.MethodPost, "/api/books/1/borrow", bytes.NewBuffer(jsonValue))
            req.Header.Set("Content-Type", "application/json")
            req.Header.Set("Authorization", "Bearer "+testToken)

            serviceMock.EXPECT().
                BorrowBook(gomock.Any()).
                Return(nil, errmap.ErrmapInvalidStock)

            w := httptest.NewRecorder()
            c := gin.CreateTestContextOnly(w, r)
            c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
            c.Request = req
			c.Set("userID", uint(1))

            h.BorrowBook(c)

            Expect(w.Code).To(Equal(http.StatusConflict))
        })
    })

    Context("ReturnBook", func() {
        It("should return book successfully", func() {
            reqBody := entity.ReturnBookRequest{
                BookID:  1,
				HistoryID: 1,
            }
            jsonValue, _ := json.Marshal(reqBody)
            req, _ := http.NewRequest(http.MethodPost, "/api/books/1/return", bytes.NewBuffer(jsonValue))
            req.Header.Set("Content-Type", "application/json")
            req.Header.Set("Authorization", "Bearer "+testToken)

            serviceMock.EXPECT().
                ReturnBook(gomock.Any()).
                Return(nil)

            w := httptest.NewRecorder()
            c := gin.CreateTestContextOnly(w, r)
            c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
            c.Request = req
			c.Set("userID", uint(1))

            h.ReturnBook(c)

            Expect(w.Code).To(Equal(http.StatusOK))
        })

        It("should return error for invalid request", func() {
            invalidPayload := `{"bookId": "test"}`
            req, _ := http.NewRequest(http.MethodPost, "/api/books/1/return", bytes.NewBuffer([]byte(invalidPayload)))
            req.Header.Set("Content-Type", "application/json")
            req.Header.Set("Authorization", "Bearer "+testToken)

            w := httptest.NewRecorder()
            c := gin.CreateTestContextOnly(w, r)
            c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
            c.Request = req
			c.Set("userID", uint(1))

            h.ReturnBook(c)

            Expect(w.Code).To(Equal(http.StatusBadRequest))
        })

        It("should return not found when borrow record doesn't exist", func() {
            reqBody := entity.ReturnBookRequest{
                BookID:  999,
				HistoryID: 999,
            }
            jsonValue, _ := json.Marshal(reqBody)
            req, _ := http.NewRequest(http.MethodPost, "/api/books/999/return", bytes.NewBuffer(jsonValue))
            req.Header.Set("Content-Type", "application/json")
            req.Header.Set("Authorization", "Bearer "+testToken)

            serviceMock.EXPECT().
                ReturnBook(gomock.Any()).
                Return(errmap.ErrmapNotFound)

            w := httptest.NewRecorder()
            c := gin.CreateTestContextOnly(w, r)
            c.Params = append(c.Params, gin.Param{Key: "id", Value: "999"})
            c.Request = req
			c.Set("userID", uint(1))


            h.ReturnBook(c)

            Expect(w.Code).To(Equal(http.StatusNotFound))
        })
    })

    Context("GetBookBorrowHistory", func() {
        It("should get book borrow history successfully", func() {
            req, _ := http.NewRequest(http.MethodGet, "/api/management/books/1/history", nil)
            req.Header.Set("Authorization", "Bearer "+testToken)

            expectedHistory := []entity.BorrowHistoryResponse{
                {
                    ID:        1,
                    BookID:    1,
                    UserID:    1,
                    BorrowedAt: time.Now(),
                },
            }

            serviceMock.EXPECT().
                GetBookBorrowHistory(uint(1)).
                Return(expectedHistory, nil)

            w := httptest.NewRecorder()
            c := gin.CreateTestContextOnly(w, r)
            c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
            c.Request = req

            h.GetBookBorrowHistory(c)

            Expect(w.Code).To(Equal(http.StatusOK))
            var response map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &response)
            Expect(err).NotTo(HaveOccurred())
            Expect(response["data"]).ToNot(BeNil())
        })

        It("should return empty borrow history", func() {
            req, _ := http.NewRequest(http.MethodGet, "/api/management/books/1/history", nil)
            req.Header.Set("Authorization", "Bearer "+testToken)

            expectedHistory := []entity.BorrowHistoryResponse{}

            serviceMock.EXPECT().
                GetBookBorrowHistory(uint(1)).
                Return(expectedHistory, nil)

            w := httptest.NewRecorder()
            c := gin.CreateTestContextOnly(w, r)
            c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
            c.Request = req

            h.GetBookBorrowHistory(c)

            Expect(w.Code).To(Equal(http.StatusOK))
            var response map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &response)
            Expect(err).NotTo(HaveOccurred())
            Expect(response["data"]).ToNot(BeNil())
        })

        It("should return error for invalid book ID", func() {
            req, _ := http.NewRequest(http.MethodGet, "/api/management/books/invalid/history", nil)
            req.Header.Set("Authorization", "Bearer "+testToken)

            w := httptest.NewRecorder()
            c := gin.CreateTestContextOnly(w, r)
            c.Params = append(c.Params, gin.Param{Key: "id", Value: "invalid"})
            c.Request = req

            h.GetBookBorrowHistory(c)

            Expect(w.Code).To(Equal(http.StatusBadRequest))
        })

    })
})