package internal

import (
	"fmt"
	"github.com/dghubble/sling"
	"net/http"
)

type Client struct {
	*sling.Sling
}

func NewClient(baseURL string) *Client {
	client := sling.New().Base(baseURL)
	return &Client{client}
}

func (c *Client) Login(request LoginRequest) (LoginResponse, error) {
	result := LoginResponse{}
	response, err := c.Post("login").BodyJSON(request).ReceiveSuccess(&result)
	if err != nil {
		err = fmt.Errorf("login: %w", err)
		return result, err
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("login: unsuccessful status code %d", response.StatusCode)
		return result, err
	}

	return result, nil
}

type verifyParams struct {
	AccessToken string `url:"access_token"`
}

func (c *Client) Verify(token string) (VerifyResponse, error) {
	result := VerifyResponse{}
	query := verifyParams{token}

	response, err := c.Get("verify").QueryStruct(query).ReceiveSuccess(&result)
	if err != nil {
		err = fmt.Errorf("verify: %w", err)
		return result, err
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("verify: unsuccessful status code %d", response.StatusCode)
		return result, err
	}

	return result, err
}

func (c *Client) Team(team string) (TeamResponse, error) {
	result := TeamResponse{}
	response, err := c.Get("teams/" + team).ReceiveSuccess(&result)
	if err != nil {
		err = fmt.Errorf("teams: %w", err)
		return result, err
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("teams: unsuccessful status code %d", response.StatusCode)
		return result, err
	}

	return result, err
}
