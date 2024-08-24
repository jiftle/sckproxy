package proto

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
)

const (
	SOCKS_VERSION = 0x05
	METHOD_CODE   = 0x00
)

type Protocol interface {
	HandleHandshake(b []byte) ([]byte, error)
	SentHandshake(conn net.Conn) error
}

/*
*

	The localConn connects to the dstServer, and sends a ver
	identifier/method selection message:
	            +----+----------+----------+
	            |VER | NMETHODS | METHODS  |
	            +----+----------+----------+
	            | 1  |    1     | 1 to 255 |
	            +----+----------+----------+
	The VER field is set to X'05' for this ver of the protocol.  The
	NMETHODS field contains the number of method identifier octets that
	appear in the METHODS field.

*
*/
type ProtocolVersion struct {
	VER      uint8
	NMETHODS uint8
	METHODS  []uint8
}

func (s *ProtocolVersion) HandleHandshake(b []byte) ([]byte, error) {
	// ver , len, method
	n := len(b)
	if n < 3 {
		return nil, fmt.Errorf("proto err,1, %v", b)
	}
	s.VER = b[0] //ReadByte reads and returns a single byte
	if s.VER != SOCKS_VERSION {
		return nil, fmt.Errorf("proto err, ver=%v", s.VER)
	}
	s.NMETHODS = b[1]
	if n != int(2+s.NMETHODS) {
		return nil, fmt.Errorf("proto err,2, %v", b)
	}
	s.METHODS = b[2 : 2+s.NMETHODS]

	// link method, 2 use account and pwd
	useMethod := byte(0x00)
	for _, v := range s.METHODS {
		if v == METHOD_CODE {
			useMethod = METHOD_CODE
		}
	}

	if useMethod != METHOD_CODE {
		return nil, errors.New("proto method err")
	}
	resp := []byte{SOCKS_VERSION, useMethod}
	return resp, nil

}

func (s *ProtocolVersion) SentHandshake(conn net.Conn) error {
	resp := []byte{SOCKS_VERSION, 0x01, METHOD_CODE}
	conn.Write(resp)
	return nil
}

/*
   This begins with the client producing a
   Username/Password request:
   +----+------+----------+------+----------+
   |VER | ULEN |  UNAME   | PLEN |  PASSWD  |
   +----+------+----------+------+----------+
   | 1  |  1   | 1 to 255 |  1   | 1 to 255 |
   +----+------+----------+------+----------+

*/

type Socks5AuthUPasswd struct {
	VER    uint8
	ULEN   uint8
	UNAME  string
	PLEN   uint8
	PASSWD string
}

func (s *Socks5AuthUPasswd) HandleAuth(b []byte) ([]byte, error) {
	n := len(b)

	s.VER = b[0]
	if s.VER != 5 {
		return nil, errors.New("just support socket5")
	}

	s.ULEN = b[1]
	s.UNAME = string(b[2 : 2+s.ULEN])
	s.PLEN = b[2+s.ULEN+1]
	s.PASSWD = string(b[n-int(s.PLEN) : n])
	log.Println(s.UNAME, s.PASSWD)

	/**
	  The server verifies the supplied UNAME and PASSWD, and sends the
	  following response:

	                          +----+--------+
	                          |VER | STATUS |
	                          +----+--------+
	                          | 1  |   1    |
	                          +----+--------+

	  A STATUS field of X'00' indicates success. If the server returns a
	  `failure' (STATUS value other than X'00') status, it MUST close the
	  connection.
	*/
	resp := []byte{SOCKS_VERSION, 0x00}
	// conn.Write(resp)

	return resp, nil
}

/*
*

	+----+-----+-------+------+----------+----------+
	|VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	+----+-----+-------+------+----------+----------+
	| 1  |  1  | X'00' |  1   | Variable |    2     |
	+----+-----+-------+------+----------+----------+

*
*/
type Socks5Resolution struct {
	VER         uint8
	CMD         uint8
	RSV         uint8
	ATYP        uint8
	DSTADDR     []byte
	DSTPORT     uint16
	DSTDOMAIN   string
	RAWADDR     *net.TCPAddr
	DestAddrStr string
}

func (s *Socks5Resolution) LSTRequest(b []byte) ([]byte, error) {
	n := len(b)
	if n < 7 {
		return nil, errors.New("proto err")
	}
	s.VER = b[0]
	if s.VER != SOCKS_VERSION {
		return nil, errors.New("just support socket5")
	}

	s.CMD = b[1]
	if s.CMD != 1 {
		return nil, errors.New("proto err")
	}
	s.RSV = b[2] //RSV

	s.ATYP = b[3]

	switch s.ATYP {
	case 1:
		//	IP V4 address: X'01'
		s.DSTADDR = b[4 : 4+net.IPv4len]
	case 3:
		//	DOMAINNAME: X'03'
		s.DSTDOMAIN = string(b[5 : n-2])
		ipAddr, err := net.ResolveIPAddr("ip", s.DSTDOMAIN)
		if err != nil {
			return nil, err
		}
		s.DSTADDR = ipAddr.IP
	case 4:
		//	IP V6 address: X'04'
		s.DSTADDR = b[4 : 4+net.IPv6len]
	default:
		return nil, errors.New("ATYP err")
	}

	s.DSTPORT = binary.BigEndian.Uint16(b[n-2 : n])
	// DSTADDR
	s.RAWADDR = &net.TCPAddr{
		IP:   s.DSTADDR,
		Port: int(s.DSTPORT),
	}

	s.DestAddrStr = s.RAWADDR.String()
	if s.DSTDOMAIN != "" {
		s.DestAddrStr = s.DSTDOMAIN + ", " + s.DestAddrStr
	}

	/**
	  +----+-----+-------+------+----------+----------+
	  |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	  +----+-----+-------+------+----------+----------+
	  | 1  |  1  | X'00' |  1   | Variable |    2     |
	  +----+-----+-------+------+----------+----------+
	*/
	resp := []byte{SOCKS_VERSION, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	return resp, nil
}
