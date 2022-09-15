package main

import (
	"collect/conf"
	"collect/etcd"
	"collect/kafka"
	"collect/tailfile"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func init() {
	err := conf.Init()

	if err != nil {
		panic(err)
	}
	err = kafka.Init(conf.GetConf().Kafka)

	if err != nil {
		panic(err)
	}
	err = etcd.Init(conf.GetConf().Etcd)

	if err != nil {
		panic(err)
	}
}

func main() {

	list, err := etcd.GetCollectPathList(conf.GetConf().Etcd.CollectKey)

	if err != nil {
		panic(err)
	}
	tailfile.InitTails(list)
	tailfile.GetTails().Collect()
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()
	select {}
}
