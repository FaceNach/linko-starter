package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"boot.dev/linko/internal/store"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	httpPort := flag.Int("port", 8899, "port to listen on")
	dataDir := flag.String("data", "./data", "directory to store data")
	flag.Parse()

	status := run(ctx, cancel, *httpPort, *dataDir)
	cancel()
	os.Exit(status)
}

func run(ctx context.Context, cancel context.CancelFunc, httpPort int, dataDir string) int {

	logger, cerrar, err := initializeLogger(os.Getenv("LINKO_LOG_FILE"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		return 1
	}

	if cerrar != nil {
		defer cerrar()
	}

	st, err := store.New(dataDir, logger)
	if err != nil {
		logger.Printf("failed to create store: %v\n", err)
		return 1
	}
	s := newServer(*st, httpPort, cancel, logger)
	var serverErr error
	go func() {
		serverErr = s.start()
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Println("Linko is shutting down")
	if err := s.shutdown(shutdownCtx); err != nil {
		logger.Printf("failed to shutdown server: %v\n", err)
		return 1
	}
	if serverErr != nil {
		logger.Printf("Server error: %v\n", serverErr)
		return 1
	}
	return 0
}

type closeFunc func() error

func initializeLogger(logFile string) (*log.Logger, closeFunc, error) {

	if len(logFile) == 0 {
		logger := log.New(os.Stderr, "", 0)
		return logger, nil, nil
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		return nil, nil, err
	}

	bufferedFile := bufio.NewWriterSize(file, 8192)

	var cerrar closeFunc = func() error {
		errFlush := bufferedFile.Flush()
		errClose := file.Close()

		return errors.Join(errFlush, errClose)
	}

	multiWriter := io.MultiWriter(os.Stderr, bufferedFile)

	logger := log.New(multiWriter, "", 0)

	return logger, cerrar, nil
}
