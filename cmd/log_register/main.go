package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"time"

	"github.com/hsequeda/drone/storage"
)

func main() {
	intervalStr, ok := os.LookupEnv("LOG_REGISTER_INTERVAL")
	if !ok {
		log.Fatal("LOG_REGISTER_INTERVAL is empty")
		return
	}

	interval, err := strconv.ParseInt(intervalStr, 10, 64)
	if err != nil {
		log.Fatal("LOG_REGISTER_INTERVAL is empty")
		return
	}

	file, err := createFile()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	defer file.Close()

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.SetOutput(file)

	st := storage.NewInMemory()

	println("Registering Drones battery level")

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	execute(ctx, st)
	exit := make(chan struct{})
	go func() {
		defer func() { exit <- struct{}{} }()
		for {
			select {
			case <-time.Tick(time.Duration(interval) * time.Second):
				execute(ctx, st)
			case <-ctx.Done():
				return
			}
		}
	}()

	// exit when the gorutine is finished
	<-exit
	println("server exited properly")
}

func execute(ctx context.Context, st *storage.InMemory) error {
	drones, err := st.Drones(ctx)
	if err != nil {
		return errors.New("fetch error")
	}

	for _, d := range drones {
		log.Printf("serial(%s)-battery_level: %d%%", d.Serial, d.BatteryCapacity)
	}

	return nil
}

func createFile() (*os.File, error) {
	pwd, _ := os.Getwd()
	dir := filepath.Join(pwd, "/logs")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, errors.New("create /logs dir")
	}

	randomName := strconv.FormatInt(time.Now().UnixNano(), 10)
	filename := filepath.Join(dir, randomName)
	file, err := os.Create(filename)
	if err != nil {
		return nil, errors.New("create log file")
	}

	return file, nil
}
