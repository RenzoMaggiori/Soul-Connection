package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"soul-connection.com/api/src/lib"
)

type AuthModel struct {
	Auth interface {
		GetApiKey() string
	}
}

type ApiKeyAuth struct {
	ApiKey string
}

func (model *AuthModel) Login(res http.ResponseWriter, req *http.Request) {
	var c struct {
		Email    string
		Password string
	}
	err := json.NewDecoder(req.Body).Decode(&c)
	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
	jwt, err := lib.Auth(lib.LoginCredentials{
		XGroupAuthentication: model.Auth.GetApiKey(),
		AuthEmail:            c.Email,
		AuthPassword:         c.Password,
	})
	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
	res.Header().Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
}

func (a ApiKeyAuth) GetApiKey() string {
	return a.ApiKey
}
