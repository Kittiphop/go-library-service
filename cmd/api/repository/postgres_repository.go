package repository

import (
	"fmt"
	"go-library-service/cmd/api/entity"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresRepository struct {
	postgres *gorm.DB
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func NewPostgresRepository(config PostgresConfig) (*PostgresRepository, error){
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
	config.Host,
	config.User,
	config.Password,
	config.DBName,
	config.Port)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})

	if err != nil {
		return nil, err
	}

	if err := postgresqlMigration(db); err != nil {
		return nil, errors.Wrap(err, "[NewPostgresRepository]: unable to connect postgres")	
	}
	
	
	return &PostgresRepository{postgres: db}, nil
}

func postgresqlMigration(db *gorm.DB) error {
	err := db.AutoMigrate(
		&entity.User{},
		&entity.Book{},
		&entity.BorrowHistory{},
	)

	return err

}

