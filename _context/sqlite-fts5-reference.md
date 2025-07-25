# SQLite FTS5 Quick Reference & Best Practices

## Installation

### Linux (Ubuntu/Debian)

```bash
# Install SQLite3 with FTS5 support (included by default in most distributions)
sudo apt update
sudo apt install sqlite3 libsqlite3-dev

# Verify installation and FTS5 support
sqlite3 --version
sqlite3 ':memory:' 'PRAGMA compile_options;' | grep -i fts5
```

### Linux (RHEL/CentOS/Rocky/Fedora)

```bash
# For RHEL/CentOS/Rocky
sudo yum install sqlite sqlite-devel
# OR for newer versions
sudo dnf install sqlite sqlite-devel

# For Fedora
sudo dnf install sqlite sqlite-devel

# Verify installation
sqlite3 --version
sqlite3 ':memory:' 'PRAGMA compile_options;' | grep -i fts5
```

### Linux (Arch)

```bash
# Install SQLite (FTS5 enabled by default)
sudo pacman -S sqlite

# Verify installation
sqlite3 --version
sqlite3 ':memory:' 'PRAGMA compile_options;' | grep -i fts5
```

### macOS (Homebrew)

```bash
# Install SQLite via Homebrew (includes FTS5)
brew install sqlite

# Verify installation and FTS5 support
sqlite3 --version
sqlite3 ':memory:' 'PRAGMA compile_options;' | grep -i fts5

# Note: macOS comes with SQLite, but it may not have FTS5
# Use 'which sqlite3' to verify you're using the Homebrew version
which sqlite3  # Should show /usr/local/bin/sqlite3 or /opt/homebrew/bin/sqlite3
```

### macOS (MacPorts)

```bash
# Install SQLite via MacPorts
sudo port install sqlite3

# Verify installation
sqlite3 --version
sqlite3 ':memory:' 'PRAGMA compile_options;' | grep -i fts5
```

### Build from Source (if FTS5 not available)

```bash
# Download SQLite source
wget https://sqlite.org/2024/sqlite-autoconf-3450000.tar.gz
tar xzf sqlite-autoconf-3450000.tar.gz
cd sqlite-autoconf-3450000

# Configure with FTS5 enabled
./configure --enable-fts5 --prefix=/usr/local
# OR with more optimization flags
./configure --enable-fts5 --enable-rtree --enable-json1 --prefix=/usr/local

# Compile and install
make
sudo make install

# Verify installation
/usr/local/bin/sqlite3 --version
/usr/local/bin/sqlite3 ':memory:' 'PRAGMA compile_options;' | grep -i fts5
```

### Verification Commands

After installation, verify FTS5 is working correctly:

```bash
# Check SQLite version
sqlite3 --version

# Verify FTS5 is compiled in
sqlite3 ':memory:' 'PRAGMA compile_options;' | grep -i fts5
# Should output: ENABLE_FTS5

# Test FTS5 functionality
sqlite3 ':memory:' << 'EOF'
.echo on
SELECT 'SQLite Version: ' || sqlite_version();
SELECT 'FTS5 Support: ' || (
    CASE WHEN EXISTS(
        SELECT 1 FROM pragma_compile_options 
        WHERE compile_options = 'ENABLE_FTS5'
    ) THEN 'YES' ELSE 'NO' END
);

-- Test basic FTS5 functionality
CREATE VIRTUAL TABLE test_fts USING fts5(content);
INSERT INTO test_fts(content) VALUES ('Hello FTS5 world');
SELECT * FROM test_fts WHERE test_fts MATCH 'Hello';
SELECT 'Test passed: FTS5 is working correctly';
.quit
EOF
```

### Go Development Setup

Once SQLite with FTS5 is installed, configure your Go environment:

```bash
# Verify Go can build with SQLite FTS5
export CGO_ENABLED=1

# Test building with FTS5 tags (in a Go project directory)
go build -tags fts5 .

# If you encounter linking issues, you may need to set CGO flags
export CGO_CFLAGS="-I/usr/local/include"
export CGO_LDFLAGS="-L/usr/local/lib"
```

### Troubleshooting

**FTS5 not found:**

- Check if your distribution's SQLite package includes FTS5
- Consider building from source with `--enable-fts5`
- On macOS, ensure you're using Homebrew/MacPorts version, not system SQLite

