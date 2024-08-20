//go:build daemon

package ddb

import (
	"net/rpc"
	"os/exec"
)

type RpcKillProcessArgs struct {
	Pid   int
	Force bool
}

type Rpc struct {
	*rpc.Client
	*exec.Cmd
}

func (c *Rpc) QueryData(query string) (*Data, error) {
	res := &Data{}
	err := c.Call("Rpc.QueryData", query, res)
	return res, err
}

func (c *Rpc) QueryTables() ([]Table, error) {
	res := &[]Table{}
	err := c.Call("Rpc.QueryTables", "", res)
	return *res, err
}

func (c *Rpc) QueryColumns(table string) ([]Column, error) {
	res := &[]Column{}
	err := c.Call("Rpc.QueryColumns", table, res)
	return *res, err
}

func (c *Rpc) QueryProcesses() ([]Process, error) {
	res := &[]Process{}
	err := c.Call("Rpc.QueryProcesses", "", res)
	return *res, err
}

func (c *Rpc) KillProcess(pid int, force bool) error {
	err := c.Call("Rpc.KillProcess", RpcKillProcessArgs{pid, force}, nil)
	return err
}

func (c *Rpc) Close() error {
	// Close the connection
	c.Client.Close()
	// Kill the process
	return c.Cmd.Process.Kill()
}
