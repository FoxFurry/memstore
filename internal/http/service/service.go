package service

import (
	"context"
	"github.com/FoxFurry/GoKeyValueStore/internal/cluster"
	"github.com/FoxFurry/GoKeyValueStore/internal/command"
	"github.com/FoxFurry/GoKeyValueStore/internal/http/httperr"
	"github.com/FoxFurry/GoKeyValueStore/internal/http/model"
	"net/http"
	"strings"
)

type Service interface {
	Execute([]model.Command) ([]string, error)
}

type service struct {
	data cluster.Cluster
}

func New(ctx context.Context) Service {
	newCluster := cluster.New()
	newCluster.Initialize(ctx)

	return &service{
		data: newCluster,
	}
}

func (s *service) Execute(modelTrns []model.Command) ([]string, error) {
	var trns []command.Command
	for _, modelTrn := range modelTrns {
		switch strings.ToUpper(modelTrn.CmdType) {
		case "GET":
			trns = append(trns, command.Get(modelTrn.Key))
		case "SET":
			trns = append(trns, command.Set(modelTrn.Key, modelTrn.Value))
		default:
			return nil, httperr.New("unknown command", http.StatusBadRequest)
		}
	}

	result, err := s.data.Execute(trns)

	if err != nil {
		return nil, httperr.WrapHttp(err, "could not execute command", http.StatusInternalServerError)
	}

	return result, nil
}
