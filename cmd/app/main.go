package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/internal/auth"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/internal/factory"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/internal/middlewares"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/internal/modules/user"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	db, err := initDatabase()
	if err != nil {
		log.Fatal(err)
	}

	if err = createSchema(db); err != nil {
		log.Panic(err)
	}

	// Init user module
	userRepository := factory.NewUserRepository(db)
	userService := factory.NewUserService(userRepository)
	userHandler := factory.NewUserHandler(userService)

	// Init list module
	listRepository := factory.NewListRepository(db)
	listService := factory.NewListService(listRepository)
	listHandler := factory.NewListHandler(listService)

	// Init item module
	itemRepository := factory.NewItemRepository(db)
	itemService := factory.NewItemService(itemRepository, listRepository, userRepository)
	itemHandler := factory.NewItemHandler(itemService)

	// Init auth module
	authService := factory.NewAuthService(userRepository)
	authHandler := factory.NewAuthHandler(authService)

	createAdminUser(userRepository)

	router := gin.Default()
	routeGroup := router.Group("/api/v1") // TODO versionize me!!

	// Auth routes
	routeGroup.POST("/authenticate", authHandler.Authenticate)
	routeGroup.POST("/authenticate/sso", authHandler.AuthenticateSSO)

	// User routes
	newPrivateEndpoint(routeGroup, http.MethodPost, "/users", userHandler.Save)
	newPrivateEndpoint(routeGroup, http.MethodGet, "/users/:id", userHandler.Get)
	newPrivateEndpoint(routeGroup, http.MethodPut, "/users/:id", userHandler.Update)

	// List routes
	newPublicEndpoint(routeGroup, http.MethodPost, "/lists", listHandler.Save)
	newPublicEndpoint(routeGroup, http.MethodGet, "/lists/:list_id", listHandler.Get)
	newPrivateEndpoint(routeGroup, http.MethodDelete, "/lists/:list_id", listHandler.Delete)

	// Item routes
	newPrivateEndpoint(routeGroup, http.MethodPost, "/lists/:list_id/items", itemHandler.Save)
	newPrivateEndpoint(routeGroup, http.MethodGet, "/lists/:list_id/items", itemHandler.GetByList)
	newPrivateEndpoint(routeGroup, http.MethodPut, "/lists/:list_id/items/:item_id", itemHandler.Update)
	newPrivateEndpoint(routeGroup, http.MethodDelete, "/lists/:list_id/items/:item_id", itemHandler.Delete)

	auth.InitJWTAuth()

	port := getRunningPort()
	router.Run(fmt.Sprintf("0.0.0.0:%s", port))
}

func initDatabase() (*gorm.DB, error) {
	dsn := os.Getenv("DNCONN_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(
		&models.User{},
		&models.List{},
		&models.Item{},
	)

	return db, nil
}

func createSchema(db *gorm.DB) error {
	path := os.Getenv("APP_PATH")
	return runSQLFromFile(db, path+"/resources/database/up.sql")
}

func runSQLFromFile(db *gorm.DB, file string) error {
	schema, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	statements := strings.Split(string(schema), ";")

	for _, statement := range statements {
		command := strings.Join(strings.Fields(statement), " ")
		if command == "" {
			continue
		}

		tx := db.Exec(command)
		if tx.Error == nil {
			continue
		}

		// Error: 1065 = Query was empty
		if strings.Contains(tx.Error.Error(), "Error 1065") {
			continue
		}

		return err
	}

	return nil
}

func createAdminUser(userRepository user.Repository) {
	admin := models.User{
		Name:     "Administrator",
		Email:    "admin@admin.com",
		Login:    "admin",
		Password: "admin",
	}

	admin.HashPassword()
	_ = userRepository.Save(&admin)
}

func newPrivateEndpoint(routeGroup *gin.RouterGroup, httpMethod string, endpoint string, handler gin.HandlerFunc) {
	switch httpMethod {
	case http.MethodPost:
		routeGroup.POST(endpoint, middlewares.Authenticate(), handler)
	case http.MethodGet:
		routeGroup.GET(endpoint, middlewares.Authenticate(), handler)
	case http.MethodPut:
		routeGroup.PUT(endpoint, middlewares.Authenticate(), handler)
	case http.MethodDelete:
		routeGroup.DELETE(endpoint, middlewares.Authenticate(), handler)
	}
}

func newPublicEndpoint(routeGroup *gin.RouterGroup, httpMethod string, endpoint string, handler gin.HandlerFunc) {
	switch httpMethod {
	case http.MethodPost:
		routeGroup.POST(endpoint, middlewares.PublicAuthenticate(), handler)
	case http.MethodGet:
		routeGroup.GET(endpoint, middlewares.PublicAuthenticate(), handler)
	case http.MethodPut:
		routeGroup.PUT(endpoint, middlewares.PublicAuthenticate(), handler)
	case http.MethodDelete:
		routeGroup.DELETE(endpoint, middlewares.PublicAuthenticate(), handler)
	}
}

func getRunningPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
