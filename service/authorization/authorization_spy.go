package authorization

import "github.com/ThuanyMendonca/project/model"

type AuthorizationServiceSpy struct {
	IAuthorizationService
	StatusCode        int
	AuthorizationResp *model.AuthorizeResponse
	AuthorizationErr  error
}

func (a *AuthorizationServiceSpy) Authorize() (int, *model.AuthorizeResponse, error) {
	return a.StatusCode, a.AuthorizationResp, a.AuthorizationErr
}
