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
	root_ctx, root_cancel := context.WithCancel(context.Background())

	srv := &http.Server{Addr: ":8080"}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello world\n")
		url := r.URL.Path[1:]
		//允许在访问/shutdown 的时候触发关闭。
		if url == "shutdown" {
			root_cancel()
		}
	})

	g := errgroup.WithCancel(root_ctx)
	g.Go(func(ctx context.Context) error {
		go func() {
			<-ctx.Done()
			fmt.Println("Ctx Done , Server shutdown")
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
			case sig := <-c:
				fmt.Println("signal exit")
				return errors.New("Signal " + sig.String())
			default:
				time.Sleep(time.Second)
			}
		}
	})

	exit_reason := g.Wait()
	fmt.Printf("main exit, reason: %v", exit_reason)
}
