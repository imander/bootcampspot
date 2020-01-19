package bcs

import (
	"net/http"
)

type AuthBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Success            bool   `json:"success"`
	ErrorCode          string `json:"errorCode"`
	ResetToken         string `json:"resetToken"`
	AuthenticationInfo struct {
		UserID     int    `json:"userId"`
		FirstLogin bool   `json:"firstLogin"`
		Active     bool   `json:"active"`
		AuthToken  string `json:"authToken"`
	} `json:"authenticationInfo"`
}

func (ab *AuthBody) GetToken() error {
	body := &AuthResponse{}
	req := RestRequest{
		Method: http.MethodPost,
		Path:   authResource,
		Data:   ab,
	}

	err := req.Send(body)
	if err != nil {
		return err
	}

	authToken = body.AuthenticationInfo.AuthToken
	return nil
}
