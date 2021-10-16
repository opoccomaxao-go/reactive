package reactive

import (
	"io"
	"log"
	"strconv"
)

var devnullLogger = log.New(io.Discard, log.Prefix(), log.Flags())

type floatKeyValue struct {
	Key   string
	Value float64
}

func (f floatKeyValue) String() string {
	return f.Key + " = " + strconv.FormatFloat(f.Value, 'g', -1, 64)
}

type boolKeyValue struct {
	Key   string
	Value bool
}

func (b boolKeyValue) String() string {
	return b.Key + " = " + strconv.FormatBool(b.Value)
}