**Go build issues:**

- Ensure `CGO_ENABLED=1` is set
- Use `-tags fts5` when building Go programs
- Check that sqlite-dev/sqlite-devel packages are installed

**Version conflicts:**

- Use `which sqlite3` to verify which SQLite binary you're using
- Consider using absolute paths if multiple versions are installed
- Check `LD_LIBRARY_PATH` and `PKG_CONFIG_PATH` for library conflicts

---

## Essential Commands

### Table Creation

```sql
-- Basic FTS5 table creation
CREATE VIRTUAL TABLE documents USING fts5(title, content);

-- With external content table (space optimization)
CREATE VIRTUAL TABLE docs_fts USING fts5(
    title, content, 
    content='documents', 
    content_rowid='id'
);

-- Custom tokenizer configuration
CREATE VIRTUAL TABLE articles USING fts5(
    content, 
    tokenize='porter unicode61'
);

-- Prefix indexing for autocomplete
CREATE VIRTUAL TABLE search_index USING fts5(
    title, body, 
    prefix='2 3 4'
);
```

### Query Patterns

```sql
-- Basic full-text search
SELECT * FROM documents WHERE documents MATCH 'search term';

-- Column-specific search
SELECT * FROM documents WHERE documents MATCH 'title:sqlite';

-- Boolean operators
SELECT * FROM documents WHERE documents MATCH 'sqlite AND fts5';
SELECT * FROM documents WHERE documents MATCH 'sqlite OR database';
SELECT * FROM documents WHERE documents MATCH 'sqlite NOT tutorial';

-- Phrase matching
SELECT * FROM documents WHERE documents MATCH '"exact phrase"';

-- Proximity search (terms within 10 positions)
SELECT * FROM documents WHERE documents MATCH 'NEAR(sqlite fts5, 10)';

-- Prefix search
SELECT * FROM documents WHERE documents MATCH 'program*';
```

### Ranking and Sorting

```sql
-- Sort by relevance (best practice)
SELECT * FROM documents 
WHERE documents MATCH 'query' 
ORDER BY rank;

-- Using BM25 scoring explicitly
SELECT *, bm25(documents) as score 
FROM documents 
WHERE documents MATCH 'query' 
ORDER BY bm25(documents);

-- Custom column weights (title=2x, content=1x)
SELECT *, bm25(documents, 2.0, 1.0) as weighted_score 
FROM documents 
WHERE documents MATCH 'query' 
ORDER BY weighted_score;
```

### Auxiliary Functions

```sql
-- Highlight search terms
SELECT highlight(documents, 0, '<mark>', '</mark>') as title
FROM documents WHERE documents MATCH 'query';

-- Extract snippets with context
SELECT snippet(documents, 1, '<mark>', '</mark>', '...', 10) as excerpt
FROM documents WHERE documents MATCH 'query';

-- Combined search with highlighting
SELECT 
    highlight(docs_fts, 0, '\u2039', '\u203a') as title,
    snippet(docs_fts, 1, '\u2039', '\u203a', '', 8) as content,
    bm25(docs_fts) as relevance
FROM docs_fts 
WHERE docs_fts MATCH ? 
ORDER BY rank;
```

## Configuration Options

### Tokenizers

- **unicode61** (default): Best for general multilingual text
- **porter**: Adds stemming (run/runs/running \u2192 run)
- **ascii**: Restrictive, ASCII-only tokenization
- **trigram**: For LIKE/GLOB pattern matching support

### Advanced Options

```sql
-- Disable content storage (index only)
CREATE VIRTUAL TABLE search_only USING fts5(
    title, content, 
    content=''
);

-- Disable columnsize for space savings
CREATE VIRTUAL TABLE compact_fts USING fts5(
    content, 
    columnsize=0
);

-- Custom rank function
CREATE VIRTUAL TABLE custom_rank USING fts5(
    content, 
    rank='bm25(10.0, 5.0)'
);
```

## Best Practices

### Schema Design

- **Use FTS5** over FTS3/FTS4 for new projects
- **No column types** in CREATE VIRTUAL TABLE statements
- **External content tables** for space optimization with large documents
- **Prefix indexing** only when needed (increases index size)

