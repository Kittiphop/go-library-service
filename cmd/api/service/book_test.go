package service_test

import (
	"encoding/json"

	"go-library-service/cmd/api/entity"
	service "go-library-service/cmd/api/service"
	"go-library-service/cmd/api/service/mock"
	errmap "go-library-service/internal/error_map"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Book Service", func() {
	var (
		ctrl *gomock.Controller
		s     *service.Service
		postgresMock  *mock.MockPostgresRepository
		redisMock *mock.MockRedisRepository

		cacheKeyLatestBooks = "latest_books"
		bookCreateRequest   entity.BookCreateRequest
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		postgresMock = mock.NewMockPostgresRepository(ctrl)
		redisMock = mock.NewMockRedisRepository(ctrl)
		s = service.NewService(&service.Dependencies{
			PostgresRepo: postgresMock,
			RedisRepo:    redisMock,
		}, &service.Config{})

		bookCreateRequest = entity.BookCreateRequest{
			Title:  "Test Book",
			Author: "Test Author",
			Price:  29.99,
			Stock:  10,
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("CreateBook", func() {
		It("should create a book successfully", func() {
			postgresMock.EXPECT().CreateBook(gomock.Any()).Return(nil)
			redisMock.EXPECT().Delete(cacheKeyLatestBooks).Return(nil)

			err := s.CreateBook(bookCreateRequest)
			Expect(err).To(BeNil())
		})
	})

	Context("GetBookByID", func() {
		It("should return a book by ID", func() {
			bookID := uint(1)
			expectedBook := &entity.BookResponse{
				ID:     bookID,
				Title:  "Test Book",
				Author: "Test Author",
				Price:  29.99,
				Stock:  10,
			}

			postgresMock.EXPECT().GetBookByID(bookID).Return(expectedBook, nil)

			book, err := s.GetBookByID(bookID)
			Expect(err).To(BeNil())
			Expect(book).To(Equal(expectedBook))
		})

		It("should return error when book not found", func() {
			bookID := uint(999)
			postgresMock.EXPECT().GetBookByID(bookID).Return(nil, errmap.ErrmapNotFound)

			book, err := s.GetBookByID(bookID)
			Expect(err).To(Equal(errmap.ErrmapNotFound))
			Expect(book).To(BeNil())
		})
	})

	Context("UpdateBook", func(){
		It("should update a book successfully", func() {
			bookID := uint(1)
			expectedBook := &entity.BookResponse{
				ID:     bookID,
				Title:  "Test Book",
				Author: "Test Author",
				Price:  29.99,
				Stock:  10,
			}

			postgresMock.EXPECT().GetBookByID(bookID).Return(expectedBook, nil)
			postgresMock.EXPECT().UpdateBook(gomock.Any()).Return(nil)
			redisMock.EXPECT().Delete(cacheKeyLatestBooks).Return(nil)

			err := s.UpdateBook(entity.BookUpdateRequest{
				ID:     bookID,
				Title:  "Test Book",
				Author: "Test Author",
				Price:  29.99,
				Stock:  10,
			})
			Expect(err).To(BeNil())
		})

		It("should return error when book not found", func() {
			bookID := uint(999)
			postgresMock.EXPECT().GetBookByID(bookID).Return(nil, errmap.ErrmapNotFound)

			err := s.UpdateBook(entity.BookUpdateRequest{
				ID:     bookID,
				Title:  "Test Book",
				Author: "Test Author",
				Price:  29.99,
				Stock:  10,
			})
			Expect(err).To(Equal(errmap.ErrmapNotFound))
		})

		It("should return error when invalid stock", func() {
			bookID := uint(1)
			postgresMock.EXPECT().GetBookByID(bookID).Return(&entity.BookResponse{
				ID:     bookID,
				Title:  "Test Book",
				Author: "Test Author",
				Price:  29.99,
				Stock:  10,
			}, nil)

			err := s.UpdateBook(entity.BookUpdateRequest{
				ID:     bookID,
				Title:  "Test Book",
				Author: "Test Author",
				Price:  29.99,
				Stock:  0,
			})
			Expect(err).To(Equal(errmap.ErrmapInvalidStock))
		})
	})

	Context("ListBook", func() {
		It("should return list books successfully", func() {
			expectedBooks := []entity.BookResponse{{
				ID:     1,
				Title:  "Book",
				Author: "Author",
				Price:  19.99,
				Stock:  5,
			}}

			postgresMock.EXPECT().ListBook(gomock.Any()).Return(expectedBooks, nil)

			books, err := s.ListBook(entity.ListBookRequest{
				Page:  1,
				Size:  10,
			})
			Expect(err).To(BeNil())
			Expect(books).To(Equal(expectedBooks))
		})

		It("should return empty when no records found", func() {
			postgresMock.EXPECT().ListBook(gomock.Any()).Return([]entity.BookResponse{}, nil)

			books, err := s.ListBook(entity.ListBookRequest{
				Page:  1,
				Size:  10,
			})
			Expect(err).To(BeNil())
			Expect(books).To(BeEmpty())
		})
	})

	Context("ListLatestBooks", func() {
		It("should return list latest books successfully", func() {
			expectedBooks := []entity.BookResponse{{
				ID:     1,
				Title:  "Cached Book",
				Author: "Cache Author",
				Price:  19.99,
				Stock:  5,
			}}

			cachedData, _ := json.Marshal(expectedBooks)
			redisMock.EXPECT().Get(cacheKeyLatestBooks).Return(string(cachedData), nil)

			books, err := s.ListLatestBooks()
			Expect(err).To(BeNil())
			Expect(books).To(Equal(expectedBooks))
		})

		It("should fetch from database and cache when cache is empty", func() {
			expectedBooks := []entity.BookResponse{{
				ID:     2,
				Title:  "New Book",
				Author: "New Author",
				Price:  24.99,
				Stock:  8,
			}}

			redisMock.EXPECT().Get(cacheKeyLatestBooks).Return("", nil)
			postgresMock.EXPECT().ListLatestBooks().Return(expectedBooks, nil)
			redisMock.EXPECT().Set(cacheKeyLatestBooks, gomock.Any(), uint(3600)).Return(nil)

			books, err := s.ListLatestBooks()
			Expect(err).To(BeNil())
			Expect(books).To(Equal(expectedBooks))
		})

		It("should handle cache unmarshal error", func() {
			redisMock.EXPECT().Get(cacheKeyLatestBooks).Return("invalid-json", nil)
			postgresMock.EXPECT().ListLatestBooks().Return([]entity.BookResponse{}, nil)
			redisMock.EXPECT().Set(cacheKeyLatestBooks, gomock.Any(), uint(3600)).Return(nil)

			_, err := s.ListLatestBooks()
			Expect(err).To(BeNil())
		})
	})
})