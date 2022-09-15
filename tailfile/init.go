package tailfile

import (
	"context"
	"github.com/hpcloud/tail"
)

var tails *Tails

func GetTails() *Tails {
	return tails
}

type CollectEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

func InitTails(list []*CollectEntry) {
	tails = &Tails{
		EntryToTail: make(map[CollectEntry]*tail.Tail),
		TailToCtx:   make(map[*tail.Tail]*CtxCancel),
	}
	for _, collectEntry := range list {
		task, err := NewTailTask(collectEntry.Path)
		if err != nil {
			continue
		}
		tails.EntryToTail[*collectEntry] = task
		ctx, cancel := context.WithCancel(context.Background())
		tails.TailToCtx[task] = &CtxCancel{ctx, cancel}
	}
	return
}
func NewTailTask(path string) (tailFile *tail.Tail, err error) {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	tailFile, err = tail.TailFile(path, config)
	if err != nil {
		return
	}
	return
}
