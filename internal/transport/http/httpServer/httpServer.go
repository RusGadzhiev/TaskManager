package httpServer

import (
	"context"
	"net/http"
	"time"

	"github.com/RusGadzhiev/TaskManager/internal/config"
	"github.com/RusGadzhiev/TaskManager/internal/transport/http/httpHandler"
	"go.uber.org/zap"
)

type HttpServer struct {
	server *http.Server
}

func NewHttpServer(ctx context.Context, h *httpHandler.HttpHandler, cfg *config.HTTPServer) *HttpServer {
	return &HttpServer{
		server: &http.Server{
			Addr:         ":" + cfg.Port,
			ReadTimeout:  cfg.Timeout,
			WriteTimeout: cfg.Timeout,
			IdleTimeout:  cfg.IdleTimeout,
			Handler:      h.Router(),
		}}
}

func (srv *HttpServer) Run(ctx context.Context, logger *zap.SugaredLogger) error {
	go func() {
		if err := srv.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("HTTP server error: %v", err)
		}
	}()

	logger.Info("Start listen at " + srv.server.Addr)
	<-ctx.Done()

	logger.Info(ctx, "shutting down server gracefully")

	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := srv.server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	return nil
}
