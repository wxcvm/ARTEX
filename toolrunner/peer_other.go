//go:build !linux

package toolrunner

import "net"

// Fail closed when explicitly requested on a platform without Linux SO_PEERCRED.
func peerGuardSupported() bool { return false }
func peerOutsidePIDNamespace(_ net.Conn) bool { return false }
