package schemas

import (
	"encoding/json"
	"errors"
	"net"
	"strings"
)

type IPAddress net.IP

type LoginEvent struct {
	Username  string    `json:"username" valid:"type(string),required"`
	TimeStamp int64     `json:"unix_timestamp" valid:"type(int64),required"`
	EventUuid string    `json:"event_uuid" valid:"uuid,required"`
	IPAddress IPAddress `json:"ip_address"`
}

func (i IPAddress) IP() net.IP {
	return net.IP(i)
}

func (i IPAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.IP().String())
}

func (i *IPAddress) UnmarshalJSON(b []byte) error {
	rawIP := string(b)
	ip := strings.Trim(rawIP, `"`)

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return errors.New("not a valid textual representation of an IP address")
	}

	*i = IPAddress(parsedIP)
	return nil
}
