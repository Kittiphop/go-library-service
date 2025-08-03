package service_test

import (
	"testing"
	"time"

	"go-library-service/cmd/api/constant"
	"go-library-service/cmd/api/entity"
	service "go-library-service/cmd/api/service"
	"go-library-service/cmd/api/service/mock"
	errmap "go-library-service/internal/error_map"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Borrow History Service", func() {
	var (
		ctrl         *gomock.Controller
		s            *service.Service
		postgresMock *mock.MockPostgresRepository
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		postgresMock = mock.NewMockPostgresRepository(ctrl)
		s = service.NewService(&service.Dependencies{
			PostgresRepo: postgresMock,
		}, &service.Config{})
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("BorrowBook", func() {
		It("should borrow a book successfully", func() {
			bookID := uint(1)
			userID := uint(1)
			req := entity.BorrowBookRequest{
				BookID: bookID,
				UserID: userID,
			}

			book := &entity.BookResponse{
				ID:    bookID,
				Stock: 5,
			}

			postgresMock.EXPECT().GetBookByID(bookID).Return(book, nil)
			postgresMock.EXPECT().BorrowBook(gomock.Any()).Return(&entity.BorrowHistory{
				ID:     1,
				BookID: bookID,
				UserID: userID,
				Status: constant.BorrowStatusBorrowed,
			}, nil)

			history, err := s.BorrowBook(req)
			Expect(err).To(BeNil())
			Expect(history.ID).To(Equal(uint(1)))
			Expect(history.BookID).To(Equal(bookID))
			Expect(history.UserID).To(Equal(userID))
			Expect(history.Status).To(Equal(constant.BorrowStatusBorrowed))
		})

		It("should return error when book not found", func() {
			req := entity.BorrowBookRequest{
				BookID: 999,
				UserID: 1,
			}

			postgresMock.EXPECT().GetBookByID(uint(999)).Return(nil, errmap.ErrmapNotFound)

			_, err := s.BorrowBook(req)
			Expect(err).To(Equal(errmap.ErrmapNotFound))
		})

		It("should return error when book is out of stock", func() {
			bookID := uint(1)
			req := entity.BorrowBookRequest{
				BookID: bookID,
				UserID: 1,
			}

			book := &entity.BookResponse{
				ID:    bookID,
				Stock: 0,
			}

			postgresMock.EXPECT().GetBookByID(bookID).Return(book, nil)

			_, err := s.BorrowBook(req)
			Expect(err).To(Equal(errmap.ErrmapInvalidStock))
		})
	})

	Context("ReturnBook", func() {
		It("should return a book successfully", func() {
			historyID := uint(1)
			bookID := uint(1)
			req := entity.ReturnBookRequest{
				HistoryID: historyID,
				BookID:    bookID,
			}

			history := &entity.BorrowHistory{
				ID:     historyID,
				BookID: bookID,
				Status: constant.BorrowStatusBorrowed,
			}

			postgresMock.EXPECT().GetBorrowHistoryByID(historyID).Return(history, nil)
			postgresMock.EXPECT().ReturnBook(historyID, bookID, gomock.Any()).Return(nil)

			err := s.ReturnBook(req)
			Expect(err).To(BeNil())
		})

		It("should return error when history not found", func() {
			historyID := uint(999)
			req := entity.ReturnBookRequest{
				HistoryID: historyID,
				BookID:    1,
			}

			postgresMock.EXPECT().GetBorrowHistoryByID(historyID).Return(nil, errmap.ErrmapNotFound)

			err := s.ReturnBook(req)
			Expect(err).To(Equal(errmap.ErrmapNotFound))
		})

		It("should return error when book is already returned", func() {
			historyID := uint(1)
			bookID := uint(1)
			req := entity.ReturnBookRequest{
				HistoryID: historyID,
				BookID:    bookID,
			}

			returnedAt := time.Now()
			history := &entity.BorrowHistory{
				ID:         historyID,
				BookID:     bookID,
				Status:     constant.BorrowStatusReturned,
				ReturnedAt: &returnedAt,
			}

			postgresMock.EXPECT().GetBorrowHistoryByID(historyID).Return(history, nil)

			err := s.ReturnBook(req)
			Expect(err).To(Equal(errmap.ErrmapConflict))
		})
	})

	Context("GetBookBorrowHistory", func() {
		It("should return borrow history for a book", func() {
			bookID := uint(1)
			expectedHistories := []entity.BorrowHistory{
				{
					ID:         1,
					BookID:     bookID,
					UserID:     1,
					Status:     constant.BorrowStatusReturned,
					BorrowedAt: &time.Time{},
					ReturnedAt: &time.Time{},
				},
				{
					ID:         2,
					BookID:     bookID,
					UserID:     2,
					Status:     constant.BorrowStatusBorrowed,
					BorrowedAt: &time.Time{},
					ReturnedAt: nil,
				},
			}

			postgresMock.EXPECT().GetBorrowHistoryByBookID(bookID).Return(expectedHistories, nil)

			histories, err := s.GetBookBorrowHistory(bookID)
			Expect(err).To(BeNil())
			Expect(histories).To(HaveLen(2))
		})

		It("should return empty history when no records found", func() {
			bookID := uint(999)
			postgresMock.EXPECT().GetBorrowHistoryByBookID(bookID).Return([]entity.BorrowHistory{}, nil)

			histories, err := s.GetBookBorrowHistory(bookID)
			Expect(err).To(BeNil())
			Expect(histories).To(BeEmpty())
		})
	})
})

func TestBorrowHistoryService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Borrow History Service Suite")
}
