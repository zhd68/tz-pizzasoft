package app

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/zhd68/tz-pizzasoft/internal/config"
	"github.com/zhd68/tz-pizzasoft/internal/storage"
	"github.com/zhd68/tz-pizzasoft/internal/storage/memory"
	"github.com/zhd68/tz-pizzasoft/internal/storage/postgresql"
	"github.com/zhd68/tz-pizzasoft/pkg/logging"
)

type App struct {
	config *config.Config
	logger *logging.Logger
	router *mux.Router
	server *http.Server
	store  storage.Storage
}

func NewApp(cfg *config.Config, logger *logging.Logger) (*App, error) {

	var db storage.Storage
	if cfg.StorageInMemory {
		db = memory.New()
	} else {
		sqlPath := fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.DB.Username,
			cfg.DB.Password,
			cfg.DB.Host,
			cfg.DB.Port,
			cfg.DB.Database,
		)
		store, err := postgresql.New(sqlPath)
		if err != nil {
			return nil, err
		}
		err = store.Init()
		if err != nil {
			return nil, err
		}
		db = store
	}
	logger.Infoln("database initialized and configured")

	a := &App{
		config: cfg,
		logger: logger,
		router: mux.NewRouter(),
		store:  db,
	}

	a.configureRouter()
	a.logger.Infoln("router initialized and configured")

	return a, nil
}

func (a *App) Run() {
	a.startHTTP()
}

func (a *App) startHTTP() {
	a.logger.Infoln("starting HTTP")

	var listener net.Listener

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", a.config.Listen.BindIP, a.config.Listen.Port))
	if err != nil {
		a.logger.Fatal(err)
	}
	a.logger.Infof("bind app to host: %s, port: %s", a.config.Listen.BindIP, a.config.Listen.Port)

	a.server = &http.Server{
		Handler:        a.router,
		MaxHeaderBytes: 1 << 20, // 1 MB
		WriteTimeout:   10 * time.Second,
		ReadTimeout:    10 * time.Second,
	}

	a.logger.Println("application successfully initialized and started")
	if err := a.server.Serve(listener); err != nil {
		a.logger.Fatal(err)
	}
}
