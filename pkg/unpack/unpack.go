package unpack

import (
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type Unpacker interface {
	UnpackAll(dir string) error
}

// for tests only
var _ Unpacker = &noopUnpacker{}

var _ Unpacker = &CompositeUnpacker{}
var _ Unpacker = &unzip{}

type noopUnpacker struct {
}

func NoopUnpacker() Unpacker {
	return &noopUnpacker{}
}

func (*noopUnpacker) UnpackAll(dir string) error {
	return nil
}

type CompositeUnpacker struct {
	Unpackers []Unpacker
	logger    zap.Logger
}

func NewUnpacker(logger *zap.SugaredLogger) Unpacker {
	return &CompositeUnpacker{
		Unpackers: []Unpacker{
			&unzip{logger: logger},
			&untar{logger: logger},
		},
	}
}

func (c CompositeUnpacker) UnpackAll(dir string) error {
	var err error
	for _, u := range c.Unpackers {
		err = multierr.Append(err, u.UnpackAll(dir))
	}

	return err
}
