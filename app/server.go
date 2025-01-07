package app

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/cors"
)

type Server struct {
	uc *UseCase
	hz *server.Hertz
}

func NewServer(uc *UseCase) *Server {

	hz := server.Default(server.WithExitWaitTime(time.Second))

	hz.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	s := &Server{
		uc: uc,
		hz: hz,
	}

	hz.GET("/presignedPostPolicy", s.presignedPostPolicy)
	hz.GET("/images", s.listImages)

	return s
}

func (s *Server) Spin() {

	s.hz.Spin()
}

func (s *Server) presignedPostPolicy(ctx context.Context, c *app.RequestContext) {
	key := c.Query("key")
	contentType := c.Query("contentType")
	url, fields, err := s.uc.PresignedPostPolicy(ctx, key, contentType)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{"message": err.Error()})
		return
	}
	c.JSON(consts.StatusOK, utils.H{"url": url, "fields": fields})
}

func (s *Server) listImages(ctx context.Context, c *app.RequestContext) {
	arr := s.uc.ListImages(ctx)

	c.JSON(consts.StatusOK, arr)
}
