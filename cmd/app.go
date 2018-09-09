package cmd

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/javking07/crawler/conf"

	"github.com/spf13/viper"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

// App ...
type App struct {
	AppServer *http.Server
	AppRouter *chi.Mux
	AppLogger *logrus.Logger
	AppConfig *conf.AppConfig
}

// RunApp runs
func (a App) RunApp(vp *viper.Viper) {
	logrus.Println("starting app")
	// extract config to initialize app loggers
	a.ExtractConfig(vp)
	logrus.Printf("using the following config %+v", *a.AppConfig)

	a.InitLogger()

	a.InitRoutes()

	// todo init database

	a.InitServer()

	a.AppServer.ListenAndServe()

}

// ExtractConfig ...

func (a *App) ExtractConfig(vp *viper.Viper) {
	logrus.Println("extracting config...")
	var config *conf.AppConfig
	err := vp.UnmarshalExact(&config)

	if err != nil {
		logrus.Panicf("error parsing config: %s", err.Error())

	}
	a.AppConfig = config
}

// InitLogger ...
func (a *App) InitLogger() {
	logrus.Println("initializing logger...")
	var level string
	logger := logrus.New()

	logger.Formatter = &logrus.JSONFormatter{}

	// Use logrus for standard log output
	// Note that `log` here references stdlib's log
	// Not logrus imported under the name `log`.
	log.SetOutput(logger.Writer())

	// extract logging level from config is exists
	level = a.AppConfig.Logging.Level

	if level, err := logrus.ParseLevel(level); err != nil {
		// if error parsin error log level, default to warn
		logrus.Warnf("error creating logger: %s", err.Error())
		logger.SetLevel(logrus.WarnLevel)
	} else {
		logrus.Infof("setting logger to level %s", level)
		logger.SetLevel(level)
	}

	a.AppLogger = logger

}

// Initroutes ...
func (a *App) InitRoutes() {
	a.AppLogger.Println("initializing routes...")
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
	a.AppRouter.HandleFunc("/health", Health)
}

// InitConfig ...
func (a *App) InitServer() {
	a.AppLogger.Println("initializing server...")
	// cfg := &tls.Config{
	// 	MinVersion:               tls.VersionTLS12,
	// 	CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
	// 	PreferServerCipherSuites: true,
	// 	CipherSuites: []uint16{
	// 		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	// 		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	// 		tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	// 		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	// 	},
	// }

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", a.AppConfig.Server.Port),
		Handler: a.AppRouter,
		// TLSConfig:    cfg,
		// TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	a.AppServer = srv

}
