package http

import (
	"github.com/ffo32167/cities/internal"
	"github.com/ffo32167/cities/internal/http/cityHandler"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

type ApiServer struct {
	port string
}

func New(port string) ApiServer {
	return ApiServer{port: port}
}

func (as ApiServer) Run(logger *zap.Logger, cache *internal.CitiesCache, reqChan chan<- internal.CityReqMessage) error {
	ch := cityHandler.New(logger, cache, reqChan)
	router := mux.NewRouter()
	router.Handle("/", ch).Methods("GET")
	err := http.ListenAndServe(as.port, router)
	if err != nil {
		return err
	}
	return nil
}
