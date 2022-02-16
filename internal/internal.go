package internal

import "sync"

type City struct {
	ID   int
	Name string
}
type CitiesCache struct {
	Cities map[int]string
	Mx     sync.Mutex
}

type CityReq struct {
	ID int
}

func New() *CitiesCache {
	return &CitiesCache{Cities: make(map[int]string)}
}
func (citiesCache *CitiesCache) Get(CityReq CityReq, reqChan chan<- CityReq, resChan <-chan City) (City, error) {
	citiesCache.Mx.Lock()
	name, ok := citiesCache.Cities[CityReq.ID]
	citiesCache.Mx.Unlock()
	if !ok {
		reqChan <- CityReq
		city := <-resChan // изменить, смотреть на 27 минуте
		name = city.Name
	}
	return City{ID: CityReq.ID, Name: name}, nil
}
