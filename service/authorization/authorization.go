package authorization

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ThuanyMendonca/project/model"
)

type IAuthorizationService interface {
	Authorize() (int, *model.AuthorizeResponse, error)
}

type AuthorizationService struct {
	client http.Client
}

func NewAuthorizationService(client http.Client) IAuthorizationService {
	return &AuthorizationService{client}
}

func (a *AuthorizationService) Authorize() (int, *model.AuthorizeResponse, error) {
	url := "https://run.mocky.io/v3/d02168c6-d88d-4ff2-aac6-9e9eb3425e31"

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("ocorreu um erro durante a criação da requisição, detalhes: %s", err)
	}

	req.Close = true

	req.Header.Add("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("ocorreu um erro estava fazendo a requisição, detalhes: %s", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("ocorreu um erro ao ler a resposta da requisição, detalhes: %s", err)
	}

	authorization := &model.AuthorizeResponse{}
	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, nil, err
	}

	if err := json.Unmarshal(body, &authorization); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, authorization, nil
}
