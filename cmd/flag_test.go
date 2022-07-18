package cmd

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestCmdFlags(t *testing.T) {
	expectPort := 9010
	expectProtocol := "udp"
	expectTimeout := 1000

	setArg("port", expectPort)
	setArg("protocol", expectProtocol)
	setArg("timeout", expectTimeout)

	flag := NewFlagCmd()

	if flag.Port() != expectPort {
		t.Errorf("expect read port from program args %d got %d", expectPort, flag.Port())
	}

	if flag.Protocol() != expectProtocol {
		t.Errorf("expect read protocol from program args %s got %s", expectProtocol, flag.Protocol())
	}

	if flag.ReadTimeout() != time.Duration(expectTimeout)*time.Second {
		t.Errorf("expect read timeout from program args %d (s) got %v", expectTimeout, flag.ReadTimeout())
	}
}

func setArg(argName string, value any) {
	os.Args = append(os.Args, fmt.Sprintf("-%s=%v", argName, value))
}
