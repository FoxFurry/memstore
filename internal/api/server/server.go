package server

import (
	"context"
	"fmt"
	"github.com/FoxFurry/memstore/internal/api/httperr"
	"github.com/FoxFurry/memstore/internal/api/model"
	"github.com/FoxFurry/memstore/internal/api/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MemStore interface {
	Start(string) error
}

type memstore struct {
	service service.Service
}

func New(ctx context.Context) MemStore {
	return &memstore{
		service: service.New(ctx),
	}
}

func (s *memstore) Start(port string) error {
	gin.DisableConsoleColor()
	//gin.SetMode(gin.ReleaseMode)

	server := gin.New()
	server.Use(gin.Logger())

	v1 := server.Group("/v1")
	{
		v1.POST("/execute", s.handleExecute)
	}

	host := fmt.Sprintf(":%s", port)

	return server.Run(host)
}

func (s *memstore) handleExecute(c *gin.Context) {
	var request = new(model.TransactionRequest)

	if err := c.BindJSON(&request); err != nil {
		httperr.Handle(c, httperr.WrapHttp(err, "could not bind request", http.StatusBadRequest))
		return
	}

	if request.Commands == nil {
		httperr.Handle(c, httperr.New("empty request", http.StatusBadRequest))
		return
	}

	execute, err := s.service.Execute(request.Commands)
	if err != nil {
		httperr.Handle(c, err)
		return
	}

	c.JSON(http.StatusOK, model.TransactionResponse{Results: execute})
}
