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

func RunRDB(ctx context.Context, rdb *RemoteDB, reqChan <-chan internal.CityReqMessage) {
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-reqChan:
			job.ResChan <- rdb.process(job.CityReq.ID)
		}
	}
}

func (rdb *RemoteDB) process(ID int) internal.CityResp {
	var city internal.CityResp
	if val, ok := rdb.data[ID]; ok {
		city = internal.CityResp{City: internal.City{ID: ID, Name: val}}
	}
	// обработать ошибку не найдено
	time.Sleep(rdb.Cooldown)
	return city
}
