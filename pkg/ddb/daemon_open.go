//go:build daemon

package ddb

import (
	"net/rpc"
	"os/exec"
	"time"
)

func Open(dsn string) (Database, error) {
	// Start daemon
	cmd := exec.Command("dconn", "-dsn", dsn, "-rpc", "127.0.0.1:25123")
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	// Wait for daemon to start and open connection
	var (
		client *rpc.Client
		start  = time.Now()
	)
	for {
		// Check timeout
		if time.Since(start) > 5*time.Second {
			return nil, err
		}
		// Pause
		time.Sleep(10 * time.Millisecond)
		// Open connection
		client, err = rpc.Dial("tcp", "127.0.0.1:25123")
		if err != nil {
			continue
		}
		// Exit
		break
	}
	// Compose and return
	return &Rpc{client, cmd}, nil
}
