package kafka

import (
	"collect/conf"
	"github.com/Shopify/sarama"
	"log"
)

var (
	client  sarama.SyncProducer
	msgChan chan *sarama.ProducerMessage
)

func GetClient() sarama.SyncProducer {
	return client
}
func Init(kafkaConf *conf.KafkaConfig) (err error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	client, err = sarama.NewSyncProducer(kafkaConf.Address, config)
	if err != nil {
		return err
	}
	msgChan = make(chan *sarama.ProducerMessage, kafkaConf.ChanSize)
	go SendToKafka()
	return nil
}
func SendToChan(msg *sarama.ProducerMessage) {
	msgChan <- msg
}
func SendToKafka() {
	for msg := range msgChan {
		pid, offset, err := client.SendMessage(msg)
		if err != nil {
			log.Println("send to kafka error", err)
		}
		log.Printf("send to kafka successful pid:%v,offser:%v", pid, offset)
	}
}
