package repository

import (
	"time"

	"go-library-service/cmd/api/constant"
	"go-library-service/cmd/api/entity"
	errmap "go-library-service/internal/error_map"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *PostgresRepository) BorrowBook(history *entity.BorrowHistory) (*entity.BorrowHistory, error) {
	tx := r.postgres.Begin()
	defer func() {
		if r := recover(); r != nil {
		  tx.Rollback()
		}
	  }()

	if tx.Error != nil {
		return nil, tx.Error
	}

	var book entity.Book
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Table("books").First(&book, history.BookID).Error; err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "[PostgresRepository.BorrowBook]: unable to get book")
	}

	if book.Stock < 1 {
		tx.Rollback()
		return nil, errmap.ErrmapInvalidStock
	}

	if err := tx.Table("books").Where("id = ?", history.BookID).Update("stock", gorm.Expr("stock - ?", 1)).Error; err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "[PostgresRepository.BorrowBook]: unable to update book stock")
	}

	if err := tx.Table("borrow_histories").Create(history).Error; err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "[PostgresRepository.BorrowBook]: unable to create borrow history")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "[PostgresRepository.BorrowBook]: unable to commit transaction")
	}

	return history, nil
}

func (r *PostgresRepository) ReturnBook(historyID, BookID uint, returnedAt time.Time) error {
	tx := r.postgres.Begin()
	defer func() {
		if r := recover(); r != nil {
		  tx.Rollback()
		}
	  }()
	  
	if tx.Error != nil {
		return tx.Error
	}

	var history entity.BorrowHistory
	if err := tx.Table("borrow_histories").First(&history, historyID).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "[PostgresRepository.ReturnBook]: unable to get borrow history")
	}

	if history.BookID != BookID {
		tx.Rollback()
		return errors.New("[PostgresRepository.ReturnBook]: book id does not match")
	}

	if err := tx.Table("borrow_histories").
		Where("id = ?", historyID).
		Updates(map[string]interface{}{
			"returned_at": returnedAt,
			"status":      constant.BorrowStatusReturned,
		}).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "[PostgresRepository.ReturnBook]: unable to update borrow history")
	}

	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Table("books").
		Where("id = ?", history.BookID).
		Update("stock", gorm.Expr("stock + ?", 1)).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "[PostgresRepository.ReturnBook]: unable to update book stock")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "[PostgresRepository.ReturnBook]: unable to commit transaction")
	}

	return nil
}

func (r *PostgresRepository) GetBorrowHistoryByBookID(bookID uint) ([]entity.BorrowHistory, error) {
	var history []entity.BorrowHistory
	err := r.postgres.Table("borrow_histories").Where("book_id = ?", bookID).Order("created_at DESC").Find(&history).Error
	if err != nil {
		log.Error(errors.Wrap(err, "[PostgresRepository.GetBorrowHistoryByBookID]: unable to get borrow history"))
		return nil, errors.Wrap(err, "[PostgresRepository.GetBorrowHistoryByBookID]: unable to get borrow history")
	}
	return history, nil
}

func (r *PostgresRepository) GetBorrowHistoryByID(id uint) (*entity.BorrowHistory, error) {
	var history entity.BorrowHistory
	err := r.postgres.Table("borrow_histories").First(&history, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errmap.ErrmapNotFound
		}
		return nil, errors.Wrap(err, "[PostgresRepository.GetBorrowHistoryByID]: unable to get borrow history")
	}
	return &history, nil
}
