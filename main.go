package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/ttblack/CredaAutherNode/config"
	"github.com/ttblack/CredaAutherNode/node"
	"github.com/ttblack/CredaAutherNode/signal"
)

func main() {

	var wg sync.WaitGroup
	// Hook interceptor for os signals.
	shutdownInterceptor, err := signal.Intercept()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	authernode, err := node.New(cfg)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	authernode.Start(&wg, &shutdownInterceptor)
	wg.Wait()
}
