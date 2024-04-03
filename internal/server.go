package internal

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type Server struct {
	*echo.Echo
	repository DataRepository
	jwt        *JwtHelper
}

func NewServer(jwt *JwtHelper, repository DataRepository) *Server {
	return &Server{
		Echo:       echo.New(),
		repository: repository,
		jwt:        jwt,
	}
}

func (s *Server) InitRoutes() {
	s.GET("health", s.buildHealthHandler())
	s.POST("login", s.buildLoginHandler())
	s.GET("verify", s.buildVerifyHandler())
	s.GET("teams/:id", s.buildTeamHandler())
}

func (s *Server) buildHealthHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	}
}

func (s *Server) buildLoginHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		request := LoginRequest{}

		if err := c.Bind(&request); err != nil {
			err = fmt.Errorf("login: %w", err)
			c.Error(err)
			return c.NoContent(http.StatusBadRequest)
		}

		key, found := s.repository.GetUserPublicKey(request.Username)
		if !found {
			return c.NoContent(http.StatusNotFound)
		}

		isValid, err := VerifyChallenge(request.Challenge, request.Username, request.Timestamp, key)
		if err != nil || !isValid {
			return c.NoContent(http.StatusUnauthorized)
		}

		accessToken, err := s.jwt.Create(request.Username, time.Now())
		if err != nil {
			err = fmt.Errorf("login: %w", err)
			c.Error(err)
			return c.NoContent(http.StatusInternalServerError)
		}

		response := LoginResponse{AccessToken: accessToken}
		return c.JSON(http.StatusOK, response)
	}
}

func (s *Server) buildVerifyHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		accessToken := c.QueryParam("access_token")
		if accessToken == "" {
			return c.NoContent(http.StatusBadRequest)
		}

		claims, err := s.jwt.Validate(accessToken)
		if err != nil {
			err = fmt.Errorf("verify: %w", err)
			c.Error(err)
			return c.NoContent(http.StatusBadRequest)
		}

		username := claims.Subject
		found := s.repository.UserExists(username)
		if !found {
			return c.NoContent(http.StatusNotFound)
		}

		response := VerifyResponse{Username: username}
		return c.JSON(http.StatusOK, response)
	}
}

func (s *Server) buildTeamHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		teamID := c.Param("id")

		members, found := s.repository.GetTeamMembers(teamID)
		if !found {
			return c.NoContent(http.StatusNotFound)
		}

		response := TeamResponse{TeamID: teamID, Members: members}
		return c.JSON(http.StatusOK, response)
	}
}

type LoginRequest struct {
	Username  string    `json:"username"`
	Timestamp time.Time `json:"timestamp"`
	Challenge string    `json:"challenge"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type TeamResponse struct {
	TeamID  string   `json:"team_id"`
	Members []string `json:"member"`
}

type VerifyResponse struct {
	Username string `json:"username"`
}
