package main

import (
	"go-library-service/cmd/api/constant"
	"go-library-service/cmd/api/docs"
	"go-library-service/cmd/api/entity"
	"go-library-service/cmd/api/handler"
	"go-library-service/cmd/api/repository"
	"go-library-service/cmd/api/service"
	"go-library-service/internal/utils"
	"os"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)


func initPostgresRepository() *repository.PostgresRepository {
	config := repository.PostgresConfig{
		Host:     requiredEnv("PG_HOST"),
		Port:     requiredEnv("PG_PORT"),
		User:     requiredEnv("PG_USER"),
		Password: requiredEnv("PG_PASSWORD"),
		DBName:   requiredEnv("PG_NAME"),
	}
	r, err :=  repository.NewPostgresRepository(config)
	if err != nil {
		log.Error("Failed to connect to database: ", err)
		panic(err)
	}

	// Generate staff id for test
	// username: staff and password: staff
	staffUser := entity.User{
		Name: "Staff",
		Username: "staff",
		Password: "$2a$10$.GRLpoToY1ZdUxfjy85Av.r.VJcnwKK9pvH/VxZNQpYIhmoRLil/C",
		Role: constant.UserTypeStaff,
	}

	_, err = r.CreateUser(staffUser)
	if err != nil {
		log.Warn("Failed to create staff or existed staff: ", err)
	}

	return r
}

func initRedisRepository() *repository.RedisRepository {
	config := repository.RedisConfig{
		Host:     requiredEnv("REDIS_HOST"),
		Port:     requiredEnv("REDIS_PORT"),
		Password: requiredEnv("REDIS_PASSWORD"),
	}
	r, err := repository.NewRedisRepository(config)

	if err != nil {
		log.Error("Failed to connect to redis: ", err)
		panic(err)
	}

	return r
}

func initService() *service.Service {
	return service.NewService(
		&service.Dependencies{
			PostgresRepo: initPostgresRepository(),
			BcryptService: utils.NewBcryptService(),
			RedisRepo: initRedisRepository(),
		},
		&service.Config{},
	)
}

func initHandler() *handler.Handler {
	return handler.NewHandler(
		&handler.Dependencies{
			Service: initService(),
			Validator: validator.New(),
		},
		&handler.Config{},
	)
}

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func initSwagger(g *gin.Engine) {
	docs.SwaggerInfo.Title = "Go Library Service"
	docs.SwaggerInfo.Description = "This is a library service"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http"}

	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, 
		ginSwagger.URL("/swagger/doc.json"),
	))
}

func requiredEnv(key string) string {
	env, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("required env %s not set", key)
	}
	return env
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	g := gin.Default()

	h := initHandler()
	initSwagger(g)

	port := os.Getenv("APP_PORT")
	
	handler.InitRoute(g.Group("/api"), h)
	endless.ListenAndServe(":"+port, g)
}
