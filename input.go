package keydigital

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/connpool"
)

// GetInputByOutput .
func (vs *KeyDigitalVideoSwitcher) GetInputByOutput(ctx context.Context, output string) (string, error) {
	var input string

	work := func(conn connpool.Conn) error {
		conn.Write([]byte(fmt.Sprintf("STA\r\n")))

		//capture response
		var buf bytes.Buffer
		io.Copy(&buf, conn)

		//regex black magic
		reg, err := regexp.Compile("Video Output : Input = ([0-9]{2}),")
		if err != nil {
			log.L.Errorf("Failed to read from %s : %s", vs.Address, err.Error())
			return fmt.Errorf("failed to create regex: %s", err)
		}
		ReturnInput := reg.FindAllStringSubmatch(fmt.Sprintf("%s", buf.Bytes()), -1)

		input := ReturnInput[0][1]
		input = input[1:]
		return nil
	}

	err := vs.Pool.Do(ctx, work)
	if err != nil {
		return "", err
	}

	return input, nil

}

// SetInputByOutput .
func (vs *KeyDigitalVideoSwitcher) SetInputByOutput(ctx context.Context, output, input string) error {

	//execute telnet command to switch input
	work := func(conn connpool.Conn) error {
		command := fmt.Sprintf("SPO0" + output + "SI0" + input)
		log.L.Infof("%s", command)
		conn.Write([]byte("\r\n"))
		conn.Write([]byte(command + "\r\n"))
		b, err := readUntil(CARRIAGE_RETURN, conn, 10)
		if err != nil {
			return fmt.Errorf("failed to read from connection: %s", err)
		}

		if strings.Contains(string(b), "FAILED") {
			return fmt.Errorf("input or output is out of range. Input recieved: %s  Output Recieved: %s", input, output)
		}
		log.L.Infof("response: %s", b)
		// response := strings.Split(fmt.Sprintf("%s", b), "")
		// test := strings.Split(fmt.Sprintf("%s", response), "x")
		// log.L.Infof("test: %s", test)
		return nil
	}

	err := vs.Pool.Do(ctx, work)
	if err != nil {
		return err
	}

	return nil
}
