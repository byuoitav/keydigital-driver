package keydigital

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/byuoitav/connpool"
)

var (
	ErrOutOfRange = errors.New("input or output is out of range")
	regGetInput   = regexp.MustCompile("Video Output  *: Input = ([0-9]{2}),")
)

// GetInputByOutput .
func (vs *VideoSwitcher) GetInputByOutput(ctx context.Context, output string) (string, error) {
	var input string

	if vs.Pool.Logger != nil {
		vs.Pool.Logger.Infof("getting the current input")
	}

	err := vs.Pool.Do(ctx, func(conn connpool.Conn) error {

		vs.Pool.Logger.Infof("writing to the connection")

		cmd := []byte("STA\r\n")
		n, err := conn.Write(cmd)
		switch {
		case err != nil:
			vs.Pool.Logger.Warnf("failed to write to connection")
			return fmt.Errorf("failed to write to connection: %w", err)
		case n != len(cmd):
			return fmt.Errorf("failed to write to connection: wrote %v/%v bytes", n, len(cmd))
		}

		vs.Pool.Logger.Infof("reading from the connection")
		var match [][]string
		for len(match) == 0 {
			c, err := conn.ReadUntil(carriageReturn, 3*time.Second)
			if err != nil {
				vs.Pool.Logger.Warnf("failed to read from connection")
				return fmt.Errorf("failed to read from connection: %w", err)
			}

			match = regGetInput.FindAllStringSubmatch(string(c), -1)
		}

		input = match[0][1]
		input = strings.TrimPrefix(input, "0")
		return nil
	})

	if err != nil {
		return "", err
	}

	if vs.Pool.Logger != nil {
		vs.Pool.Logger.Infof(fmt.Sprintf("returning input - current input: %s", input))
	}

	return input, nil

}

// SetInputByOutput .
func (vs *VideoSwitcher) SetInputByOutput(ctx context.Context, output, input string) error {
	return vs.Pool.Do(ctx, func(conn connpool.Conn) error {

		if vs.Pool.Logger != nil {
			vs.Pool.Logger.Infof(fmt.Sprintf("writing command to change input - change to input: %s", input))
		}

		cmd := []byte(fmt.Sprintf("SPO0%sSI0%s\r\n", output, input))
		n, err := conn.Write(cmd)
		switch {
		case err != nil:
			return fmt.Errorf("failed to write to connection: %w", err)
		case n != len(cmd):
			return fmt.Errorf("failed to write to connection: wrote %v/%v bytes", n, len(cmd))
		}

		vs.Pool.Logger.Infof("reading from connection to see if there was an error")

		buf, err := conn.ReadUntil(carriageReturn, 3*time.Second)
		if err != nil {
			return fmt.Errorf("failed to read from connection: %w", err)
		}

		if strings.Contains(string(buf), "FAILED") {
			return ErrOutOfRange
		}

		if vs.Pool.Logger != nil {
			vs.Pool.Logger.Infof(fmt.Sprintf("successfully changed the input - current input: %s", input))
		}

		return nil
	})
}
