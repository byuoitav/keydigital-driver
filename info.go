package keydigital

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"regexp"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/common/structs"
	"github.com/byuoitav/connpool"
)

//GetHardwareInfo .
func (vs *KeyDigitalVideoSwitcher) GetHardwareInfo(ctx context.Context) (structs.HardwareInfo, error) {
	var verdata string
	var macaddr string
	var ipaddr string
	var resp structs.HardwareInfo

	work := func(conn connpool.Conn) error {
		conn.Write([]byte(fmt.Sprintf("STA\r\n")))

		//capture response
		var buf bytes.Buffer
		io.Copy(&buf, conn)

		//regex black magic
		reg, err := regexp.Compile("Host IP Address = ([0-9]{3}.[0-9]{3}.[0-9]{3}.[0-9]{3})")
		if err != nil {
			log.L.Errorf("Failed to read from %s : %s", vs.Address, err.Error())
			return fmt.Errorf("failed to create regex: %s", err)
		}
		ReturnInput := reg.FindAllStringSubmatch(fmt.Sprintf("%s", buf.Bytes()), -1)
		ipaddr := ReturnInput[0][1]
		log.L.Infof("IP Addr: %s", ipaddr)

		//Version
		reg, err = regexp.Compile("Version : ([0-9]+.[0-9]+)")
		if err != nil {
			log.L.Errorf("Failed to read from %s : %s", vs.Address, err.Error())
			return nerr.Translate(err).Add("failed to create regex")
		}
		ReturnInput = reg.FindAllStringSubmatch(fmt.Sprintf("%s", buf.Bytes()), -1)
		verdata := ReturnInput[0][1]
		log.L.Infof("Version: %s", verdata)

		//MacAddress
		reg, err = regexp.Compile("MAC Address = ([A-Z,0-9]{2}:[A-Z,0-9]{2}:[A-Z,0-9]{2}:[A-Z,0-9]{2}:[A-Z,0-9]{2}:[A-Z,0-9]{2})")
		if err != nil {
			log.L.Errorf("Failed to read from %s : %s", vs.Address, err.Error())
			return nerr.Translate(err).Add("failed to create regex")
		}
		ReturnInput = reg.FindAllStringSubmatch(fmt.Sprintf("%s", buf.Bytes()), -1)
		macaddr := ReturnInput[0][1]
		log.L.Infof("MAC: %s", macaddr)

		return nil
	}

	err := vs.Pool.Do(ctx, work)
	if err != nil {
		return resp, err
	}

	resp.FirmwareVersion = verdata
	resp.NetworkInfo.MACAddress = macaddr
	resp.NetworkInfo.IPAddress = ipaddr
	return resp, nil
}

//GetInfo .
func (vs *KeyDigitalVideoSwitcher) GetInfo(ctx context.Context) (interface{}, error) {
	var info interface{}
	return info, fmt.Errorf("not currently implemented")
}
