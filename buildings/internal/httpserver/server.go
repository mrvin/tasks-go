package httpserver

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	// import docs generated with swag init.
	_ "github.com/mrvin/tasks-go/buildings/docs"
	"github.com/mrvin/tasks-go/buildings/internal/httpserver/handlers"
	"github.com/mrvin/tasks-go/buildings/internal/storage"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type Conf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Server struct {
	*gin.Engine
}

// @title           Buildings API
// @version         1.0
// @description     This is a sample server buildings server.
// @host      		localhost:8081
func New(st storage.Storage) *Server {
	router := gin.Default()

	router.GET("/health", handlers.Health)

	router.POST("/buildings", handlers.NewCreateBuilding(st))
	router.GET("/buildings", handlers.NewListBuildings(st))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &Server{router}
}

func (s *Server) Start(conf *Conf) {
	addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	log.Print("Start HTTP server: http://" + addr)
	if err := s.Run(addr); err != nil {
		log.Fatalf("Error while starting HTTP server: %v", err)
	}
}
