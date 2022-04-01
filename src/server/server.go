package server

import (
	"github.com/gin-gonic/gin"
	"github.com/zytekaron/zog-server/src/config"
	"github.com/zytekaron/zog-server/src/database"
	"github.com/zytekaron/zog-server/src/server/logs"
	"net/http"
)

type Server struct {
	Config *config.Config
	Addr   string
	Logs   database.LogController
	Tokens database.TokenController
}

func New(cfg *config.Config, lc database.LogController, tc database.TokenController) *Server {
	return &Server{
		Addr:   cfg.Server.Addr,
		Logs:   lc,
		Tokens: tc,
	}
}

func (s *Server) Start() error {
	r := gin.Default()

	logs.New(s.Config, s.Tokens, s.Logs).Register(r)

	return http.ListenAndServe(s.Addr, r)
}
