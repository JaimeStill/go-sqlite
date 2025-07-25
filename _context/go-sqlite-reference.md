# Go Development Quick Reference: SQLite FTS5 & BM25

## Setup and Dependencies

### Required Dependencies

```go
// go.mod
module fts5-project

go 1.24

require (
    github.com/mattn/go-sqlite3 v1.14.17
    // Optional: dedicated FTS5 enabler
    github.com/knaka/go-sqlite3-fts5 v0.0.0-20230407203622-8aec9c2adfd2
)
```

### Build Configuration

```bash
# Enable FTS5 with build tags
go build -tags "fts5" ./cmd/app

# Or set CGO flags
export CGO_CFLAGS="-DSQLITE_ENABLE_FTS5"
go build ./cmd/app

# Verify FTS5 availability
go run -tags fts5 -c "PRAGMA compile_options;" | grep FTS5
```

### Basic Setup

```go
package main

import (
    "database/sql"
    "log"
    
    _ "github.com/mattn/go-sqlite3"
    // Optional: ensures FTS5 is available
    _ "github.com/knaka/go-sqlite3-fts5"
)

func main() {
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // Verify FTS5 availability
    if err := verifyFTS5Support(db); err != nil {
        log.Fatal("FTS5 not available:", err)
    }
}

func verifyFTS5Support(db *sql.DB) error {
    var available bool
    err := db.QueryRow(`
        SELECT COUNT(*) > 0 
        FROM pragma_compile_options 
        WHERE compile_options = 'ENABLE_FTS5'
    `).Scan(&available)
    
    if err != nil || !available {
        return fmt.Errorf("FTS5 not enabled")
    }
    return nil
}
```

## Database Patterns

### Connection Management

```go
type Database struct {
    db *sql.DB
}

func NewDatabase(dataSourceName string) (*Database, error) {
    // Essential connection parameters for FTS5
    dsn := fmt.Sprintf("%s?_journal=WAL&_timeout=5000&_fk=true", dataSourceName)
    
    db, err := sql.Open("sqlite3", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    
    // Configure connection pool
    db.SetMaxOpenConns(1)  // SQLite is single-writer
    db.SetMaxIdleConns(1)
    db.SetConnMaxLifetime(time.Hour)
    
    return &Database{db: db}, nil
}

func (d *Database) Close() error {
    return d.db.Close()
}
```

### Table Creation Patterns

```go
// Document represents a searchable document
type Document struct {
    ID      int64     `db:"id"`
    Title   string    `db:"title"`
    Content string    `db:"content"`
    Created time.Time `db:"created"`
}

func (d *Database) InitSchema(ctx context.Context) error {
    schemas := []string{
        // Base table
        `CREATE TABLE IF NOT EXISTS documents (
            id INTEGER PRIMARY KEY,
            title TEXT NOT NULL,
            content TEXT NOT NULL,
            created DATETIME DEFAULT CURRENT_TIMESTAMP
        )`,
        
        // FTS5 virtual table with external content
        `CREATE VIRTUAL TABLE IF NOT EXISTS documents_fts USING fts5(
            title, content,
            content='documents',
            content_rowid='id',
            tokenize='porter unicode61'
        )`,
        
        // Sync triggers
        `CREATE TRIGGER IF NOT EXISTS documents_after_insert 
         AFTER INSERT ON documents BEGIN
             INSERT INTO documents_fts(rowid, title, content) 
             VALUES (new.id, new.title, new.content);
         END`,
         
        `CREATE TRIGGER IF NOT EXISTS documents_after_update 
         AFTER UPDATE ON documents BEGIN
             INSERT INTO documents_fts(documents_fts, rowid, title, content) 
             VALUES('delete', old.id, old.title, old.content);
             INSERT INTO documents_fts(rowid, title, content) 
             VALUES (new.id, new.title, new.content);
         END`,
         
        `CREATE TRIGGER IF NOT EXISTS documents_after_delete 
         AFTER DELETE ON documents BEGIN
             INSERT INTO documents_fts(documents_fts, rowid, title, content) 
             VALUES('delete', old.id, old.title, old.content);
         END`,
    }
    
    for _, schema := range schemas {
        if _, err := d.db.ExecContext(ctx, schema); err != nil {
            return fmt.Errorf("failed to create schema: %w", err)
        }
    }
    
    return nil
}
```

## CRUD Operations

