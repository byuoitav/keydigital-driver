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
	regGetInput   = regexp.MustCompile("Video Output : Input = ([0-9]{2}),")
)

// GetInputByOutput .
func (vs *KeyDigitalVideoSwitcher) GetInputByOutput(ctx context.Context, output string) (string, error) {
	var input string

	err := vs.Pool.Do(ctx, func(conn connpool.Conn) error {

		cmd := []byte("STA\r\n")
		n, err := conn.Write(cmd)
		switch {
		case err != nil:
			return fmt.Errorf("failed to write to connection: %w", err)
		case n != len(cmd):
			return fmt.Errorf("failed to write to connection: wrote %v/%v bytes", n, len(cmd))
		}

		var match [][]string
		for len(match) == 0 {
			c, err := conn.ReadUntil(carriageReturn, 3*time.Second)
			if err != nil {
				return fmt.Errorf("failed to read from connection: %w", err)
			}

			match = regGetInput.FindAllStringSubmatch(string(c), -1)
		}

		input = match[0][1]
		input = input[1:]
		return nil
	})

	if err != nil {
		return "", err
	}

	return input, nil

}

// SetInputByOutput .
func (vs *KeyDigitalVideoSwitcher) SetInputByOutput(ctx context.Context, output, input string) error {
	return vs.Pool.Do(ctx, func(conn connpool.Conn) error {
		cmd := []byte(fmt.Sprintf("SPO0%sSI0%s\r\n", output, input))
		n, err := conn.Write(cmd)
		switch {
		case err != nil:
			return fmt.Errorf("failed to write to connection: %w", err)
		case n != len(cmd):
			return fmt.Errorf("failed to write to connection: wrote %v/%v bytes", n, len(cmd))
		}

		buf, err := conn.ReadUntil(carriageReturn, 3*time.Second)
		if err != nil {
			return fmt.Errorf("failed to read from connection: %w", err)
		}

		if strings.Contains(string(buf), "FAILED") {
			return ErrOutOfRange
		}

		return nil
	})
}
