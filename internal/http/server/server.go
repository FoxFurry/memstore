package server

import (
	"context"
	"github.com/FoxFurry/GoKeyValueStore/internal/http/httperr"
	"github.com/FoxFurry/GoKeyValueStore/internal/http/model"
	"github.com/FoxFurry/GoKeyValueStore/internal/http/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Cluster interface {
	Start() error
}

type cluster struct {
	service service.Service
}

func New(ctx context.Context) Cluster {
	return &cluster{
		service: service.New(ctx),
	}
}

func (c *cluster) Start() error {
	server := gin.Default()

	server.POST("/execute", c.handleExecute)

	server.Run()

	return nil
}

func (c *cluster) handleExecute(ctx *gin.Context) {
	var request model.TransactionRequest

	if err := ctx.BindJSON(&request); err != nil {
		httperr.BadRequest(ctx, httperr.Wrap(err, "could not bind request"))
		return
	}

	execute, err := c.service.Execute(request.Commands)
	if err != nil {
		httperr.InternalError(ctx, httperr.Wrap(err, "internal error"))
		return
	}

	ctx.JSON(http.StatusOK, model.TransactionResponse{Results: execute})
}
