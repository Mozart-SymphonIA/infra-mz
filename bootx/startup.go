package bootx

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"
)

type Runner func(ctx context.Context) error

func Start(name string, run Runner) {
	if run == nil {
		panic(fmt.Sprintf("bootx: run function is required for %s", name))
	}
	
	LoadEnv()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	defer func() {
		if r := recover(); r != nil {
			log.Printf("[%s] PANIC: %v\n%s", name, r, debug.Stack())
			os.Exit(1)
		}
	}()

	if err := run(ctx); err != nil {
		log.Printf("[%s] fatal: %v", name, err)
		os.Exit(1)
	}
}

func Service[C any, S any](
	name string,
	loadConfig func() (C, error),
	buildServer func(C) (S, error),
	run func(context.Context, C, S) error,
) {
	Start(name, func(ctx context.Context) error {
		for {
			if ctx.Err() != nil {
				return nil
			}

			cfg, err := loadConfig()
			if err != nil {
				log.Printf("[%s] config error: %v (retrying in 3s...)", name, err)
				retryWait(ctx)
				continue
			}

			srv, err := buildServer(cfg)
			if err != nil {
				log.Printf("[%s] build error: %v (retrying in 3s...)", name, err)
				retryWait(ctx)
				continue
			}

			err = run(ctx, cfg, srv)

			if closer, ok := any(srv).(interface{ Close() error }); ok {
				_ = closer.Close()
			}

			if ctx.Err() != nil {
				return nil
			}

			retryWait(ctx)
		}
	})
}

func retryWait(ctx context.Context) {
	select {
	case <-time.After(3 * time.Second):
	case <-ctx.Done():
	}
}
