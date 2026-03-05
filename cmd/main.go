package main

import (
	"log"
	"net/http"

	"connectrpc.com/connect"
	"github.com/shizumico/arcane/pkg/logger"
	"github.com/shizumico/arcane/pkg/sqlite"
	"github.com/thebeyond-net/node-agent/config"
	amneziawgAdapter "github.com/thebeyond-net/node-agent/internal/adapters/amneziawg"
	"github.com/thebeyond-net/node-agent/internal/adapters/grpc/interceptors"
	peerHandlers "github.com/thebeyond-net/node-agent/internal/adapters/grpc/peers"
	"github.com/thebeyond-net/node-agent/internal/adapters/repositories/sqlite/ipallocation"
	"github.com/thebeyond-net/node-agent/internal/core/application/peers"
	"github.com/thebeyond-net/node-agent/pkg/amneziawg/v1/amneziawgv1connect"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	appLogger, err := logger.Init(cfg.LogLevel)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v\n", err)
	}
	defer appLogger.Sync()

	appLogger.Info("Starting Node Agent server")

	db, err := sqlite.New(cfg.DatabasePath, cfg.MigrationsPath)
	if err != nil {
		appLogger.Fatal("Failed to connect to sqlite", zap.Error(err))
	}

	ipallocationRepo, err := ipallocation.NewRepository(db)
	if err != nil {
		appLogger.Fatal("Failed to create ip allocation repository", zap.Error(err))
	}

	amneziawgAdapter, err := amneziawgAdapter.New(
		cfg.NodeIP,
		cfg.NetworkCIDR,
		cfg.ClientDNS,
	)
	if err != nil {
		appLogger.Fatal("Failed to create amneziawg adapter", zap.Error(err))
	}

	peersInteractor, err := peers.NewInteractor(
		amneziawgAdapter,
		ipallocationRepo,
		"10.10.0.0/20",
	)
	if err != nil {
		appLogger.Fatal("Failed to create peers interactor", zap.Error(err))
	}

	amneziawgHandlers := peerHandlers.NewAmneziaWGServiceServer(peersInteractor)
	authInterceptor := connect.WithInterceptors(interceptors.NewAuthInterceptor(cfg.AuthSecret))

	mux := http.NewServeMux()
	path, handler := amneziawgv1connect.NewAmneziaWGServiceHandler(
		amneziawgHandlers,
		authInterceptor,
	)
	mux.Handle(path, handler)

	appLogger.Info("Listening on :" + cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, h2c.NewHandler(mux, &http2.Server{})))
}
