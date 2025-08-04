package service

import (
	"time"

	"go-library-service/cmd/api/constant"
	"go-library-service/cmd/api/entity"
	errmap "go-library-service/internal/error_map"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)


func (s *Service) BorrowBook(req entity.BorrowBookRequest) (*entity.BorrowHistory, error) {
	book, err := s.deps.PostgresRepo.GetBookByID(req.BookID)
	if err != nil {
		if errors.Is(err, errmap.ErrmapNotFound) {
			return nil, errmap.ErrmapNotFound
		}
		return nil, errors.Wrap(err, "[Service.BorrowBook]: unable to get book")
	}

	if book.Stock < 1 {
		return nil, errmap.ErrmapInvalidStock
	}

	borrowedAt := time.Now()

	history := &entity.BorrowHistory{
		BookID:     req.BookID,
		UserID:     req.UserID,
		BorrowedAt: &borrowedAt,
		Status:     constant.BorrowStatusBorrowed,
	}

	history, err = s.deps.PostgresRepo.BorrowBook(history)
	if err != nil {
		return nil, err
	}

	return history, nil
}

func (s *Service) ReturnBook(req entity.ReturnBookRequest) error {
	history, err := s.deps.PostgresRepo.GetBorrowHistoryByID(req.HistoryID)
	if err != nil {
		if errors.Is(err, errmap.ErrmapNotFound) {
			return errmap.ErrmapNotFound
		}
		log.Error(errors.Wrap(err, "[Service.ReturnBook]: unable to get borrow history"))
		return errors.Wrap(err, "[Service.ReturnBook]: unable to get borrow history")
	}

	if history.ReturnedAt != nil {
		log.Error(errors.Wrap(err, "[Service.ReturnBook]: book is already returned"))
		return errmap.ErrmapConflict
	}

	if history.Status != constant.BorrowStatusBorrowed {
		log.Error(errors.Wrap(err, "[Service.ReturnBook]: book is not borrowed"))
		return errmap.ErrmapConflict
	}

	if err := s.deps.PostgresRepo.ReturnBook(req.HistoryID, req.BookID, time.Now()); err != nil {
		log.Error(errors.Wrap(err, "[Service.ReturnBook]: failed to return book"))
		return errors.Wrap(err, "[Service.ReturnBook]: failed to return book")
	}

	return nil
}

func (s *Service) GetBookBorrowHistory(bookID uint) ([]entity.BorrowHistoryResponse, error) {
	histories, err := s.deps.PostgresRepo.GetBorrowHistoryByBookID(bookID)
	if err != nil {
		return nil, errors.Wrap(err, "[Service.GetBookBorrowHistory]: unable to get book borrow history")
	}

	return histories, nil
}
