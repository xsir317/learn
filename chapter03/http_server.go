package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kratos/kratos/pkg/sync/errgroup"
)

func main() {
	srv := &http.Server{Addr: ":8080"}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello world\n")
	})

	g := errgroup.WithCancel(context.Background())
	g.Go(func(ctx context.Context) error {
		go func() {
			<-ctx.Done()
			fmt.Println("Ctx Done in server")
			srv.Close()
		}()
		return srv.ListenAndServe()
	})

	g.Go(func(ctx context.Context) error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Ctx Done in signal select")
				return nil
			case <-c:
				fmt.Println("signal exit")
				return errors.New("got signal")
			default:
				time.Sleep(time.Second)
			}
		}
		return nil
	})

	g.Wait()
}
