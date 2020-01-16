package bcs

import (
	"fmt"
	"net/http"
	"os"
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

func (ab *AuthBody) GetToken() {
	body := &AuthResponse{}
	req := RestRequest{
		Method: http.MethodPost,
		Path:   authResource,
		Data:   ab,
	}

	err := req.Send(body)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		os.Exit(1)
	}

	authToken = body.AuthenticationInfo.AuthToken
}
