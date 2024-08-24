package proxy

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jiftle/sckproxy/internal/proto"
	"github.com/jiftle/sckproxy/internal/utils"

	"github.com/gogf/gf/v2/frame/g"
)

func StartSocket5Proxy(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		g.Log().Warningf(context.Background(), "Error listening: %s", err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	g.Log().Infof(context.Background(), "socket5 proxy listen on %v", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			g.Log().Warningf(context.Background(), "Error accepting: %v", err.Error())
			os.Exit(1)
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	clientAddr := conn.RemoteAddr().String()

	reader := bufio.NewReader(conn)
	buf := make([]byte, 1024)
	nr, err := reader.Read(buf)
	if err != nil {
		g.Log().Warningf(context.Background(), "handshake read err,%v", err)
		return
	}

	var pt proto.ProtocolVersion

	resp, err := pt.HandleHandshake(buf[0:nr])
	if err != nil {
		g.Log().Warningf(context.Background(), "handshake err,%v", err)
		return
	}
	_, err = conn.Write(resp)
	if err != nil {
		g.Log().Warningf(context.Background(), "handshake write err,%v", err)
		return
	}

	for {
		buf = make([]byte, 1024*32)
		nr, err := reader.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			g.Log().Warningf(context.Background(), "handshake read err,%v", err)
			return
		}

		var request proto.Socks5Resolution
		resp, err = request.LSTRequest(buf[0:nr])
		if err != nil {
			g.Log().Warningf(context.Background(), "LST request err,%v", err)
			return
		}
		_, err = conn.Write(resp)
		if err != nil {
			g.Log().Warningf(context.Background(), "handshake write err,%v", err)
			return
		}
		g.Log().Infof(context.Background(), "%s accepted %s:%d[%s]", conn.RemoteAddr().String(), request.DSTDOMAIN, request.DSTPORT, request.RAWADDR.String())

		dstServer, err := net.DialTimeout("tcp", request.RAWADDR.String(), time.Second*3)
		if err != nil {
			g.Log().Warningf(context.Background(), "connect %s err,%s", request.RAWADDR.String(), err.Error())
			return
		}
		defer dstServer.Close()

		wg := new(sync.WaitGroup)
		wg.Add(2)

		go func() {
			defer wg.Done()
			n, err := utils.IoCopy(conn, dstServer)
			if err != nil {
				if strings.Contains(err.Error(), "connection reset by peer") {
					return
				} else if strings.Contains(err.Error(), "write: broken pipe") {
					return
				}
				g.Log().Warningf(context.Background(), "%v->%s, send fail,%v", clientAddr, request.RAWADDR.String(), err)
			} else {
				g.Log().Infof(context.Background(), "%v->%s,len=%s", clientAddr, request.RAWADDR.String(), utils.BytesSize2Str(n))
			}
		}()

		go func() {
			defer wg.Done()
			n, err := utils.IoCopy(dstServer, conn)
			if err != nil {
				if strings.Contains(err.Error(), "connection reset by peer") {
					return
				} else if strings.Contains(err.Error(), "write: broken pipe") {
					return
				}
				g.Log().Warningf(context.Background(), "%s->%v, send fail,%v", request.RAWADDR.String(), clientAddr, err)
			} else {
				g.Log().Infof(context.Background(), "%s->%v, ,len=%s", request.RAWADDR.String(), clientAddr, utils.BytesSize2Str(n))
			}
		}()
		wg.Wait()
	}
}
