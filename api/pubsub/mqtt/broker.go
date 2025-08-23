package mqtt

import (
	"fmt"
	"net"
	"net/url"
	"testing"

	"github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
)

func getFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "127.0.0.1:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			//nolint:errcheck
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}

// NewBroker creates a local MQTT broker that can be used for testing.
func NewBroker(t *testing.T) (*mqtt.Server, *url.URL) {
	server := mqtt.New(&mqtt.Options{InlineClient: true})

	err := server.AddHook(new(auth.AllowHook), nil)
	if err != nil {
		t.Fatalf("adding auth hook: %v", err)
	}

	port, err := getFreePort()
	if err != nil {
		t.Fatalf("getting free port: %v", err)
	}

	ws := listeners.NewWebsocket(listeners.Config{ID: "broker1", Address: fmt.Sprintf("127.0.0.1:%d", port)})
	err = server.AddListener(ws)
	if err != nil {
		t.Fatalf("adding ws listener: %v", err)
	}

	addr, err := url.Parse(fmt.Sprintf("ws://%s", ws.Address()))
	if err != nil {
		t.Fatalf("parsing broker url: %v", err)
	}

	return server, addr
}
