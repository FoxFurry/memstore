package service

import (
	"context"
	"errors"
	"github.com/FoxFurry/GoKeyValueStore/internal/cluster"
	"github.com/FoxFurry/GoKeyValueStore/internal/command"
	"github.com/FoxFurry/GoKeyValueStore/internal/http/model"
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
			return nil, errors.New("bad command")
		}
	}

	result, err := s.data.Execute(trns)

	if err != nil {
		return nil, err
	}

	return result, nil
}
