package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/acmk189/todo_app_REST_API/config"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatalf("failed to run: %+v", err)
		os.Exit(1)
	}
}

func run(cxt context.Context) error {
	ctx, stop := signal.NotifyContext(cxt, os.Interrupt, syscall.SIGTERM)
	defer stop()
	cfg, err := config.New()
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port: %d: %v", cfg.Port, err)
	}
	url := fmt.Sprintf("http://localhost:%d", l.Addr().String())
	log.Printf("start server listening at: %s", url)

	s := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(5 * time.Second)
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}

	eg, ctx := errgroup.WithContext(cxt)
	// 別ゴルーチンでHTTPサーバを起動
	eg.Go(func() error {
		// http.ErrServerClosed は http.Server.Shutdown() の正常終了なのでエラーではない
		if err := s.Serve(l); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	// チャネルからの終了通知を受け取る
	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}
	return eg.Wait()
}