### Insert Operations

```go
func (d *Database) InsertDocument(ctx context.Context, doc *Document) error {
    query := `
        INSERT INTO documents (title, content, created) 
        VALUES (?, ?, ?)
        RETURNING id`
    
    err := d.db.QueryRowContext(ctx, query, 
        doc.Title, doc.Content, doc.Created).Scan(&doc.ID)
    
    if err != nil {
        return fmt.Errorf("failed to insert document: %w", err)
    }
    
    // FTS index is updated automatically via triggers
    return nil
}

func (d *Database) BatchInsertDocuments(ctx context.Context, docs []*Document) error {
    tx, err := d.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback()
    
    stmt, err := tx.PrepareContext(ctx, 
        `INSERT INTO documents (title, content, created) VALUES (?, ?, ?)`)
    if err != nil {
        return fmt.Errorf("failed to prepare statement: %w", err)
    }
    defer stmt.Close()
    
    for _, doc := range docs {
        if _, err := stmt.ExecContext(ctx, doc.Title, doc.Content, doc.Created); err != nil {
            return fmt.Errorf("failed to insert document: %w", err)
        }
    }
    
    return tx.Commit()
}
```

### Search Operations

```go
type SearchResult struct {
    Document
    Score     float64 `db:"score"`
    Rank      int     `db:"rank_num"`
    Highlight string  `db:"highlight"`
    Snippet   string  `db:"snippet"`
}

func (d *Database) SearchDocuments(ctx context.Context, query string, limit int) ([]*SearchResult, error) {
    // Escape search query to prevent FTS syntax errors
    escapedQuery := escapeSearchQuery(query)
    
    sql := `
        SELECT 
            d.id,
            highlight(documents_fts, 0, '<mark>', '</mark>') as title,
            snippet(documents_fts, 1, '<mark>', '</mark>', '...', 8) as content,
            d.created,
            bm25(documents_fts) as score,
            ROW_NUMBER() OVER (ORDER BY rank) as rank_num
        FROM documents d
        JOIN documents_fts ON documents_fts.rowid = d.id
        WHERE documents_fts MATCH ?
        ORDER BY rank
        LIMIT ?`
    
    rows, err := d.db.QueryContext(ctx, sql, escapedQuery, limit)
    if err != nil {
        return nil, fmt.Errorf("failed to execute search: %w", err)
    }
    defer rows.Close()
    
    var results []*SearchResult
    for rows.Next() {
        var result SearchResult
        err := rows.Scan(
            &result.ID,
            &result.Title,
            &result.Content,
            &result.Created,
            &result.Score,
            &result.Rank,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan result: %w", err)
        }
        results = append(results, &result)
    }
    
    return results, rows.Err()
}

func escapeSearchQuery(query string) string {
    // Escape quotes in FTS5 queries
    escaped := strings.ReplaceAll(query, `"`, `""`)
    // Wrap in quotes for phrase matching
    return `"` + escaped + `"`
}
```

### Advanced Search Patterns

```go
type SearchOptions struct {
    Query       string
    Limit       int
    Offset      int
    ColumnBoost map[string]float64
    FieldFilter string
}

