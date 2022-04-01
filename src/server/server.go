package server

import (
	"github.com/gin-gonic/gin"
	"github.com/zytekaron/zog-server/src/config"
	"github.com/zytekaron/zog-server/src/database"
	"github.com/zytekaron/zog-server/src/server/logs"
	"github.com/zytekaron/zog-server/src/types"
	"net/http"
)

type Server struct {
	Config *config.Config
	Addr   string
	Logs   database.Controller[*types.Log]
	Tokens database.Controller[*types.Token]
}

func New(cfg *config.Config, lc database.Controller[*types.Log], tc database.Controller[*types.Token]) *Server {
	return &Server{
		Addr:   cfg.Server.Addr,
		Logs:   lc,
		Tokens: tc,
	}
}

func (s *Server) Start() error {
	r := gin.Default()

	logs.New(s.Config, s.Logs, s.Tokens).Register(r)

	return http.ListenAndServe(s.Addr, r)
}
