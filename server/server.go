package server

import (
	"gin-server/middleware"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/kjk/dailyrotate"
)

type Router struct {
	*gin.Engine
}

const (
	ENV_PRODUCTION  = "pro"
	ENV_DEVELOPMENT = "dev"
)

func New(env string) *Router {
	router := &Router{gin.New()}
	router.ForwardedByClientIP = true
	router.Use(gin.Recovery())
	if env == ENV_PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = ioutil.Discard
	} else {
		router.Use(gin.Logger())
	}

	return router
}

func (r *Router) LoadStatic(urlPrefix, folder string) {
	r.Use(static.Serve(urlPrefix, static.LocalFile(folder, false)))
}

func (r *Router) LoadTemplate(pattern string) {
	r.LoadHTMLGlob(pattern)
}

func (r *Router) CORS(config middleware.CORSConfig) {
	r.Use(middleware.CORSMiddleware(config))
}

func (r *Router) SecureHeader() {
	r.Use(middleware.SecureHeader())
}

func (r *Router) AllowHosts(domains ...string) {
	r.Use(middleware.CheckHost(domains...))
}

func (r *Router) FuckBot() {
	r.Use(middleware.BlockUserAgentMalicious())
}

func (r *Router) AccessLogDaily(dir string) error {
	if len(dir) > 0 {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	pathFormat := filepath.Join(dir, "access_02_01_2006.log")
	w, err := dailyrotate.NewFile(pathFormat, nil)
	if err != nil {
		return err
	}
	r.Use(middleware.Logger(w))

	return nil
}
