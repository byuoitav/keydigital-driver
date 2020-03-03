package keydigital

import (
	"context"
	"net"
	"time"

	"github.com/byuoitav/connpool"
)

type KeyDigitalVideoSwitcher struct {
	Address string
	Pool    *connpool.Pool
}

const (
	carriageReturn = 0x0D
)

var (
	_defaultTTL   = 30 * time.Second
	_defaultDelay = 250 * time.Millisecond
)

func CreateVideoSwitcher(ctx context.Context, addr string) (*KeyDigitalVideoSwitcher, error) {
	p := &KeyDigitalVideoSwitcher{
		Address: addr,
		Pool: &connpool.Pool{
			TTL:   _defaultTTL,
			Delay: _defaultDelay,
			NewConnection: func(ctx context.Context) (net.Conn, error) {
				dial := net.Dialer{}
				return dial.DialContext(ctx, "tcp", addr+":3629")
			},
		},
	}

	return p, nil
}
