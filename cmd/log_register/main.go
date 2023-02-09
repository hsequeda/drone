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

	"github.com/hsequeda/drone"
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

	st := drone.NewStorage()

	println("Registering Drones battery level")

	execute(st)
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	exit := make(chan struct{})
	go func() {
		defer func() { exit <- struct{}{} }()
		for {
			select {
			case <-time.Tick(time.Duration(interval) * time.Second):
				execute(st)
			case <-ctx.Done():
				return
			}
		}
	}()

	// exit when the gorutine is finished
	<-exit
	println("server exited properly")
}

func execute(st *drone.Storage) {
	for _, d := range st.Drones() {
		log.Printf("serial(%s)-battery_level: %d%%", d.Serial, d.BatteryCapacity)
	}
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