func (d *Database) AdvancedSearch(ctx context.Context, opts SearchOptions) ([]*SearchResult, error) {
    var queryBuilder strings.Builder
    var args []interface{}
    
    // Build base query
    queryBuilder.WriteString(`
        SELECT 
            d.id,
            highlight(documents_fts, 0, '<mark>', '</mark>') as title,
            snippet(documents_fts, 1, '<mark>', '</mark>', '...', 8) as content,
            d.created,`)
    
    // Add scoring with optional column weights
    if len(opts.ColumnBoost) > 0 {
        titleWeight := opts.ColumnBoost["title"]
        contentWeight := opts.ColumnBoost["content"]
        if titleWeight == 0 {
            titleWeight = 1.0
        }
        if contentWeight == 0 {
            contentWeight = 1.0
        }
        
        queryBuilder.WriteString(fmt.Sprintf(
            `bm25(documents_fts, %.1f, %.1f) as score,`, 
            titleWeight, contentWeight))
    } else {
        queryBuilder.WriteString(`bm25(documents_fts) as score,`)
    }
    
    queryBuilder.WriteString(`
            ROW_NUMBER() OVER (ORDER BY rank) as rank_num
        FROM documents d
        JOIN documents_fts ON documents_fts.rowid = d.id
        WHERE documents_fts MATCH ?`)
    
    // Build search query with optional field filter
    searchQuery := opts.Query
    if opts.FieldFilter != "" {
        searchQuery = fmt.Sprintf("%s:%s", opts.FieldFilter, searchQuery)
    }
    
    args = append(args, escapeSearchQuery(searchQuery))
    
    queryBuilder.WriteString(` ORDER BY rank`)
    
    if opts.Limit > 0 {
        queryBuilder.WriteString(` LIMIT ?`)
        args = append(args, opts.Limit)
        
        if opts.Offset > 0 {
            queryBuilder.WriteString(` OFFSET ?`)
            args = append(args, opts.Offset)
        }
    }
    
    rows, err := d.db.QueryContext(ctx, queryBuilder.String(), args...)
    if err != nil {
        return nil, fmt.Errorf("failed to execute advanced search: %w", err)
    }
    defer rows.Close()
    
    var results []*SearchResult
    for rows.Next() {
        var result SearchResult
        err := rows.Scan(
            &result.ID,
            &result.Title,
            &result.Content,
            &result.Created,
            &result.Score,
            &result.Rank,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan result: %w", err)
        }
        results = append(results, &result)
    }
    
    return results, rows.Err()
}
```

## Testing Patterns

### Test Setup

```go
func setupTestDB(t *testing.T) *Database {
    db, err := NewDatabase(":memory:")
    if err != nil {
        t.Fatalf("Failed to create test database: %v", err)
    }
    
    ctx := context.Background()
    if err := db.InitSchema(ctx); err != nil {
        t.Fatalf("Failed to initialize schema: %v", err)
    }
    
    return db
}

func TestDocumentSearch(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    
    ctx := context.Background()
    
    // Insert test data
    testDocs := []*Document{
        {Title: "Go Programming", Content: "Learn Go programming language"},
        {Title: "SQLite Guide", Content: "Database operations with SQLite"},
        {Title: "FTS5 Tutorial", Content: "Full-text search with SQLite FTS5"},
    }
    
    for _, doc := range testDocs {
        if err := db.InsertDocument(ctx, doc); err != nil {
            t.Fatalf("Failed to insert test document: %v", err)
        }
    }
    
    // Test search
    results, err := db.SearchDocuments(ctx, "SQLite", 10)
    if err != nil {
        t.Fatalf("Search failed: %v", err)
    }
    
    if len(results) == 0 {
        t.Error("Expected search results, got none")
    }
    
    // Verify relevance ordering (lower scores = better)
    for i := 1; i < len(results); i++ {
        if results[i-1].Score > results[i].Score {
            t.Errorf("Results not ordered by relevance: %f > %f", 
                results[i-1].Score, results[i].Score)
        }
    }
}
```

### Benchmark Patterns

```go
func BenchmarkSearch(b *testing.B) {
    db := setupBenchmarkDB(b)
    defer db.Close()
    
    ctx := context.Background()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := db.SearchDocuments(ctx, "programming", 20)
        if err != nil {
            b.Fatalf("Search failed: %v", err)
        }
    }
}

func BenchmarkBatchInsert(b *testing.B) {
    db := setupBenchmarkDB(b)
    defer db.Close()
    
    docs := generateTestDocuments(1000)
    ctx := context.Background()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        if err := db.BatchInsertDocuments(ctx, docs); err != nil {
            b.Fatalf("Batch insert failed: %v", err)
        }
    }
}
```

## Error Handling

### Common Error Patterns

```go
var (
    ErrFTS5NotAvailable = errors.New("FTS5 extension not available")
    ErrInvalidQuery     = errors.New("invalid FTS5 query syntax")
    ErrDocumentNotFound = errors.New("document not found")
)

func (d *Database) handleFTSError(err error) error {
    if err == nil {
        return nil
    }
    
    errStr := err.Error()
    switch {
    case strings.Contains(errStr, "no such module: fts5"):
        return ErrFTS5NotAvailable
    case strings.Contains(errStr, "fts5: syntax error"):
        return ErrInvalidQuery
    case errors.Is(err, sql.ErrNoRows):
        return ErrDocumentNotFound
    default:
        return fmt.Errorf("database error: %w", err)
    }
}
```

### Graceful Degradation

```go
func (d *Database) SearchWithFallback(ctx context.Context, query string, limit int) ([]*SearchResult, error) {
    // Try FTS5 search first
    results, err := d.SearchDocuments(ctx, query, limit)
    if err == nil {
        return results, nil
    }
    
    // Fallback to LIKE search if FTS5 fails
    if errors.Is(err, ErrInvalidQuery) {
        return d.searchWithLike(ctx, query, limit)
    }
    
    return nil, err
}

