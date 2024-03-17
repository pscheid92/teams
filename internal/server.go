package internal

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Server struct {
	*gin.Engine
	repository DataRepository
	jwt        *JwtHelper
}

func NewServer(jwt *JwtHelper, repository DataRepository) *Server {
	return &Server{
		Engine:     gin.New(),
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

func (s *Server) buildHealthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	}
}

func (s *Server) buildLoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := LoginRequest{}
		if err := c.BindJSON(&request); err != nil {
			err = fmt.Errorf("login: %w", err)
			c.Error(err)
			c.Status(http.StatusBadRequest)
			return
		}

		key, found := s.repository.GetUserPublicKey(request.Username)
		if !found {
			c.Status(http.StatusNotFound)
			return
		}

		isValid, err := VerifyChallenge(request.Challenge, request.Username, request.Timestamp, key)
		if err != nil || !isValid {
			c.Status(http.StatusUnauthorized)
			return
		}

		accessToken, err := s.jwt.Create(request.Username, time.Now())
		if err != nil {
			err = fmt.Errorf("login: %w", err)
			c.Error(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		response := LoginResponse{AccessToken: accessToken}
		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) buildVerifyHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.Query("access_token")
		if accessToken == "" {
			c.Status(http.StatusBadRequest)
			return
		}

		claims, err := s.jwt.Validate(accessToken)
		if err != nil {
			err = fmt.Errorf("verify: %w", err)
			c.Error(err)
			c.Status(http.StatusBadRequest)
			return
		}

		username := claims.Subject
		found := s.repository.UserExists(username)
		if !found {
			c.Status(http.StatusNotFound)
			return
		}

		response := VerifyResponse{Username: username}
		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) buildTeamHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		teamID := c.Param("id")

		members, found := s.repository.GetTeamMembers(teamID)
		if !found {
			c.Status(http.StatusNotFound)
			return
		}

		response := TeamResponse{TeamID: teamID, Members: members}
		c.JSON(http.StatusOK, response)
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
