package cityHandler

import (
	"encoding/json"
	"fmt"
	"github.com/ffo32167/test2/internal"
	"go.uber.org/zap"
	"net/http"
)

type CityHandler struct {
	log         *zap.Logger
	citiesCache *internal.CitiesCache
	reqChan     chan<- internal.CityReq
	resChan     <-chan internal.City
}

func New(log *zap.Logger, citiesCache *internal.CitiesCache, reqChan chan<- internal.CityReq, resChan <-chan internal.City) CityHandler {
	return CityHandler{log: log, citiesCache: citiesCache, reqChan: reqChan, resChan: resChan}
}

type CityRequest struct {
	ID int `json:"ID"`
}

type CityResponse struct {
	ID   int    `json:"ID"`
	Name string `json:"Name"`
}

func (h CityHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var city CityRequest
	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&city)
	if err != nil {
		h.log.Error("CityHandler: cant decode user position:", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(res, err.Error())
		return
	}

	cityResp, err := h.citiesCache.Get(city.DecodeCityReq(), h.reqChan, h.resChan)

	if err != nil {
		h.log.Error("CityHandler: cant encode response:", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(res, err.Error())
		return
	}
	// Сделать переходник между внутренней city и структурой для ответа наружу
	err = json.NewEncoder(res).Encode(EncodeCiyResp(cityResp))
	if err != nil {
		h.log.Error("CityHandler: cant encode response:", zap.Error(err))
	}
}

func (c CityRequest) DecodeCityReq() internal.CityReq {
	return internal.CityReq{ID: c.ID}
}
func EncodeCiyResp(city internal.City) CityResponse {
	return CityResponse{ID: city.ID, Name: city.Name}
}