func (d *Database) searchWithLike(ctx context.Context, query string, limit int) ([]*SearchResult, error) {
    likeQuery := "%" + query + "%"
    
    sql := `
        SELECT id, title, content, created, 0 as score, 0 as rank_num
        FROM documents 
        WHERE title LIKE ? OR content LIKE ?
        LIMIT ?`
    
    rows, err := d.db.QueryContext(ctx, sql, likeQuery, likeQuery, limit)
    if err != nil {
        return nil, fmt.Errorf("fallback search failed: %w", err)
    }
    defer rows.Close()
    
    var results []*SearchResult
    for rows.Next() {
        var result SearchResult
        err := rows.Scan(
            &result.ID, &result.Title, &result.Content, 
            &result.Created, &result.Score, &result.Rank)
        if err != nil {
            return nil, fmt.Errorf("failed to scan fallback result: %w", err)
        }
        results = append(results, &result)
    }
    
    return results, rows.Err()
}
```

## Performance Optimization

### Index Maintenance

```go
func (d *Database) OptimizeFTSIndex(ctx context.Context) error {
    // Optimize FTS5 index for better performance
    _, err := d.db.ExecContext(ctx, 
        `INSERT INTO documents_fts(documents_fts) VALUES('optimize')`)
    if err != nil {
        return fmt.Errorf("failed to optimize FTS index: %w", err)
    }
    
    return nil
}

func (d *Database) RebuildFTSIndex(ctx context.Context) error {
    // Completely rebuild FTS5 index
    _, err := d.db.ExecContext(ctx, 
        `INSERT INTO documents_fts(documents_fts) VALUES('rebuild')`)
    if err != nil {
        return fmt.Errorf("failed to rebuild FTS index: %w", err)
    }
    
    return nil
}
```

### Connection Pooling

```go
func NewDatabasePool(dataSourceName string) (*Database, error) {
    db, err := sql.Open("sqlite3", dataSourceName)
    if err != nil {
        return nil, err
    }
    
    // SQLite-specific pool settings
    db.SetMaxOpenConns(1)           // SQLite is single-writer
    db.SetMaxIdleConns(1)           
    db.SetConnMaxLifetime(time.Hour)
    db.SetConnMaxIdleTime(30 * time.Minute)
    
    return &Database{db: db}, nil
}
```

### Context and Timeout Handling

```go
func (d *Database) SearchWithTimeout(query string, timeout time.Duration) ([]*SearchResult, error) {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    
    return d.SearchDocuments(ctx, query, 50)
}
```

## Production Patterns

### Configuration Management

```go
type Config struct {
    DatabasePath    string        `env:"DATABASE_PATH" default:"./app.db"`
    MaxConnections  int           `env:"MAX_CONNECTIONS" default:"1"`
    QueryTimeout    time.Duration `env:"QUERY_TIMEOUT" default:"30s"`
    EnableWAL       bool          `env:"ENABLE_WAL" default:"true"`
    EnableForeignKeys bool        `env:"ENABLE_FK" default:"true"`
}

func (c Config) DSN() string {
    var params []string
    
    if c.EnableWAL {
        params = append(params, "_journal=WAL")
    }
    
    if c.EnableForeignKeys {
        params = append(params, "_fk=true")
    }
    
    params = append(params, fmt.Sprintf("_timeout=%d", 
        int(c.QueryTimeout.Milliseconds())))
    
    return fmt.Sprintf("%s?%s", c.DatabasePath, strings.Join(params, "&"))
}
```

### Logging and Monitoring

```go
func (d *Database) SearchDocumentsWithLogging(ctx context.Context, query string, limit int) ([]*SearchResult, error) {
    start := time.Now()
    
    results, err := d.SearchDocuments(ctx, query, limit)
    
    duration := time.Since(start)
    
    logger.Info("FTS search completed",
        "query", query,
        "results", len(results),
        "duration", duration,
        "error", err,
    )
    
    return results, err
}
```
