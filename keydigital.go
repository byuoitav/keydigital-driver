package keydigital

import (
	"context"
	"net"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/connpool"
	"github.com/fatih/color"
)

type KeyDigitalVideoSwitcher struct {
	Address string
	Pool    *connpool.Pool
}

var (
	_defaultTTL   = 45 * time.Second
	_defaultDelay = 100 * time.Millisecond
)

type options struct {
	ttl   time.Duration
	delay time.Duration
}

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

func WithTTL(t time.Duration) Option {
	return optionFunc(func(o *options) {
		o.ttl = t
	})
}

func WithDelay(t time.Duration) Option {
	return optionFunc(func(o *options) {
		o.delay = t
	})
}

func NewVideoSwitcher(addr string, opts ...Option) *KeyDigitalVideoSwitcher {
	options := options{
		ttl:   _defaultTTL,
		delay: _defaultDelay,
	}

	for _, o := range opts {
		o.apply(&options)
	}

	p := &KeyDigitalVideoSwitcher{
		Address: addr,
		Pool: &connpool.Pool{
			TTL:   options.ttl,
			Delay: options.delay,
		},
	}

	p.Pool.NewConnection = func(ctx context.Context) (net.Conn, error) {
		// addr, err := net.ResolveTCPAddr("tcp", addr+":23")
		// if err != nil {
		// 	return nil, err
		// }

		// conn, err := net.DialTCP("tcp", nil, addr)
		// if err != nil {
		// 	return nil, err
		// }

		dial := net.Dialer{}
		conn, err := dial.DialContext(ctx, "tcp", p.Address+":3629")
		if err != nil {
			return nil, err
		}

		pconn := connpool.Wrap(conn)

		//This was used in the older code, not sure if we still need it so I am keeping it and defaulting it to true
		readWelcome := true
		if readWelcome {
			color.Set(color.FgMagenta)
			log.L.Infof("Reading welcome message")
			color.Unset()
			_, err := readUntil(CARRIAGE_RETURN, pconn, 3)
			if err != nil {
				return conn, err
			}
		}
		return conn, err
	}

	return p
}