### Query Optimization

- **Always use MATCH** operator, never LIKE with FTS tables
- **ORDER BY rank** for relevance-sorted results
- **Column filters** (title:term) for better precision
- **Escape user input** in queries to prevent syntax errors

### Performance Guidelines

- **Batch inserts** in transactions for better performance
- **OPTIMIZE command** after bulk operations
- **Triggers for sync** when using external content tables
- **Monitor index size** - can be larger than source data

### Index Maintenance

```sql
-- Optimize FTS index after bulk operations
INSERT INTO documents_fts(documents_fts) VALUES('optimize');

-- Rebuild index completely
INSERT INTO documents_fts(documents_fts) VALUES('rebuild');

-- Check index integrity
INSERT INTO documents_fts(documents_fts) VALUES('integrity-check');
```

### External Content Pattern

```sql
-- Create base table
CREATE TABLE articles (
    id INTEGER PRIMARY KEY,
    title TEXT,
    content TEXT,
    created DATETIME
);

-- Create FTS index
CREATE VIRTUAL TABLE articles_fts USING fts5(
    title, content, 
    content='articles', 
    content_rowid='id'
);

-- Sync triggers
CREATE TRIGGER articles_after_insert 
AFTER INSERT ON articles BEGIN
    INSERT INTO articles_fts(rowid, title, content) 
    VALUES (new.id, new.title, new.content);
END;

CREATE TRIGGER articles_after_update 
AFTER UPDATE ON articles BEGIN
    INSERT INTO articles_fts(articles_fts, rowid, title, content) 
    VALUES('delete', old.id, old.title, old.content);
    INSERT INTO articles_fts(rowid, title, content) 
    VALUES (new.id, new.title, new.content);
END;

CREATE TRIGGER articles_after_delete 
AFTER DELETE ON articles BEGIN
    INSERT INTO articles_fts(articles_fts, rowid, title, content) 
    VALUES('delete', old.id, old.title, old.content);
END;
```

## Common Pitfalls

### Avoid These Mistakes

- **Using LIKE** instead of MATCH on FTS tables
- **Adding column types** in virtual table creation
- **Forgetting ORDER BY rank** for relevance sorting
- **Not escaping quotes** in search queries
- **Ignoring trigger maintenance** for external content

### Error Handling

```sql
-- Check if FTS5 is available
PRAGMA compile_options; -- Look for ENABLE_FTS5

-- Verify virtual table exists
PRAGMA table_info(table_name);

-- Debug query execution
EXPLAIN QUERY PLAN SELECT * FROM docs WHERE docs MATCH 'term';
```

### Query Escaping

```sql
-- Escape quotes in search terms
SELECT * FROM docs WHERE docs MATCH '"user ""query"" here"';

-- Handle special characters
-- Wrap entire query in quotes for phrase search
-- Use double quotes to escape internal quotes
```

## Performance Characteristics

### Speed Comparisons

- **FTS5 vs LIKE**: ~50x faster for text search operations
- **Index size**: Typically 40-60% of original text size
- **Memory usage**: Efficient for datasets up to millions of documents
- **Build time**: Fast incremental updates, slower bulk rebuilds

### Scaling Guidelines

- **Small datasets** (<10k docs): All approaches work well
- **Medium datasets** (10k-1M docs): External content recommended
- **Large datasets** (>1M docs): Consider columnsize=0, optimize regularly
- **Real-time updates**: Batch operations in transactions when possible

### Storage Optimization

- **contentless tables** save space but limit functionality
- **columnsize=0** disables document length tracking (breaks BM25)
- **external content** reduces duplication while maintaining features
- **prefix indexing** significantly increases index size

## Version Compatibility

### SQLite Requirements

- **Minimum version**: SQLite 3.20.0 (2017) for stable FTS5
- **Recommended**: SQLite 3.43.0+ (2023) for latest features
- **Build requirement**: SQLITE_ENABLE_FTS5 compile flag
- **Default availability**: Enabled in most SQLite distributions

### Feature Evolution

- **3.20.0**: FTS5 stable release
- **3.40.0**: Performance improvements
- **3.43.0**: Enhanced auxiliary functions
- **3.46.0**: Improved phrase handling
- **3.50.0**: Latest optimizations and bug fixes
