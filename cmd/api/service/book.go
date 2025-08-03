package service

import (
	"encoding/json"
	"go-library-service/cmd/api/entity"
	errmap "go-library-service/internal/error_map"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	cacheKeyLatestBooks = "latest_books"
)

// CreateBook creates a new book
func (s *Service) CreateBook(req entity.BookCreateRequest) error {
	book := entity.Book{
		Title:  req.Title,
		Author: req.Author,
		Price:  req.Price,
		Stock:  req.Stock,
	}

	if req.Stock < 1 {
		return errmap.ErrmapInvalidStock
	}

	err := s.deps.PostgresRepo.CreateBook(book)
	if err != nil {
		log.Error(errors.Wrap(err, "[Service.CreateBook]: unable to create book"))
		return errors.Wrap(err, "[Service.CreateBook]: unable to create book")
	}

	if err := s.deps.RedisRepo.Delete(cacheKeyLatestBooks); err != nil {
		log.Error(errors.Wrap(err, "[Service.CreateBook]: unable to delete cache"))
	}

	return nil
}

// GetBookByID get a book by id
func (s *Service) GetBookByID(bookID uint) (*entity.BookResponse, error) {
	book, err := s.deps.PostgresRepo.GetBookByID(bookID)
	if err != nil {
		if errors.Is(err, errmap.ErrmapNotFound) {
			log.Error(errors.Wrap(err, "[Service.GetBookByID]: book not found"))
			return nil, errmap.ErrmapNotFound
		}
		log.Error(errors.Wrap(err, "[Service.GetBookByID]: unable to get book"))
		return nil, errors.Wrap(err, "[Service.GetBookByID]: unable to get book")
	}

	return book, nil
}

// UpdateBook updates a book
func (s *Service) UpdateBook(req entity.BookUpdateRequest) error {
	_, err := s.deps.PostgresRepo.GetBookByID(req.ID)
	if err != nil {
		if errors.Is(err, errmap.ErrmapNotFound) {
			log.Error(errors.Wrap(err, "[Service.UpdateBook]: book not found"))
			return errmap.ErrmapNotFound
		}
		log.Error(errors.Wrap(err, "[Service.UpdateBook]: unable to get book"))
		return errors.Wrap(err, "[Service.UpdateBook]: unable to get book")
	}
	if req.Stock < 1 {
		return errmap.ErrmapInvalidStock
	}

	book := entity.Book{
		ID:     req.ID,
		Title:  req.Title,
		Author: req.Author,
		Price:  req.Price,
		Stock:  req.Stock,
	}

	if err := s.deps.PostgresRepo.UpdateBook(book);err != nil {
		log.Error(errors.Wrap(err, "[Service.UpdateBook]: unable to update book"))
		return errors.Wrap(err, "[Service.UpdateBook]: unable to update book")
	}

	if err := s.deps.RedisRepo.Delete(cacheKeyLatestBooks); err != nil {
		log.Error(errors.Wrap(err, "[Service.UpdateBook]: unable to delete cache"))
	}

	return nil
}

// ListBook lists books with pagination
func (s *Service) ListBook(req entity.ListBookRequest) ([]entity.BookResponse, error) {
	books, err := s.deps.PostgresRepo.ListBook(req)
	if err != nil {
		log.Error(errors.Wrap(err, "[Service.ListBook]: unable to get books"))
		return nil, errors.Wrap(err, "[Service.ListBook]: unable to get books")
	}

	return books, nil
}

// ListLatestBooks lists latest books
func (s *Service) ListLatestBooks() ([]entity.BookResponse, error) {
    cachedBooks, err := s.deps.RedisRepo.Get(cacheKeyLatestBooks)
    if err == nil && cachedBooks != "" {
        var books []entity.BookResponse
        if err := json.Unmarshal([]byte(cachedBooks), &books); err == nil {
            return books, nil
        }
    }
	
	books, err := s.deps.PostgresRepo.ListLatestBooks()
	if err != nil {
		log.Error(errors.Wrap(err, "[Service.ListLatestBooks]: unable to list latest book"))
		return nil, errors.Wrap(err, "[Service.ListLatestBooks]: unable to list latest book")
	}

	   
	booksJSON, err := json.Marshal(books)
	if err == nil {
		expiration := uint(60 * 60) // 1 hour
		if err := s.deps.RedisRepo.Set(cacheKeyLatestBooks, string(booksJSON), expiration); err != nil {
			log.Error(errors.Wrap(err, "[Service.ListLatestBooks]: unable to cache latest books"))
		}
	}

	return books, nil
}