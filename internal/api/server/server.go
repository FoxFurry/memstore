package server

import (
	"context"
	"github.com/FoxFurry/memstore/internal/api/httperr"
	"github.com/FoxFurry/memstore/internal/api/model"
	"github.com/FoxFurry/memstore/internal/api/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type KeyValueServer interface {
	Start() error
}

type keyValueServer struct {
	service service.Service
}

func New(ctx context.Context) KeyValueServer {
	return &keyValueServer{
		service: service.New(ctx),
	}
}

func (s *keyValueServer) Start() error {
	gin.DisableConsoleColor()

	server := gin.New()
	server.Use(gin.Logger())

	v1 := server.Group("/v1")
	{
		v1.POST("/execute", s.handleExecute)
	}

	server.Run()

	return nil
}

func (s *keyValueServer) handleExecute(c *gin.Context) {
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
