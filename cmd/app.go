package cmd

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gocql/gocql"

	"github.com/javking07/food-crawler/model"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/coocood/freecache"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/javking07/food-crawler/conf"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	cassandraDb = "cassandra"
	postgresDb  = "postgres"
)

// App ...
type App struct {
	AppServer   *http.Server
	AppClient   *http.Client
	AppRouter   *chi.Mux
	AppDatabase *gocql.Session
	AppCache    *freecache.Cache
	AppLogger   zerolog.Logger
	AppConfig   *conf.AppConfig
}

// RunApp runs
func (a App) RunApp() {
	a.BootstrapApp()

	// only run app is db connection present
	if a.AppDatabase != nil {
		defer a.AppDatabase.Close() // for graceful db shutdown
		a.AppServer.ListenAndServe()
	}
}

// BootstrapApp prepares app with config to run
func (a *App) BootstrapApp() {
	log.Info().Msgf("bootstrapping app with the following config:/n %+V", a.AppConfig)

	a.InitLogger()
	a.InitRoutes()
	a.InitCache()
	if err := a.InitDatabase(); err != nil {
		log.Fatal().Msgf("error bootstrapping database: %s", err.Error())
	}
	a.InitClient()
	a.InitServer()
}

// ExtractConfig ...

func (a *App) ExtractConfig(vp *viper.Viper) {
	log.Print("extracting config...")
	var config *conf.AppConfig
	err := vp.UnmarshalExact(&config)
	if err != nil {
		logrus.Panicf("error parsing config: %s", err.Error())
	}
	a.AppConfig = config
}

// InitLogger ...
func (a *App) InitLogger() {
	var level string
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	log.Output(logger)

	// extract logging level from config is exists
	level = a.AppConfig.Logging.Level

	if level, err := zerolog.ParseLevel(level); err != nil {
		// if error parsing error log level, default to warn
		log.Warn().Msgf("error creating logger: %s", err.Error())
		logger.Level(zerolog.WarnLevel)
	} else {
		logger.Level(level)
	}

	a.AppLogger = logger
	a.AppLogger.Info().Msgf("initialized logger to level `%s`", level)
}

// InitServer bootstraps app server
func (a *App) InitRoutes() {
	a.AppRouter = chi.NewRouter()
	// A good base middleware stack
	a.AppRouter.Use(middleware.RequestID)
	a.AppRouter.Use(middleware.RealIP)
	a.AppRouter.Use(middleware.Logger)
	a.AppRouter.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	a.AppRouter.Use(middleware.Timeout(60 * time.Second))

	// add actual api routes
	a.AppRouter.HandleFunc("/food-crawler/v1/health", a.Health)
	a.AppRouter.Handle("/food-food-crawler/v1/metrics", promhttp.Handler())
}

// InitServer bootstraps app server
func (a *App) InitCache() {
	cacheSize := a.AppConfig.Cache.Size
	log.Info().Msgf("initializing cache with size of `%d` bytes", cacheSize)
	cache := freecache.NewCache(cacheSize)
	a.AppCache = cache
}

// InitS bootstraps app server
func (a *App) InitDatabase() error {
	db, err := model.BootstrapCassandra(a.AppConfig.Database)
	if err != nil {
		return err
	}
	a.AppDatabase = db
	return nil
}

// InitServer bootstraps app server
func (a *App) InitClient() {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
		},
		Timeout: 0,
	}
	a.AppClient = client
}

// InitServer bootstraps app server
func (a *App) InitServer() {

	var cfg *tls.Config
	if !a.AppConfig.Server.TLS {
		cfg = &tls.Config{}
	} else {
		cert, err := tls.LoadX509KeyPair(
			a.AppConfig.Server.Cert,
			a.AppConfig.Server.Key)

		if err != nil {
			log.Fatal().Msgf("Unable to load cert/key: %s", err)
		}

		cfg = &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			InsecureSkipVerify: false,
			Certificates:       []tls.Certificate{cert},
		}
		cfg.BuildNameToCertificate()
	}

	addr := fmt.Sprintf(":%s", a.AppConfig.Server.Port)
	a.AppServer = &http.Server{
		Addr:      addr,
		Handler:   a.AppRouter,
		TLSConfig: cfg,
		// TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	a.AppLogger.Info().Msgf("initialized server on port %+v", a.AppServer.Addr)

}
