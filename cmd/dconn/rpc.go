package main

import (
	"net"
	"net/rpc"

	"github.com/yznts/dsh/pkg/ddb"
	"go.kyoto.codes/zen/v3/async"
)

type RpcKillProcessArgs struct {
	Pid   int
	Force bool
}

type Rpc struct{}

func (s *Rpc) QueryData(query string, res *ddb.Data) error {
	data, err := db.QueryData(query)
	if err != nil {
		return err
	}
	*res = *data
	return nil
}

func (s *Rpc) QueryTables(empty string, res *[]ddb.Table) error {
	tables, err := db.QueryTables()
	if err != nil {
		return err
	}
	*res = tables
	return nil
}

func (s *Rpc) QueryColumns(table string, res *[]ddb.Column) error {
	columns, err := db.QueryColumns(table)
	if err != nil {
		return err
	}
	*res = columns
	return nil
}

func (s *Rpc) QueryProcesses(empty string, res *[]ddb.Process) error {
	processes, err := db.QueryProcesses()
	if err != nil {
		return err
	}
	*res = processes
	return nil
}

func (s *Rpc) KillProcess(args RpcKillProcessArgs, res *bool) error {
	err := db.KillProcess(args.Pid, args.Force)
	if err != nil {
		*res = false
	}
	return err
}

func rpcserver(addr string) *async.Future[bool] {
	return async.New(func() (bool, error) {
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			panic(err)
		}
		server := &Rpc{}
		rpc.Register(server)
		rpc.Accept(ln)
		return false, nil
	})
}
