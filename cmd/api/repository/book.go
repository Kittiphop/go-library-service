package repository

import (
	"go-library-service/cmd/api/entity"
	errmap "go-library-service/internal/error_map"

	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreateBook creates a new book
func (r *PostgresRepository) CreateBook(book entity.Book) error {
	err := r.postgres.Table("books").Create(&book).Error
	if err != nil {
		return errors.Wrap(err, "[PostgresRepository.CreateBook]: unable to create book")
	}
	return nil
}

// UpdateBook updates a book 
func (r *PostgresRepository) UpdateBook(book entity.Book) error {
    tx := r.postgres.Begin()
	defer func() {
		if r := recover(); r != nil {
		  tx.Rollback()
		}
	  }()
	  
    if tx.Error != nil {
        return errors.Wrap(tx.Error, "[PostgresRepository.UpdateBook]: unable to begin transaction")
    }

    err := tx.Table("books").Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", book.ID).First(&entity.Book{}).Error
    if err != nil {
        tx.Rollback()
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return errmap.ErrmapNotFound
        }
        return errors.Wrap(err, "[PostgresRepository.UpdateBook]: unable to find book")
    }

    err = tx.Table("books").Where("id = ?", book.ID).Updates(book).Error
    if err != nil {
        tx.Rollback()
        return errors.Wrap(err, "[PostgresRepository.UpdateBook]: unable to update book")
    }

    if err := tx.Commit().Error; err != nil {
        tx.Rollback()
        return errors.Wrap(err, "[PostgresRepository.UpdateBook]: unable to commit transaction")
    }

    return nil
}

// GetBookByID retrieves a book by ID
func (r *PostgresRepository) GetBookByID(bookID uint) (*entity.BookResponse, error) {
	var book entity.BookResponse
	err := r.postgres.Table("books").First(&book, bookID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errmap.ErrmapNotFound
		}
		return nil, errors.Wrap(err, "[PostgresRepository.GetBookByID]: unable to get book")
	}
	return &book, nil
}

// ListBook lists books with pagination
func (r *PostgresRepository) ListBook(req entity.ListBookRequest) ([]entity.BookResponse, error) {
	var books []entity.BookResponse
	query := r.postgres.Table("books")

	if req.Search != nil {
		query = query.Where("title ILIKE ?", "%"+*req.Search+"%")
	}

	err := query.Debug().Offset((req.Page - 1) * req.Size).
		Limit(req.Size).
		Find(&books).Error
	if err != nil {
		return nil, errors.Wrap(err, "[PostgresRepository.ListBook]: unable to get books")
	}
	return books, nil
}

// ListLatestBooks lists latest books
func (r *PostgresRepository) ListLatestBooks() ([]entity.BookResponse, error) {
	var books []entity.BookResponse
	err := r.postgres.Table("books").Order("created_at DESC").Limit(5).Find(&books).Error
	if err != nil {
		return nil, errors.Wrap(err, "[PostgresRepository.ListLatestBooks]: unable to get latest books")
	}
	return books, nil
}