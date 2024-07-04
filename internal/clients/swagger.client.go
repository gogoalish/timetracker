package clients

import (
	"context"
	"net/http"
	"net/url"

	"github.com/gogoalish/timetracker/config"
	"github.com/gogoalish/timetracker/internal/clients/swagger"
	"github.com/gogoalish/timetracker/internal/service"
)

type APIService struct {
	client *swagger.APIClient
}

func NewAPIService(cfg *config.Config) (*APIService, error) {
	apiURL, err := url.Parse(cfg.APIURL)
	if err != nil {
		return nil, err
	}
	swagCfg := swagger.NewConfiguration()
	swagCfg.Host = apiURL.Host
	swagCfg.Scheme = apiURL.Scheme

	client := swagger.NewAPIClient(swagCfg)

	return &APIService{client}, nil
}

func (s *APIService) InfoGet(ctx context.Context, passportSerie int32, passportNumber int32) (*service.Person, error) {
	person, resp, err := s.client.DefaultApi.InfoGet(ctx, passportSerie, passportNumber)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return &service.Person{
			Surname:    person.Surname,
			Name:       person.Name,
			Patronymic: person.Patronymic,
			Address:    person.Address,
		}, nil
	case http.StatusBadRequest:
		return nil, service.ErrBadRequest
	default:
		return nil, service.ErrApiInternal
	}

}
