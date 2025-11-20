package scylla

import (
	"strings"
	"time"

	"github.com/gocql/gocql"
	"golang-social-media/pkg/config"
	"golang-social-media/pkg/logger"
)

func NewSession(hosts []string, keyspace string) (*gocql.Session, error) {
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.ProtoVersion = 4

	// Configure consistency level from environment variable
	// Default: Quorum (production-ready, requires 2+ nodes)
	// Options: QUORUM, ONE, LOCAL_ONE, ALL, etc.
	consistencyLevel := getConsistencyLevel()
	cluster.Consistency = consistencyLevel

	// Set reasonable timeouts to avoid long bootstrap times
	// ConnectTimeout: time to establish initial connection
	// Timeout: time for query execution
	cluster.ConnectTimeout = 10 * time.Second
	cluster.Timeout = 5 * time.Second
	// Retry policy: retry up to 3 times with exponential backoff
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 3}

	logger.Component("scylla.session").
		Info().
		Strs("hosts", hosts).
		Str("keyspace", keyspace).
		Str("consistency", consistencyLevel.String()).
		Msg("scylla session configured")

	return cluster.CreateSession()
}

// getConsistencyLevel parses consistency level from environment variable
// Default: Quorum (production-ready for 3-node cluster)
func getConsistencyLevel() gocql.Consistency {
	level := strings.ToUpper(config.GetEnv("SCYLLA_CONSISTENCY_LEVEL", "QUORUM"))

	switch level {
	case "QUORUM":
		return gocql.Quorum
	case "ONE":
		return gocql.One
	case "LOCAL_ONE":
		return gocql.LocalOne
	case "ALL":
		return gocql.All
	case "LOCAL_QUORUM":
		return gocql.LocalQuorum
	case "EACH_QUORUM":
		return gocql.EachQuorum
	case "ANY":
		return gocql.Any
	case "TWO":
		return gocql.Two
	case "THREE":
		return gocql.Three
	default:
		logger.Component("scylla.session").
			Warn().
			Str("level", level).
			Msg("unknown consistency level, using QUORUM")
		return gocql.Quorum
	}
}
