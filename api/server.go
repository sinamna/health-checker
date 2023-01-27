package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Router *gin.Engine
}

func CreateHttpServer() *Server {
	return &Server{
		Router: gin.Default(),
	}
}

func (s *Server) StartServer(port string) error {
	err := s.Router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		return err
	}
	return nil
}
