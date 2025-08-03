package service_test

import (
	"testing"

	"go-library-service/cmd/api/entity"
	"go-library-service/cmd/api/service"
	"go-library-service/cmd/api/service/mock"
	errmap "go-library-service/internal/error_map"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

var _ = Describe("User Service", func() {
	var (
		ctrl *gomock.Controller
		s     *service.Service
		postgresMock    *mock.MockPostgresRepository
		redisMock   *mock.MockRedisRepository
		bcryptMock  *mock.MockBcryptService
		// sampleUser  *entity.User
		// sampleReq   entity.UserCreateRequest
		// sampleUserResponse *entity.UserResponse
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		postgresMock = mock.NewMockPostgresRepository(ctrl)
		redisMock = mock.NewMockRedisRepository(ctrl)
		bcryptMock = mock.NewMockBcryptService(ctrl)

		s = service.NewService(&service.Dependencies{
			PostgresRepo: postgresMock,
			RedisRepo:    redisMock,
			BcryptService: bcryptMock,
		}, &service.Config{})




		

	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("CreateUser", func() {
		sampleReq := entity.UserCreateRequest{
			Name:     "Test User",
			Username: "testuser",
			Password: "hashedpassword",
		}

		sampleUser := &entity.User{
			ID:       1,
			Name:     "Test User",
			Username: "testuser",
			Password: "hashedpassword",
			Role:     "user",
		}

		It("should create a new user successfully", func() {
			expectedUserID := uint(1)

			postgresMock.EXPECT().GetUserByUsername(sampleReq.Username).
				Return(nil, errmap.ErrmapNotFound)

			bcryptMock.EXPECT().Hash(sampleReq.Password).
				Return("hashedpassword", nil)

			postgresMock.EXPECT().CreateUser(gomock.Any()).
				Return(&expectedUserID, nil)

			userID, err := s.CreateUser(sampleReq)

			Expect(err).NotTo(HaveOccurred())
			Expect(*userID).To(Equal(expectedUserID))
		})

		It("should return error when username already exists", func() {
			postgresMock.EXPECT().GetUserByUsername(sampleReq.Username).
				Return(sampleUser, nil)

			_, err := s.CreateUser(sampleReq)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("username already exists")))
		})

		It("should return error when password hashing fails", func() {
			postgresMock.EXPECT().GetUserByUsername(sampleReq.Username).
				Return(nil, errmap.ErrmapNotFound)

			bcryptMock.EXPECT().Hash(sampleReq.Password).
				Return("", assert.AnError)

			_, err := s.CreateUser(sampleReq)

			Expect(err).To(HaveOccurred())
		})
	})

	Context("LoginUser", func() {
		sampleReq := entity.UserCreateRequest{
			Name:     "Test User",
			Username: "testuser",
			Password: "hashedpassword",
		}
		
		sampleUser := &entity.User{
			ID:       1,
			Name:     "Test User",
			Username: "testuser",
			Password: "hashedpassword",
			Role:     "user",
		}

		It("should login user successfully", func() {
			postgresMock.EXPECT().GetUserByUsername(sampleReq.Username).
				Return(sampleUser, nil)


			bcryptMock.EXPECT().Compare(sampleUser.Password, sampleReq.Password).
				Return(nil)

			user, err := s.LoginUser(sampleReq.Username, sampleReq.Password)

			Expect(err).NotTo(HaveOccurred())
			Expect(user.Username).To(Equal(sampleUser.Username))
		})

		It("should return error when user not found", func() {
			postgresMock.EXPECT().GetUserByUsername(sampleReq.Username).
				Return(nil, errmap.ErrmapNotFound)

			_, err := s.LoginUser(sampleReq.Username, sampleReq.Password)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("not found")))
		})

		It("should return error for invalid credentials", func() {
			postgresMock.EXPECT().GetUserByUsername(sampleReq.Username).
				Return(sampleUser, nil)

			bcryptMock.EXPECT().Compare(sampleUser.Password, sampleReq.Password).
				Return(errmap.ErrmapInvalidPassword)

			_, err := s.LoginUser(sampleReq.Username, sampleReq.Password)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("invalid password")))
		})
	})

	Context("GetUserByID", func() {
		sampleUser := &entity.User{
			ID:       1,
			Name:     "Test User",
			Username: "testuser",
			Password: "hashedpassword",
			Role:     "user",
		}

		sampleUserResponse := &entity.UserResponse{
			ID:       1,
			Name:     "Test User",
		}
		It("should get user by ID successfully", func() {
			postgresMock.EXPECT().GetUserByID(sampleUser.ID).
				Return(sampleUserResponse, nil)

			user, err := s.GetUserByID(sampleUser.ID)

			Expect(err).NotTo(HaveOccurred())
			Expect(user).To(Equal(sampleUserResponse))
		})

		It("should return error when user not found", func() {
			postgresMock.EXPECT().GetUserByID(sampleUser.ID).
				Return(nil, errmap.ErrmapNotFound)

			_, err := s.GetUserByID(sampleUser.ID)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("not found")))
		})
	})

	Context("UpdateUser", func() {
		expectedUserName := "Test User Updated"
		sampleUser := &entity.User{
			ID:       1,
			Name:     "Test User",
			Username: "testuser",
			Password: "hashedpassword",
			Role:     "user",
		}

		sampleUserResponse := &entity.UserResponse{
			ID:       1,
			Name:     "Test User",
		}
		It("should update user successfully", func() {
			postgresMock.EXPECT().GetUserByID(sampleUser.ID).
				Return(sampleUserResponse, nil)

			postgresMock.EXPECT().UpdateUser(gomock.Any()).
				Return(nil)

			err := s.UpdateUser(entity.UserUpdateRequest{
				ID:   sampleUser.ID,
				Name: expectedUserName,
			})

			Expect(err).NotTo(HaveOccurred())
		})

		It("should return error when user not found", func() {
			postgresMock.EXPECT().GetUserByID(uint(10)).
				Return(nil, errmap.ErrmapNotFound)

			err := s.UpdateUser(entity.UserUpdateRequest{
				ID:   uint(10),
				Name: expectedUserName,
			})

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("not found")))
		})
	})

	Context("DeleteUser", func() {
		sampleUser := &entity.User{
			ID:       1,
			Name:     "Test User",
			Username: "testuser",
			Password: "hashedpassword",
			Role:     "user",
		}

		sampleUserResponse := &entity.UserResponse{
			ID:       1,
			Name:     "Test User",
		}
		It("should delete user successfully", func() {
			postgresMock.EXPECT().GetUserByID(sampleUser.ID).
				Return(sampleUserResponse, nil)

			postgresMock.EXPECT().DeleteUser(sampleUser.ID).
				Return(nil)

			err := s.DeleteUser(sampleUser.ID)

			Expect(err).NotTo(HaveOccurred())
		})

		It("should return error when user not found", func() {
			postgresMock.EXPECT().GetUserByID(uint(10)).
				Return(nil, errmap.ErrmapNotFound)

			err := s.DeleteUser(uint(10))

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("not found")))
		})
	})
})

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Suite")
}
