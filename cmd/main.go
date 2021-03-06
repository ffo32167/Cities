package main

import (
	"context"
	"fmt"
	"github.com/ffo32167/cities/internal"
	"github.com/ffo32167/cities/internal/http"
	"github.com/ffo32167/cities/internal/remoteDB"
	"go.uber.org/zap"
	"os"
	"time"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("cant start logger: %s \n", err)
		os.Exit(1)
	}
	defer func() {
		err := logger.Sync()
		if err != nil {
			fmt.Printf("cant sync logger: %s \n", err)
		}
	}()
	cityCache := internal.New()
	cityCache.Cities[0] = "Moscow"

	rdb := remoteDB.New(1 * time.Second)
	ctx, _ := context.WithCancel(context.TODO()) //cancel нужно будет добавить только при обвязке завершения программы os.Signal и всё вот это вот
	reqChan := make(chan internal.CityReqMessage, 100)

	go remoteDB.RunRDB(ctx, rdb, reqChan)

	server := http.New(":8080")
	err = server.Run(logger, cityCache, reqChan)
	if err != nil {
		logger.Error("cant start api server:", zap.Error(err))
		os.Exit(1)
	}
}

// кэш использовать библиотечный
// сделать тип  type response struct{ City, error }
// кэш и бд хранят одинаковое количество данных(по 5)
//
