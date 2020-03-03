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

func CreateVideoSwitcher(ctx context.Context, addr string) (*KeyDigitalVideoSwitcher, error) {
	var err error

	p := &KeyDigitalVideoSwitcher{
		Address: addr,
		Pool: &connpool.Pool{
			TTL:   _defaultTTL,
			Delay: _defaultDelay,
		},
	}

	p.Pool.NewConnection = func(ctx context.Context) (net.Conn, error) {

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

	if err != nil {
		return p, err
	}

	return p, nil
}
