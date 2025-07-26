# Phase 1: SQLite FTS5 Foundation Learning Tool

A CLI application for learning SQLite FTS5 (Full-Text Search) fundamentals and BM25 scoring through hands-on experimentation.

## Project Overview

This educational tool demonstrates core FTS5 concepts through practical examples:

- FTS5 virtual table creation with unicode61 tokenizer
- Document insertion with automatic indexing
- Full-text search using MATCH operator with BM25 scoring
- Complete CRUD operations on FTS5 tables
- Error handling and validation patterns

## Learning Objectives

By using this tool, you will understand:

1. **FTS5 Virtual Tables**: How to create and configure FTS5 tables
2. **BM25 Scoring**: How SQLite implements BM25 ranking (negative scores, lower = better)
3. **Search Patterns**: Basic, category-filtered, and field-specific searches
4. **Index Management**: How FTS5 automatically maintains indexes during CRUD operations
5. **Error Handling**: Type-safe error patterns in Go applications

## Prerequisites

- **Go 1.24+** installed
- **SQLite with FTS5 support** (typically included in modern SQLite installations)
- Basic understanding of SQL and command-line interfaces

## Installation & Build

### 1. Clone and Navigate

```bash
cd /path/to/go-sqlite/src/01-foundation
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Run with FTS5 Support

```bash
# Run directly (recommended)
go run -tags fts5 ./fts5-foundation

# Or test basic functionality
go run -tags fts5 ./fts5-foundation --help
```

⚠️ **Important**: Always include the `-tags fts5` flag to enable FTS5 support in the SQLite driver.

## Usage Examples

### Quick Start

```bash
# 1. Create the FTS5 table
go run -tags fts5 ./fts5-foundation document create-table --database learning.db

# 2. Insert sample documents
go run -tags fts5 ./fts5-foundation document batch-insert --database learning.db

# 3. Search documents
go run -tags fts5 ./fts5-foundation document search "go programming" --database learning.db

# 4. List all documents
go run -tags fts5 ./fts5-foundation document list --database learning.db
```

### Document Management

#### Create Table

```bash
go run -tags fts5 ./fts5-foundation document create-table --database mydb.db
```

#### Insert Single Document

```bash
go run -tags fts5 ./fts5-foundation document insert \
  --title "My Document" \
  --content "Document content here..." \
  --category "example" \
  --database mydb.db
```

#### Insert Sample Documents

```bash
go run -tags fts5 ./fts5-foundation document batch-insert --database mydb.db
```

#### Update Document

```bash
go run -tags fts5 ./fts5-foundation document update 1 \
  --title "Updated Title" \
  --database mydb.db
```

#### Delete Document

```bash
go run -tags fts5 ./fts5-foundation document delete 1 --database mydb.db
```

### Search Operations

#### Basic Search

```bash
go run -tags fts5 ./fts5-foundation document search "database indexing" --database mydb.db
```

#### Search with BM25 Scores

```bash
go run -tags fts5 ./fts5-foundation document search "golang" --scores --database mydb.db
```

#### Category-Filtered Search

```bash
go run -tags fts5 ./fts5-foundation document search-category "programming" "database" --database mydb.db
```

#### Field-Specific Search

```bash
go run -tags fts5 ./fts5-foundation document search-field "Introduction" "title" --database mydb.db
```

#### Limit Results

```bash
go run -tags fts5 ./fts5-foundation document search "database" --limit 3 --database mydb.db
```

## Architecture

### Project Structure

```
fts5-foundation/
├── main.go                 # Application entry point
├── commands/              # CLI command definitions
│   ├── command_group.go   # Hierarchical command structure pattern
│   ├── root.go           # Root command and global flags
│   └── document.go       # Document-related sub-commands
├── config/               # Configuration management
│   └── config.go        # Global configuration with Viper integration
├── database/             # Database layer
│   └── database.go      # Global database instance and FTS5 operations
├── handlers/             # Business logic layer
│   └── document.go      # Stateless document handlers using global instances
├── models/              # Data structures
│   └── document.go     # Document, SearchResult, and DocumentInfo models
└── errors/             # Centralized error handling system
    └── errors.go      # Typed errors, sentinel errors, and display functions
