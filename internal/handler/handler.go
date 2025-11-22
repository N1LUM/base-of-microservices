package handler

import (
	"site-constructor/internal/auth"
	"site-constructor/internal/middleware"
	"site-constructor/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	services   *service.Service
	jwtManager *auth.JWTManager
}

func NewHandler(services *service.Service, jwtManager *auth.JWTManager) *Handler {
	logrus.Info("[Handler] Initialized Handler")
	return &Handler{
		services:   services,
		jwtManager: jwtManager,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	logrus.Info("[Handler.InitRoutes] Trying to initialize routes")

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	authorization := router.Group("/api/auth")
	{
		authorization.POST("/login", h.Login)     // Логин
		authorization.POST("/refresh", h.Refresh) // Рефреш токена
	}

	usersOpen := router.Group("/api/users")
	{
		usersOpen.POST("/", h.createUser) // Регистрация
	}

	api := router.Group("/api")
	api.Use(middleware.JWTAuthMiddleware(h.jwtManager))
	{
		users := api.Group("/users")
		{
			users.GET("/", h.getUsers)                               // Получить список пользователей
			users.GET("/:id", h.getUserByID)                         // Получить пользователя по ID
			users.PATCH("/:id", h.updateUser)                        // Обновить пользователя по ID
			users.DELETE("/:id", h.deleteUser)                       // Удалить пользователя по ID
			users.GET("/by-username/:username", h.getUserByUsername) // Получить пользователя по username
		}
	}

	logrus.Info("[Handler.InitRoutes] Successfully initialized routes")

	return router
}
