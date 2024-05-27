package api

import (
	"contactcenter-api/common/constants"
	authMdw "contactcenter-api/middleware/auth"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Engine *gin.Engine
}

func NewServer() *Server {
	engine := gin.New()
	authMdw.SetupGoGuardian()
	engine.Use(gin.Recovery())
	engine.Use(CORSMiddleware())
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": constants.SERVICE_NAME,
			"version": constants.VERSION,
			"time":    time.Now().Unix(),
		})
	})
	engine.Static("/v1/crm/contact/upload", "/root/go/src/public/upload")
	engine.Static("/v1/crm/images", "./public/logos/")
	server := &Server{Engine: engine}
	return server
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	}
}

func (server *Server) Start(port string) {
	v := make(chan struct{})
	go func() {
		if err := server.Engine.Run(":" + port); err != nil {
			log.WithError(err).Error("failed to start service")
			close(v)
		}
	}()
	log.Infof("service %v listening on port %v", constants.SERVICE_NAME, port)
	<-v
}
