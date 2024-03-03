package handler

import (
	"api/pkg/service"
	"github.com/gin-gonic/gin" // Импорт пакета Gin для создания маршрутов и обработки HTTP-запросов.
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// Handler представляет обработчик запросов
type Handler struct {
	services *service.Service
}

// NewHandler создает новый экземпляр обработчика с переданным сервисом
func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

// InitRoutes Метод инициализирует все маршруты для приложения и возвращает экземпляр маршрутизатора Gin
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	config := router.Group("/api/v1/good")
	{
		config.POST("/create/:projectId", h.createGood)
		config.PATCH("/update", h.updateGood)
		config.DELETE("/remove/:id", h.deleteGood)
		config.GET("/list", h.getGoods)
		config.PATCH("/reprioritize/", h.reprioritize)
	}

	return router
}
