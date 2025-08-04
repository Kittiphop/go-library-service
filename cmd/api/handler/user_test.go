package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

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

var _ = Describe("User Handler", func() {
	var (
		h         	*handler.Handler
		serviceMock *mock.MockService
		r           *gin.Engine
		ctrl        *gomock.Controller
		testUser *entity.UserResponse
		testToken string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		serviceMock = mock.NewMockService(ctrl)
		h = handler.NewHandler(&handler.Dependencies{
			Service: serviceMock,
			Validator: validator.New(),
		}, &handler.Config{})

		gin.SetMode(gin.TestMode)
		r = gin.Default()
		handler.RegisterUserRoutes(r.Group("/api"), h)

		testUser = &entity.UserResponse{
			ID:       1,
			Name:     "Test User",
		}

		var err error
		testToken, err = middleware.GenerateToken(uint(1), "USER")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("CreateUser", func() {
		It("should create user successfully", func() {
			reqBody := entity.UserCreateRequest{
				Name:     "Test User",
				Username: "testuser",
				Password: "testpass123",
			}
			jsonValue, _ := json.Marshal(reqBody)
			req, _ := http.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c:= gin.CreateTestContextOnly(w, r)
			c.Request = req

			expectedUserID := uint(1)
			serviceMock.EXPECT().
				CreateUser(gomock.Any()).
				Return(&expectedUserID, nil)

			h.CreateUser(c)

			Expect(w.Code).To(Equal(http.StatusOK))
		})

		It("should return error for invalid request", func() {
			invalidPayload := `
			{
				"name": 123,
				"username": "testuser",
				"password": "testpass123"
			}
			`
			req, _ := http.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer([]byte(invalidPayload)))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c:= gin.CreateTestContextOnly(w, r)
			c.Request = req
			h.CreateUser(c)


			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return conflict for duplicate username", func() {
			reqBody := entity.UserCreateRequest{
				Name:     "Test User",
				Username: "existinguser",
				Password: "testpass123",
			}
			jsonValue, _ := json.Marshal(reqBody)
			req, _ := http.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")

			serviceMock.EXPECT().
				CreateUser(gomock.Any()).
				Return(nil, errmap.ErrmapConflict)

			w := httptest.NewRecorder()
			c:= gin.CreateTestContextOnly(w, r)
			c.Request = req
			h.CreateUser(c)

			Expect(w.Code).To(Equal(http.StatusConflict))
		})
	})

	Context("Login", func() {
		It("should login user successfully", func() {
			reqBody := entity.UserLoginRequest{
				Username: "testuser",
				Password: "testpass123",
			}
			jsonValue, _ := json.Marshal(reqBody)
			req, _ := http.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")

			serviceMock.EXPECT().
				LoginUser(reqBody.Username, reqBody.Password).
				Return(&entity.User{
					ID:       1,
					Name: "Test User",
					Username: "testuser",
					Password: "testpass123",
				}, nil)

			w := httptest.NewRecorder()
			c:= gin.CreateTestContextOnly(w, r)
			c.Request = req
			h.Login(c)

			Expect(w.Code).To(Equal(http.StatusOK))
			
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			Expect(err).NotTo(HaveOccurred())
			Expect(response["data"]).To(HaveKey("token"))
		})

		It("should return unauthorized for invalid password", func() {
			reqBody := entity.UserLoginRequest{
				Username: "wronguser",
				Password: "wrongpass",
			}
			jsonValue, _ := json.Marshal(reqBody)
			req, _ := http.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")

			serviceMock.EXPECT().
				LoginUser(reqBody.Username, reqBody.Password).
				Return(nil, errmap.ErrmapInvalidPassword)

			w := httptest.NewRecorder()
			c:= gin.CreateTestContextOnly(w, r)
			c.Request = req
			h.Login(c)

			Expect(w.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should return user not found", func() {
			reqBody := entity.UserLoginRequest{
				Username: "nonexistent",
				Password: "password",
			}
			jsonValue, _ := json.Marshal(reqBody)
			req, _ := http.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")

			serviceMock.EXPECT().
				LoginUser(reqBody.Username, reqBody.Password).
				Return(nil, errmap.ErrmapNotFound)

			w := httptest.NewRecorder()
			c:= gin.CreateTestContextOnly(w, r)
			c.Request = req
			h.Login(c)

			Expect(w.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("GetUser", func() {
		It("should get user info successfully", func() {
			req, _ := http.NewRequest(http.MethodGet, "/api/users/info", nil)
			req.Header.Set("Authorization", "Bearer "+testToken)

			serviceMock.EXPECT().
				GetUserByID(uint(1)).
				Return(testUser, nil)

			w := httptest.NewRecorder()
			c:= gin.CreateTestContextOnly(w, r)
			c.Request = req
			c.Set("userID", uint(1))
			c.Set("role", "USER")

			h.GetUser(c)

			Expect(w.Code).To(Equal(http.StatusOK))
			
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			Expect(err).NotTo(HaveOccurred())
			Expect(response["data"]).To(HaveKey("id"))
			Expect(response["data"]).To(HaveKey("name"))
			Expect(response["data"]).To(HaveKey("createdAt"))
			Expect(response["data"]).To(HaveKey("updatedAt"))
		})

		It("should return error when user not found", func() {
			req, _ := http.NewRequest(http.MethodGet, "/api/users/info", nil)
			req.Header.Set("Authorization", "Bearer "+testToken)

			serviceMock.EXPECT().
				GetUserByID(uint(1)).
				Return(nil, errmap.ErrmapNotFound)

			w := httptest.NewRecorder()
			c:= gin.CreateTestContextOnly(w, r)
			c.Request = req
			c.Set("userID", uint(1))
			c.Set("role", "USER")
			h.GetUser(c)

			Expect(w.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("UpdateUser", func() {
		It("should update user successfully", func() {
			reqBody := entity.UserUpdateRequest{
				Name: "Updated User",
			}
			jsonValue, _ := json.Marshal(reqBody)
			req, _ := http.NewRequest(http.MethodPut, "/api/users/info", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+testToken)

			serviceMock.EXPECT().
				UpdateUser(gomock.Any()).
				Return(nil)

			w := httptest.NewRecorder()
			c:= gin.CreateTestContextOnly(w, r)
			c.Request = req
			c.Set("userID", uint(1))
			c.Set("role", "USER")

			h.UpdateUser(c)

			Expect(w.Code).To(Equal(http.StatusOK))
		})

		It("should return error for invalid request", func() {
			req, _ := http.NewRequest(http.MethodPut, "/api/users/info", bytes.NewBuffer([]byte("{invalid json}")))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+testToken)

			w := httptest.NewRecorder()
			c:= gin.CreateTestContextOnly(w, r)
			c.Request = req
			c.Set("userID", uint(1))
			c.Set("role", "USER")

			h.UpdateUser(c)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("DeleteUser", func() {
		It("should delete user successfully", func() {
			req, _ := http.NewRequest(http.MethodDelete, "/api/management/users/1", nil)
			req.Header.Set("Authorization", "Bearer "+testToken)

			serviceMock.EXPECT().
				DeleteUser(uint(1)).
				Return(nil)

			w := httptest.NewRecorder()
			c:= gin.CreateTestContextOnly(w, r)
			c.Request = req
			c.Set("userID", uint(2))
			c.Set("role", "ADMIN")
			c.Params = []gin.Param{
				{Key: "id", Value: "1"},
			}

			h.DeleteUser(c)

			Expect(w.Code).To(Equal(http.StatusOK))
		})

		It("should return user not found", func() {
			req, _ := http.NewRequest(http.MethodDelete, "/api/management/users/999", nil)
			req.Header.Set("Authorization", "Bearer "+testToken)

			serviceMock.EXPECT().
				DeleteUser(uint(999)).
				Return(errmap.ErrmapNotFound)

			w := httptest.NewRecorder()
			c:= gin.CreateTestContextOnly(w, r)
			c.Request = req
			c.Set("userID", uint(2))
			c.Set("role", "ADMIN")
			c.Params = []gin.Param{
				{Key: "id", Value: "999"},
			}

			h.DeleteUser(c)

			Expect(w.Code).To(Equal(http.StatusNotFound))
		})
	})
})

