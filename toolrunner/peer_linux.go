//go:build linux

package toolrunner

import (
	"net"

	"golang.org/x/sys/unix"
)

func peerGuardSupported() bool { return true }

// SO_PEERCRED reports a PID relative to the receiving process's PID namespace.
// A controller in a separate Docker PID namespace has no visible PID here;
// local Bash/PTY/background child processes do. Container PID namespaces must
// remain separate (no pid: host or shared pid namespace).
func peerOutsidePIDNamespace(conn net.Conn) bool {
	u, ok := conn.(*net.UnixConn)
	if !ok {
		return false
	}
	rc, err := u.SyscallConn()
	if err != nil {
		return false
	}
	allowed := false
	var sockErr error
	err = rc.Control(func(fd uintptr) {
		var cred *unix.Ucred
		cred, sockErr = unix.GetsockoptUcred(int(fd), unix.SOL_SOCKET, unix.SO_PEERCRED)
		allowed = sockErr == nil && cred.Pid == 0
	})
	return err == nil && sockErr == nil && allowed
}
