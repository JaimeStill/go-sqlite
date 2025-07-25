# SQLite FTS5 Setup Validation

A comprehensive validation tool to ensure your SQLite FTS5 environment is properly configured for the go-sqlite learning project.

## Learning Objectives

This validation project teaches and demonstrates:

- **Native SQLite FTS5 Setup**: Verifying local SQLite installation includes FTS5 support
- **Go Build Tags**: Understanding when and how to use `-tags fts5` with the mattn/go-sqlite3 driver  
- **BM25 Scoring**: Confirming BM25 relevance ranking works with SQLite's inverted scoring system
- **Error Handling**: Robust validation with clear error messages and recovery guidance
- **Build System Design**: Creating foolproof build configurations that prevent common mistakes

## Key Concepts

### SQLite FTS5 Integration

FTS5 (Full-Text Search version 5) is SQLite's advanced full-text search engine that provides:

- Virtual table interface for full-text indexing
- BM25 relevance scoring algorithm  
- Advanced query operators (AND, OR, NOT, phrase matching)
- Configurable tokenizers and ranking functions

### Go Build Tags Requirement

The `github.com/mattn/go-sqlite3` driver requires explicit FTS5 compilation:

- **Without tags**: `go run *.go` → FTS5 unavailable, "no such module" errors
- **With tags**: `go run -tags fts5 *.go` → FTS5 enabled and functional

### BM25 Scoring Behavior

SQLite FTS5 uses inverted BM25 scores where **lower (more negative) scores indicate better matches**:

- Score -1.5 ranks higher than -3.2
- Default parameters: k1=1.2, b=0.75 (hardcoded, non-configurable)
- Use `ORDER BY bm25(table)` or `ORDER BY rank` for relevance sorting

## Usage Instructions

### Recommended Usage

Run validation directly with the required FTS5 build tag:

```bash
# Quick validation (recommended)
go run -tags fts5 *.go validate

# Build standalone binary
go build -tags fts5 -o setup-validator *.go
./setup-validator validate
```

### Important Build Requirements

The project requires the FTS5 build tag for proper SQLite integration:

```bash
# Correct - includes FTS5 support
go run -tags fts5 *.go validate
go build -tags fts5 -o validator *.go

# Wrong - will fail with helpful error message
go run *.go validate
go build *.go
```

### Available Commands

```bash
go run -tags fts5 *.go validate    # Run validation
go run -tags fts5 *.go connect     # Test database connection
go run -tags fts5 *.go fts5        # Test FTS5 functionality  
go run -tags fts5 *.go testdata    # Test sample data generation
go run -tags fts5 *.go bm25        # Test BM25 scoring
```

## Validation Checks

The tool performs comprehensive validation:

1. **Build Check**: Verifies FTS5 build tag was included at compile time
2. **Connection Test**: Confirms SQLite database connectivity  
3. **FTS5 Support**: Tests virtual table creation and basic functionality
4. **BM25 Scoring**: Validates relevance ranking with proper score ordering

### Sample Output

```
SQLite FTS5 Setup Validation
============================

✅ Build Check: FTS5 support compiled in
✅ Connection: SQLite database opened successfully  
✅ FTS5 Support: FTS5 is available and functional
✅ BM25 Scoring: BM25 ranking works correctly

🎉 All validation checks passed!
Your SQLite FTS5 environment is ready for development.
```

## Learning Reflections

### Architecture Decisions

**Native vs Containerized Development**: Initially planned with Docker, but SQLite's embedded nature made native installation simpler and more appropriate for learning objectives. This teaches real-world SQLite deployment patterns.

**Build Tag Management**: The challenge of FTS5 build tags led to implementing comprehensive build tooling. This demonstrates the importance of build system design in Go projects with C dependencies.

**Error-First Design**: The validation tool fails fast with helpful error messages, teaching defensive programming patterns and user experience considerations.

### Key Insights

1. **SQLite FTS5 Availability**: Most modern SQLite distributions include FTS5, but Go drivers require explicit compilation flags
2. **BM25 Implementation**: SQLite's inverted scoring system differs from traditional search engines - understanding this prevents ranking confusion
3. **Build System Complexity**: Simple Go programs can require sophisticated build tooling when integrating with C libraries
4. **Validation Strategy**: Comprehensive upfront validation prevents issues in subsequent development phases

### Trade-offs Considered

- **Simplicity vs Robustness**: Chose robust build system over simple `go run` to prevent common FTS5 mistakes
- **Error Messages**: Verbose error output trades brevity for developer guidance and faster problem resolution  
- **Cross-Platform Support**: Added complexity for Windows/macOS support despite Linux-focused development

## Technical Details

- **Go Version**: 1.24+
- **SQLite Driver**: github.com/mattn/go-sqlite3 v1.14.24
- **CLI Framework**: Cobra + Viper for consistent command patterns
- **Build Requirements**: CGO_ENABLED=1, FTS5 build tags
- **Test Coverage**: Integration tests for FTS5 functionality and BM25 scoring

This validation project establishes a solid foundation for subsequent learning phases while teaching essential concepts about SQLite FTS5 integration, Go build systems, and robust validation design patterns.
