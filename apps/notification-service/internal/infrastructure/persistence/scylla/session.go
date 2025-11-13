package scylla

import (
	"github.com/gocql/gocql"
)

func NewSession(hosts []string, keyspace string) (*gocql.Session, error) {
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.ProtoVersion = 4
	cluster.Consistency = gocql.Quorum
	cluster.ConnectTimeout = 0
	cluster.Timeout = 0

	return cluster.CreateSession()
}
