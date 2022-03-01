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

type CityReqMessage struct {
	CityReq CityReq
	ResChan chan CityResp
}

type CityResp struct {
	City City
	err  error
}

func New() *CitiesCache {
	return &CitiesCache{Cities: make(map[int]string)}
}

func (citiesCache *CitiesCache) Get(CityReq CityReq, reqChan chan<- CityReqMessage) (c City, err error) {
	citiesCache.Mx.Lock()
	name, ok := citiesCache.Cities[CityReq.ID]
	citiesCache.Mx.Unlock()
	if ok {
		c.ID = CityReq.ID
		c.Name = name
	} else {
		resChan := make(chan CityResp, 1)
		reqChan <- CityReqMessage{CityReq: CityReq, ResChan: resChan}
		resp := <-resChan
		c = resp.City
		err = resp.err
	}
	return c, err
}
