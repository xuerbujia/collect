package tailfile

import (
	"collect/kafka"
	"context"
	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	"log"
	"time"
)

type Tails struct {
	EntryToTail map[CollectEntry]*tail.Tail
	TailToCtx   map[*tail.Tail]*CtxCancel
}
type CtxCancel struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}

func (t *Tails) Collect() {
	for collectEntry, tailTask := range t.EntryToTail {
		go Run(t.TailToCtx[tailTask].Ctx, tailTask, collectEntry.Topic)
	}
}
func Run(ctx context.Context, tailFile *tail.Tail, topic string) {
	for {
		select {
		case line, ok := <-tailFile.Lines:
			if !ok {
				log.Println("read line failed")
				time.Sleep(time.Second)
				continue
			}
			msg := sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(line.Text)}
			kafka.SendToChan(&msg)
		case <-ctx.Done():
			return
		}
	}
}
