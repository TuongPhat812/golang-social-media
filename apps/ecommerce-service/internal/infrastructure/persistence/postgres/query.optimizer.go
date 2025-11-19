package postgres

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// QueryOptimizer provides query analysis and optimization helpers
type QueryOptimizer struct {
	db *gorm.DB
}

// NewQueryOptimizer creates a new QueryOptimizer
func NewQueryOptimizer(db *gorm.DB) *QueryOptimizer {
	return &QueryOptimizer{db: db}
}

// AnalyzeQuery explains a query execution plan
func (q *QueryOptimizer) AnalyzeQuery(ctx context.Context, query *gorm.DB) (string, error) {
	var result struct {
		QueryPlan string `gorm:"column:QUERY PLAN"`
	}

	// Get the SQL and arguments
	sql := query.ToSQL(func(query *gorm.DB) *gorm.DB {
		return query
	})

	// Use EXPLAIN ANALYZE
	explainSQL := fmt.Sprintf("EXPLAIN ANALYZE %s", sql)
	if err := q.db.WithContext(ctx).Raw(explainSQL).Scan(&result).Error; err != nil {
		return "", err
	}

	return result.QueryPlan, nil
}

// GetTableStats returns statistics about a table
func (q *QueryOptimizer) GetTableStats(ctx context.Context, tableName string) (map[string]interface{}, error) {
	var stats struct {
		TableName      string
		RowCount       int64
		TableSize      string
		IndexesSize    string
		TotalSize      string
		IndexCount     int
	}

	query := `
		SELECT 
			tablename as table_name,
			(SELECT COUNT(*) FROM ` + tableName + `) as row_count,
			pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as total_size,
			pg_size_pretty(pg_relation_size(schemaname||'.'||tablename)) as table_size,
			pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename) - pg_relation_size(schemaname||'.'||tablename)) as indexes_size,
			(SELECT COUNT(*) FROM pg_indexes WHERE tablename = $1) as index_count
		FROM pg_tables
		WHERE tablename = $1
	`

	if err := q.db.WithContext(ctx).Raw(query, tableName).Scan(&stats).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"table_name":    stats.TableName,
		"row_count":     stats.RowCount,
		"table_size":    stats.TableSize,
		"indexes_size":  stats.IndexesSize,
		"total_size":    stats.TotalSize,
		"index_count":   stats.IndexCount,
	}, nil
}

// SuggestIndexes analyzes queries and suggests missing indexes
func (q *QueryOptimizer) SuggestIndexes(ctx context.Context, tableName string) ([]string, error) {
	// Query pg_stat_statements for slow queries (requires extension)
	query := `
		SELECT 
			DISTINCT ON (left(query, 100))
			query,
			mean_exec_time,
			calls
		FROM pg_stat_statements
		WHERE query LIKE '%' || $1 || '%'
		ORDER BY left(query, 100), mean_exec_time DESC
		LIMIT 10
	`

	var suggestions []string
	rows, err := q.db.WithContext(ctx).Raw(query, tableName).Rows()
	if err != nil {
		// pg_stat_statements might not be enabled
		return []string{"Enable pg_stat_statements extension for query analysis"}, nil
	}
	defer rows.Close()

	for rows.Next() {
		var queryText string
		var meanTime float64
		var calls int64
		if err := rows.Scan(&queryText, &meanTime, &calls); err != nil {
			continue
		}
		if meanTime > 100 { // Suggest index for queries taking > 100ms
			suggestions = append(suggestions, fmt.Sprintf("Consider adding index for query: %s (avg: %.2fms, calls: %d)", queryText[:100], meanTime, calls))
		}
	}

	return suggestions, nil
}

// VacuumAnalyze runs VACUUM ANALYZE on a table
func (q *QueryOptimizer) VacuumAnalyze(ctx context.Context, tableName string) error {
	sql := fmt.Sprintf("VACUUM ANALYZE %s", tableName)
	return q.db.WithContext(ctx).Exec(sql).Error
}