```

### Design Patterns

**Command Pattern**: Hierarchical CLI structure using Cobra

```
fts5-foundation
└── document
    ├── create-table
    ├── insert
    ├── batch-insert
    ├── search
    ├── search-category
    ├── search-field
    ├── list
    ├── update
    └── delete
```

**Error Handling**: Type-safe sentinel error system

- `ErrValidation`: Input validation failures with detailed context
- `ErrDatabase`: Database operation errors with SQLite details
- `ErrFTS5`: FTS5-specific errors with build hints
- `ErrNotFound`: Resource not found errors with specific identifiers
- `ErrTransaction`: Transaction failures with rollback information
- `DisplayError()`: Automatic verbose/simple error display based on flags

**Layered Architecture**: Clear separation of concerns

- Commands: CLI interface using CommandGroup pattern for hierarchical organization
- Config: Global configuration management with factory functions
- Database: Global database instance with initialization in PersistentPreRun
- Handlers: Stateless business logic accessing global config and database instances
- Models: Domain-specific data structures (Document, SearchResult, DocumentInfo)
- Errors: Centralized error handling with type safety and user-friendly display

## FTS5 Key Concepts

### BM25 Scoring

SQLite FTS5 uses BM25 with these characteristics:

- **Negative scores**: Lower values indicate better matches (-1.5 ranks higher than -3.2)
- **Fixed parameters**: k1=1.2, b=0.75 (hardcoded, non-configurable)
- **Column weighting**: Use `bm25(table_name, weight1, weight2, ...)` for custom importance

### Virtual Tables

```sql
CREATE VIRTUAL TABLE documents USING fts5(
    title,           -- Document title for headline searches
    content,         -- Main document body content  
    category,        -- Document classification for filtering
    tokenize='unicode61 remove_diacritics 1'
);
```

### Search Syntax

- **Basic**: `documents MATCH 'golang programming'`
- **Field-specific**: `documents MATCH 'title:Introduction'`
- **Category filter**: `documents MATCH 'category:database AND indexing'`

## Flags and Options

### Global Flags

- `--database, -d`: Database file path (default: ":memory:")
- `--verbose, -v`: Show detailed error information
- `--format, -f`: Output format (text, json) (default: "text")
- `--config`: Configuration file path
- `--help, -h`: Show help information

### Command-Specific Flags

- **insert**: `--title`, `--content`, `--category`
- **search**: `--limit`, `--scores`
- **update**: `--title`, `--content`, `--category`
- **list**: `--limit`

## Error Handling

The application provides clear, categorized error messages:

```bash
# Validation errors
Validation Error: validation failed: title cannot be empty

# Database errors  
Database Error: database operation failed: no such table: documents

# FTS5 errors
FTS5 Error: fts5 operation failed: SQLite not compiled with FTS5 support
Hint: Ensure SQLite is compiled with FTS5 support (go build -tags fts5)

# Not found errors
Not Found: not found: document with rowid 999
```

Use `--verbose` flag for detailed error chains:

```bash
go run -tags fts5 ./fts5-foundation document search "" --verbose
# Shows full error context and stack traces
```

## Development

### Running Tests

```bash
go test -tags fts5 ./fts5-foundation/...
```

### Adding New Commands

1. Define command in appropriate file under `commands/`
2. Add to command group in `setupCommands()`
3. Implement handler functions in `handlers/`
4. Use typed errors from `errors/` package

## Troubleshooting

### FTS5 Not Available

```
FTS5 Error: SQLite not compiled with FTS5 support
```

**Solution**: Ensure you're using the `-tags fts5` build flag.

### Table Not Found

```
Database Error: no such table: documents
```

**Solution**: Run `create-table` command first, or ensure you're using the same database file.

### Permission Denied

```
Database Error: unable to open database file
```

**Solution**: Check file permissions and ensure the directory exists.

## Learning Resources

- [SQLite FTS5 Documentation](https://www.sqlite.org/fts5.html)
- [BM25 Algorithm Explanation](https://en.wikipedia.org/wiki/Okapi_BM25)
- [Go SQLite Driver Documentation](https://github.com/mattn/go-sqlite3)

## Next Steps

This foundation project prepares you for advanced FTS5 topics:

- Custom tokenizers and ranking functions
- Multi-table search and joins
- Performance optimization techniques
- Integration with larger applications

---

**Phase 1 Learning Goal**: Master FTS5 fundamentals through hands-on CLI experimentation.
