package web

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"

	"pipego/internal/efs"
	"pipego/internal/manager"

	"github.com/gin-gonic/gin"
)

type Server struct {
	manager  *manager.Manager
	username string
	password string
	port     int
}

func New(mgr *manager.Manager, username, password string, port int) *Server {
	return &Server{
		manager:  mgr,
		username: username,
		password: password,
		port:     port,
	}
}

func (s *Server) Start() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	auth := r.Group("/", gin.BasicAuth(gin.Accounts{
		s.username: s.password,
	}))

	auth.GET("/", func(c *gin.Context) {
		data, _ := efs.FS.ReadFile("index.html")
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	auth.GET("/api/routes", s.readRoutes)
	auth.POST("/api/routes", s.writeRoutes)

	addr := fmt.Sprintf("0.0.0.0:%d", s.port)
	log.Printf("[INFO] web management starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Printf("[ERROR] web server failed: %v", err)
	}
}

func (s *Server) readRoutes(c *gin.Context) {
	data, err := os.ReadFile("pipego.routes")
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusOK, gin.H{"content": ""})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "read error: " + err.Error()})
		return
	}
	encoded := base64.StdEncoding.EncodeToString(data)
	c.JSON(http.StatusOK, gin.H{"content": encoded})
}

func (s *Server) writeRoutes(c *gin.Context) {
	var req struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request: " + err.Error()})
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "base64 decode error: " + err.Error()})
		return
	}

	if err := os.WriteFile("pipego.routes", decoded, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "write error: " + err.Error()})
		return
	}

	log.Println("[INFO] routes file updated via web, reloading listeners...")

	if err := s.manager.LoadAndStart("pipego.routes"); err != nil {
		log.Printf("[ERROR] reload routes failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "file saved but reload failed: " + err.Error(),
		})
		return
	}

	log.Println("[INFO] routes reloaded successfully")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
