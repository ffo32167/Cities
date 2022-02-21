package remoteDB

import (
	"context"
	"github.com/ffo32167/cities/internal"
	"time"
)

type RemoteDB struct {
	data     map[int]string
	Cooldown time.Duration
}

func New(cooldown time.Duration) *RemoteDB {
	d := make(map[int]string)
	d[0] = "Moscow"
	d[1] = "SPB"
	d[2] = "Tula"
	d[3] = "Tver"
	d[4] = "Kaluga"
	return &RemoteDB{data: d, Cooldown: cooldown}
}

func RunRDB(ctx context.Context, rdb *RemoteDB, reqChan <-chan internal.CityReq, resChan chan<- internal.City) {
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-reqChan:
			resChan <- rdb.process(job)
		}
	}
}

func (rdb *RemoteDB) process(cityReq internal.CityReq) internal.City {
	var city internal.City
	if val, ok := rdb.data[cityReq.ID]; ok {
		city = internal.City{ID: cityReq.ID, Name: val}
	}
	time.Sleep(rdb.Cooldown)
	return city
}
