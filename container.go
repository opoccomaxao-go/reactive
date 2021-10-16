package reactive

import "log"

type Container struct {
	Float *Float64 `json:"float"`
	Bool  *Bool    `json:"bool"`

	logger *log.Logger
}

func NewContainer() *Container {
	return &Container{
		Float:  NewFloat64(),
		Bool:   NewBool(),
		logger: devnullLogger,
	}
}

func (c *Container) SetLogger(logger *log.Logger) {
	if logger != nil {
		c.logger = logger
		c.Float.SetLogger(logger)
		c.Bool.SetLogger(logger)
	}
}

func (c *Container) SafeCopy() map[string]interface{} {
	return map[string]interface{}{
		"b": c.Bool.SafeCopy(),
		"f": c.Float.SafeCopy(),
	}
}
