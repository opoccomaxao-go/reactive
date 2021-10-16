package reactive

import (
	"log"
	"sync"

	"github.com/opoccomaxao-go/event"
)

type BoolMap map[string]bool

type Bool struct {
	boolMap BoolMap
	*event.Emitter
	mu      sync.RWMutex
	history BoolMap
	logger  *log.Logger
}

func NewBool() *Bool {
	return &Bool{
		boolMap: BoolMap{},
		history: BoolMap{},
		Emitter: event.NewEmitter(),
		logger:  devnullLogger,
	}
}

func (b *Bool) SetLogger(logger *log.Logger) {
	b.logger = logger
}

func (b *Bool) Set(name string, value bool) {
	b.mu.Lock()
	b.boolMap[name] = value
	b.commit()
	b.mu.Unlock()
}

func (b *Bool) Get(name string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.history[name]
}

func (b *Bool) batchSet(batch BoolMap) {
	for k, v := range batch {
		b.boolMap[k] = v
	}
}

func (b *Bool) BatchSet(batch BoolMap) {
	b.mu.Lock()
	b.batchSet(batch)
	b.commit()
	b.mu.Unlock()
}

func (b *Bool) Update(name string, fn func(bool) bool) {
	b.mu.Lock()
	b.boolMap[name] = fn(b.history[name])
	b.commit()
	b.mu.Unlock()
}

func (b *Bool) BatchUpdate(fn func(BoolMap) BoolMap) {
	b.mu.Lock()
	b.batchSet(fn(b.copy()))
	b.commit()
	b.mu.Unlock()
}

func (b *Bool) commit() {
	var toCommit []boolKeyValue
	for k, v := range b.boolMap {
		if hV, ok := b.history[k]; !ok || hV != v {
			toCommit = append(toCommit, boolKeyValue{
				Key:   k,
				Value: v,
			})
			b.history[k] = v
		}
	}
	if len(toCommit) > 0 {
		b.logger.Printf("Bool.Commit: %v\n", toCommit)
	}
	go b.fire(toCommit)
}

func (b *Bool) fire(update []boolKeyValue) {
	for _, pair := range update {
		b.Emit(pair.Key, pair.Value)
	}
}

func (b *Bool) SafeCopy() map[string]interface{} {
	b.mu.RLock()
	defer b.mu.RUnlock()
	res := make(map[string]interface{})
	for key, value := range b.history {
		res[key] = value
	}
	return res
}

func (b *Bool) copy() BoolMap {
	res := BoolMap{}
	for key, value := range b.history {
		res[key] = value
	}
	return res
}

func (b *Bool) SafeCopyBool() BoolMap {
	b.mu.RLock()
	res := b.copy()
	b.mu.RUnlock()
	return res
}
