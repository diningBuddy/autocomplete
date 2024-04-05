package app

import (
	"context"
	"crypto/tls"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"net/http"
	"sync"
	"time"

	"github.com/skku/autocomplete/pkg/config"

	log "github.com/sirupsen/logrus"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/skku/autocomplete/app/handler"
	"github.com/skku/autocomplete/app/middleware"
	"github.com/skku/autocomplete/app/model"
)

type App struct {
	Router *mux.Router
	Redis  *model.AutocompleteRedis

	server *http.Server

	m       *sync.RWMutex
	Version *model.Version
}

func (a *App) Initialize(config *config.Properties) {
	autocompleteRedis := redis.NewClient(newRedisOption(config.SearchAutocompleteRedis.Addr, config.SearchAutocompleteRedis.Password))

	a.m = &sync.RWMutex{}

	a.Redis = &model.AutocompleteRedis{
		Search: autocompleteRedis,
	}
	a.setVersion()
	a.Router = mux.NewRouter()
	a.setRouters()
	a.setMiddlewares(middleware.WithRecover)

	a.setMiddlewares(middleware.WithLogging)
	log.SetLevel(log.WarnLevel)

	log.SetFormatter(
		&log.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		},
	)
	log.SetLevel(log.DebugLevel)
	a.setMiddlewares(middleware.WithRecover)

	wrappedHandler := httptrace.WrapHandler(a.Router, "autocomplete", "http.request")
	a.server = &http.Server{Addr: config.Addr, Handler: wrappedHandler}
	go a.runVersionChecker()
}

func (a *App) setRouters() {
	a.Get("/healthcheck", a.handleRequest(handler.HealthCheck))

	a.Get("/restaurant", a.handleRequest(handler.RestaurantAutocomplete))
}

func (a *App) setMiddlewares(f func(http.Handler) http.Handler) {
	a.Router.Use(f)
}

func (a *App) setVersion() {
	var cv map[string]string

	cv = map[string]string{}
	for _, v := range handler.RestaurantKeyVersions {
		rv, err := a.Redis.Search.Get("restaurant:" + v + ":version").Result()

		if err != nil && v == handler.RestaurantVersion {
			log.Errorf("failed to get default restaurant version with err : %s", err.Error())
		}
		cv[v] = rv
		log.Debugf("restaurant version : %s", rv)
	}

	a.Version = &model.Version{
		Restaurant: cv,
	}
}

func (a *App) runVersionChecker() {
	for range time.Tick(time.Minute * 1) {
		for _, v := range handler.RestaurantKeyVersions {
			com, err := a.Redis.Search.Get("restaurant:" + v + ":version").Result()
			if err != nil {
				log.Error(err.Error())
				continue
			}
			if com != a.Version.Restaurant[v] {
				log.Infof("setting restaurant version %s : %s", v, com)
				a.m.Lock()
				a.Version.Restaurant[v] = com
				a.m.Unlock()
			}
		}
	}
}

// Get wraps the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

func (a *App) Run() {

	log.Infof("starting server at %s", a.server.Addr)
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error in ListenAndServe: %s", err)
	}
}

type RequestHandlerFunction func(rd *model.AutocompleteRedis, v *model.Version, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.m.RLock()
		v := a.Version
		a.m.RUnlock()
		handler(a.Redis, v, w, r)
	}
}

func (a *App) GracefulShutdown(ctx context.Context) error {
	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func newRedisOption(addr, password string) *redis.Options {
	log.Infof("set up redis client %s : %s", addr, password)
	if len(password) > 0 {
		return &redis.Options{
			Addr:     addr,
			Password: password,
			TLSConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
			DialTimeout: 500 * time.Millisecond,
			ReadTimeout: 500 * time.Millisecond,
		}
	} else {
		return &redis.Options{
			Addr:        addr,
			DialTimeout: 500 * time.Millisecond,
			ReadTimeout: 500 * time.Millisecond,
		}
	}

}
