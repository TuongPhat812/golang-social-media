package scylla

import (
	"github.com/gocql/gocql"
)

func NewSession(hosts []string, keyspace string) (*gocql.Session, error) {
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.ProtoVersion = 4
	// Use Quorum for multi-node cluster (requires 2 out of 3 nodes to acknowledge)
	// This provides better consistency guarantees than ONE
	cluster.Consistency = gocql.Quorum
	cluster.ConnectTimeout = 0
	cluster.Timeout = 0

	return cluster.CreateSession()
}
