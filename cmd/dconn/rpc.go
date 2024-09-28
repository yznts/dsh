package main

import (
	"net"
	"net/rpc"

	"github.com/yznts/dsh/pkg/ddb"
	"go.kyoto.codes/zen/v3/async"
)

// RpcKillProcessArgs holds arguments for Rpc.KillProcess.
type RpcKillProcessArgs struct {
	Pid   int
	Force bool
}

// Rpc provides a set of RPC-compatible wrap methods
// around ddb.Database.
type Rpc struct{}

// QueryData is a wrap method around ddb.Database.QueryData.
func (s *Rpc) QueryData(query string, res *ddb.Data) error {
	data, err := db.QueryData(query)
	if err != nil {
		return err
	}
	*res = *data
	return nil
}

// QueryTables is a wrap method around ddb.Database.QueryTables.
func (s *Rpc) QueryTables(empty string, res *[]ddb.Table) error {
	tables, err := db.QueryTables()
	if err != nil {
		return err
	}
	*res = tables
	return nil
}

// QueryColumns is a wrap method around ddb.Database.QueryColumns.
func (s *Rpc) QueryColumns(table string, res *[]ddb.Column) error {
	columns, err := db.QueryColumns(table)
	if err != nil {
		return err
	}
	*res = columns
	return nil
}

// QueryProcesses is a wrap method around ddb.Database.QueryProcesses.
func (s *Rpc) QueryProcesses(empty string, res *[]ddb.Process) error {
	processes, err := db.QueryProcesses()
	if err != nil {
		return err
	}
	*res = processes
	return nil
}

// KillProcess is a wrap method around ddb.Database.KillProcess.
func (s *Rpc) KillProcess(args RpcKillProcessArgs, res *bool) error {
	err := db.KillProcess(args.Pid, args.Force)
	if err != nil {
		*res = false
	}
	return err
}

// rpcserver starts an RPC server on the given address.
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
