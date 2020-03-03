package keydigital

import (
	"crypto/tls"
	"errors"
	"fmt"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/connpool"

	telnet "github.com/reiver/go-telnet"
)

const (
	CARRIAGE_RETURN           = 0x0D
	LINE_FEED                 = 0x0A
	SPACE                     = 0x20
	DELAY_BETWEEN_CONNECTIONS = time.Second * 10
)

var tlsConfig *tls.Config
var caller telnet.Caller

func init() {
	tlsConfig = &tls.Config{}
}

func readUntil(delimeter byte, conn connpool.Conn, timeoutInSeconds int) ([]byte, error) {
	conn.SetReadDeadline(time.Now().Add(time.Duration(int64(timeoutInSeconds)) * time.Second))

	buffer := make([]byte, 128)
	message := []byte{}

	for !charInBuffer(delimeter, buffer) {
		_, err := conn.Read(buffer)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error reading response: %s", err.Error()))
			log.L.Infof("%s", err.Error())
			return message, err
		}

		message = append(message, buffer...)
	}

	return removeNil(message), nil
}

func removeNil(b []byte) (ret []byte) {
	for _, c := range b {
		switch c {
		case '\x00':
			break
		default:
			ret = append(ret, c)
		}
	}
	return ret
}

func charInBuffer(toCheck byte, buffer []byte) bool {
	for _, b := range buffer {
		if toCheck == b {
			return true
		}
	}

	return false
}
