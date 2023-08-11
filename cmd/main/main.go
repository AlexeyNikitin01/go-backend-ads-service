package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ads/internal/adapters/adrepo"
	"ads/internal/adapters/pgrepo"
	"ads/internal/adapters/userrepo"
	"ads/internal/app"
	grpcPort "ads/internal/ports/grpc"
	"ads/internal/ports/httpgin"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"github.com/joho/godotenv"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	var envs map[string]string
    envs, err := godotenv.Read(".env")
    if err != nil {
        log.Fatal("Error loading .env file")
    }
	db, err := pgrepo.NewPostgresDB(pgrepo.Config{
		Host:     envs["HOST"],
		Port:     envs["PORT"],
		Username: envs["USERNAME"],
		DBName:   envs["DBNAME"],
		SSLMode:  envs["SSLMODE"],
		Password: envs["POSTGRES_PASSWORD"],
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	a := app.NewApp(adrepo.New(), userrepo.New(), pgrepo.NewAuthPostgres(db))

	svr := httpgin.NewHTTPServer(":18080", a)

	httpServer := &http.Server{
		Addr:    ":18080",
		Handler: svr.Handler,
	}

	port := ":50054"

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(grpcPort.UnaryServerInterceptorPanicMethod),
		grpc.ChainUnaryInterceptor(grpcPort.UnaryServerInterceptorLogMethod),
	)

	svc := grpcPort.NewService(a)
	grpcPort.RegisterAdServiceServer(grpcServer, svc)

	eg, ctx := errgroup.WithContext(context.Background())

	sigQuit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			log.Printf("captured signal: %v\n", s)
			return fmt.Errorf("captured signal: %v", s)
		case <-ctx.Done():
			return nil
		}
	})

	// run grpc server
	eg.Go(func() error {
		log.Printf("starting grpc server, listening on %s\n", port)
		defer log.Printf("close grpc server listening on %s\n", port)

		errCh := make(chan error)

		defer func() {
			grpcServer.GracefulStop()
			_ = lis.Close()

			close(errCh)
		}()

		go func() {
			if err := grpcServer.Serve(lis); err != nil {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("grpc server can't listen and serve requests: %w", err)
		}
	})
	//Run rest
	eg.Go(func() error {
		log.Printf("starting http server, listening on %s\n", ":18080")
		defer log.Printf("close http server listening on %s\n", ":18080")

		errCh := make(chan error)

		defer func() {
			shCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if err := httpServer.Shutdown(shCtx); err != nil {
				log.Printf("can't close http server listening on %s: %s", ":18080", err.Error())
			}

			close(errCh)
		}()

		go func() {
			if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("http server can't listen and serve requests: %w", err)
		}
	})

	if err := eg.Wait(); err != nil {
		log.Printf("gracefully shutting down the servers: %s\n", err.Error())
	}

	log.Println("servers were successfully shutdown")
}
