package util

import (
	"log"
	"net"
	"os"
	"runtime/debug"
	"strings"
)

// PanicRecover 有时候可能需要自己手动解决一些panic
func PanicRecover(stack bool) {
	if err := recover(); err != nil {
		// Check for a broken connection, as it is not really a
		// condition that warrants a panic stack trace.
		var brokenPipe bool
		if ne, ok := err.(*net.OpError); ok {
			if se, ok := ne.Err.(*os.SyscallError); ok {
				if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
					brokenPipe = true
				}
			}
		}

		if brokenPipe {
			log.Printf("[Recovery brokenPipe]: %+v \n\n", err)
			return
		}

		if stack {
			log.Printf("[Recovery from panic]: %+v \n %s \n", err, string(debug.Stack()))
		} else {
			log.Printf("[Recovery from panic]: %+v \n\n", err)
		}
	}
}
