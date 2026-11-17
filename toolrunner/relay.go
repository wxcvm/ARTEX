package toolrunner

import (
	"context"
	"errors"
	"io"
	"net"
)

// Relay forwards only a single fixed TCP destination (the ARTEX recording
// proxy). It deliberately provides no general-purpose socket or HTTP API.
func Relay(ctx context.Context, listen, target string) error {
	if listen == "" || target == "" {
		return errors.New("proxy relay listen and target are required")
	}
	ln, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}
	defer ln.Close()
	go func() {
		<-ctx.Done()
		_ = ln.Close()
	}()
	for {
		in, err := ln.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return err
		}
		go func() {
			defer in.Close()
			out, err := net.Dial("tcp", target)
			if err != nil {
				return
			}
			defer out.Close()
			go func() {
				_, _ = io.Copy(out, in)
				if tcp, ok := out.(*net.TCPConn); ok {
					_ = tcp.CloseWrite()
				}
			}()
			_, _ = io.Copy(in, out)
		}()
	}
}
