/*
Package server
Copyright © 2022 Arthur Isac isacartur@gmail.com

Describes http layer which wraps service layer. Provides a server interface with single Start method
*/
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

// MemStore wraps memstore service into http server and provides single method: Start
// TODO: Revisit what is a server and what it should do
type MemStore interface {
	Start(string) error
}

type memstore struct {
	service service.Service
}

// New creates a MemStore with initialized service
func New(ctx context.Context) MemStore {
	return &memstore{
		service: service.New(ctx),
	}
}

// Start creates a simple gin server and starts it. This method is blocking
func (s *memstore) Start(port string) error {
	gin.DisableConsoleColor()
	//gin.SetMode(gin.ReleaseMode)	// TODO: Add environment variable for develop/live

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
