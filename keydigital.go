package keydigital

import (
	"context"
	"net"
	"time"

	"github.com/byuoitav/connpool"
)

type VideoSwitcher struct {
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

func CreateVideoSwitcher(ctx context.Context, addr string, log connpool.Logger) (*VideoSwitcher, error) {
	p := &VideoSwitcher{
		Address: addr,
		Pool: &connpool.Pool{
			TTL:   _defaultTTL,
			Delay: _defaultDelay,
			NewConnection: func(ctx context.Context) (net.Conn, error) {
				dial := net.Dialer{}
				return dial.DialContext(ctx, "tcp", addr+":23")
			},
			Logger: log,
		},
	}

	return p, nil
}
