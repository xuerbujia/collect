package etcd

import (
	"collect/conf"
	"collect/tailfile"
	"context"
	"encoding/json"
	etcd "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

var client = new(etcd.Client)

func Init(config *conf.EtcdConfig) (err error) {
	client, err = etcd.New(etcd.Config{
		DialTimeout: time.Second * 5,
		Endpoints:   config.Address,
	})
	if err != nil {
		return err
	}
	go watch(config)
	return nil
}

func GetCollectPathList(key string) (entryList []*tailfile.CollectEntry, err error) {
	gr, err := client.Get(context.Background(), key)
	if err != nil {
		log.Println("get from etcd error", err)
		return
	}
	if len(gr.Kvs) == 0 {
		log.Println("get len:0 from etcd")
		return
	}
	err = json.Unmarshal(gr.Kvs[0].Value, &entryList)
	if err != nil {
		log.Println("entryList unmarshal error", err)
		return
	}
	return
}

func watch(config *conf.EtcdConfig) {
	watchChan := client.Watch(context.Background(), config.CollectKey)
	for {
		select {
		case <-watchChan:
			list, err := GetCollectPathList(config.CollectKey)
			if err != nil {
				continue
			}
			tails := tailfile.GetTails()
			var hash = map[tailfile.CollectEntry]bool{}
			for _, collectEntry := range list {
				hash[*collectEntry] = true
				if _, ok := tails.EntryToTail[*collectEntry]; !ok {
					newTailTask, err := tailfile.NewTailTask(collectEntry.Path)
					if err != nil {
						log.Println("a new tailTask was created failed", err)
						continue
					}
					ctx, cancel := context.WithCancel(context.Background())
					tails.TailToCtx[newTailTask] = &tailfile.CtxCancel{Ctx: ctx, Cancel: cancel}
					tails.EntryToTail[*collectEntry] = newTailTask
					go tailfile.Run(ctx, newTailTask, collectEntry.Topic)
				}
			}

			for collectEntry, tail := range tails.EntryToTail {
				if _, ok := hash[collectEntry]; !ok {
					ctxCancel := tails.TailToCtx[tail]
					ctxCancel.Cancel()
					_ = tail.Stop()
					delete(tails.TailToCtx, tail)
					delete(tails.EntryToTail, collectEntry)
				}
			}

		}
	}
}
