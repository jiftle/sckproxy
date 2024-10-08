package proxy

import (
	"context"
	"net"
	"os"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/jiftle/sckproxy/internal/utils"
)

func StartHttpProxy(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		g.Log().Warningf(context.Background(), "Error listening: %s", err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	g.Log().Infof(context.Background(), "http proxy listen on %v", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			g.Log().Warningf(context.Background(), "Error accepting: %s", err.Error())
			os.Exit(1)
		}

		go handleHttpRequest(conn)
	}
}

func handleHttpRequest(conn net.Conn) {
	req, err := utils.NewHTTPRequest(&conn, 4094, false, nil)
	if err != nil {
		conn.Close()
		return
	}

	clientAddr := conn.RemoteAddr().String()
	address := req.Host

	dstServer, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		g.Log().Warningf(context.Background(), "connect %s err,%s", address, err.Error())
		return
	}
	defer dstServer.Close()

	if req.IsHTTPS() {
		req.HTTPSReply()
	} else {
		_, err = dstServer.Write(req.HeadBuf)
		if err != nil {
			g.Log().Warningf(context.Background(), "write %s err,%s", address, err.Error())
			return
		}
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		defer wg.Done()
		n, err := utils.IoCopy(conn, dstServer)
		if err != nil {
			g.Log().Warningf(context.Background(), "c->s, send err, %v", err)
		} else {
			g.Log().Infof(context.Background(), "%v->%v,len=%s", clientAddr, address, utils.BytesSize2Str(n))
		}
	}()

	go func() {
		defer wg.Done()
		n, err := utils.IoCopy(dstServer, conn)
		if err != nil {
			g.Log().Warningf(context.Background(), "s->c, send err, %v", err)
		} else {
			g.Log().Infof(context.Background(), "%v->%v,len=%s", address, clientAddr, utils.BytesSize2Str(n))
		}
	}()
	wg.Wait()

	g.Log().Warningf(context.Background(), "%v close", address)
}
