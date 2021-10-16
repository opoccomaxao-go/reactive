package reactive

import (
	"log"
	"sync"

	"github.com/opoccomaxao-go/event"
)

type Float64Map map[string]float64

type Float64 struct {
	float64Map Float64Map
	*event.Emitter
	mu      sync.RWMutex
	history Float64Map
	logger  *log.Logger
}

func NewFloat64() *Float64 {
	return &Float64{
		float64Map: Float64Map{},
		history:    Float64Map{},
		Emitter:    event.NewEmitter(),
		logger:     devnullLogger,
	}
}

func (f *Float64) SetLogger(logger *log.Logger) {
	f.logger = logger
}

func (f *Float64) Set(name string, value float64) {
	f.mu.Lock()
	f.float64Map[name] = value
	f.commit()
	f.mu.Unlock()
}

func (f *Float64) Get(name string) (res float64) {
	f.mu.RLock()
	res = f.history[name]
	f.mu.RUnlock()
	return
}

func (f *Float64) GetOk(name string) (res float64, ok bool) {
	f.mu.RLock()
	res, ok = f.history[name]
	f.mu.RUnlock()
	return
}

func (f *Float64) batchSet(batch Float64Map) {
	for k, v := range batch {
		f.float64Map[k] = v
	}
}

func (f *Float64) BatchSet(batch Float64Map) {
	f.mu.Lock()
	f.batchSet(batch)
	f.commit()
	f.mu.Unlock()
}

func (f *Float64) Update(name string, fn func(float64) float64) {
	f.mu.Lock()
	f.float64Map[name] = fn(f.history[name])
	f.commit()
	f.mu.Unlock()
}

func (f *Float64) BatchUpdate(fn func(Float64Map) Float64Map) {
	f.mu.Lock()
	f.batchSet(fn(f.copy()))
	f.commit()
	f.mu.Unlock()
}

func (f *Float64) commit() {
	var toCommit []floatKeyValue
	for k, v := range f.float64Map {
		if hV, ok := f.history[k]; !ok || hV != v {
			toCommit = append(toCommit, floatKeyValue{
				Key:   k,
				Value: v,
			})
			f.history[k] = v
		}
	}
	if len(toCommit) > 0 {
		f.logger.Printf("Float.Commit: %v\n", toCommit)
	}
	go f.fire(toCommit)
}

func (f *Float64) fire(update []floatKeyValue) {
	for _, pair := range update {
		f.Emit(pair.Key, pair.Value)
	}
}

func (f *Float64) SafeCopy() map[string]interface{} {
	f.mu.RLock()
	res := make(map[string]interface{})
	for key, value := range f.history {
		res[key] = value
	}
	f.mu.RUnlock()
	return res
}

func (f *Float64) copy() Float64Map {
	res := Float64Map{}
	for key, value := range f.history {
		res[key] = value
	}
	return res
}

func (f *Float64) SafeCopyFloat() Float64Map {
	f.mu.RLock()
	res := f.copy()
	f.mu.RUnlock()
	return res
}
