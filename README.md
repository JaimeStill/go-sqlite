# SQLite FTS5 & BM25 Learning Project

A hands-on learning project to build understanding of SQLite FTS5 and BM25 scoring through incremental Go experiments.

## Project Overview

This repository contains a series of isolated learning phases, each focusing on specific SQLite FTS5 and BM25 concepts. The goal is to build foundational understanding through practical experimentation rather than theoretical study.

## Repository Structure

```
├── src/                          # All learning phases and utilities
│   ├── 00-setup-validation/      # ✅ Environment validation tool (COMPLETED)
│   ├── 01-foundation/            # ✅ Basic FTS5 setup and operations (COMPLETED)
│   ├── 02-bm25-fundamentals/     # ✅ BM25 scoring and interpretation (COMPLETED)
│   ├── 03-query-operations/      # Advanced FTS5 query patterns
│   ├── 04-ranking-relevance/     # Custom ranking strategies
│   ├── 05-advanced-features/     # FTS5 auxiliary functions
│   └── 06-integration-patterns/  # Production-ready patterns
└── _context/                     # Reference materials and documentation
```

## Quick Start

### 1. Prerequisites

- Go 1.24+
- SQLite with FTS5 support installed locally
- See `_context/sqlite-fts5-reference.md` for installation instructions

### 2. Validate Your Setup

Before starting any learning phases, validate your environment:

```bash
cd src/00-setup-validation
go run -tags fts5 ./setup-validation validation validate
```

This ensures SQLite FTS5 is properly installed and accessible from Go.

### 3. Explore Learning Phases

Each phase is a standalone project with its own README:

```bash
# Phase 0: Setup Validation (COMPLETED)
cd src/00-setup-validation
go run -tags fts5 ./setup-validation --help

# Phase 1: Foundation (COMPLETED)
cd src/01-foundation
go run -tags fts5 ./fts5-foundation --help
go run -tags fts5 ./fts5-foundation document --help

# Phase 2: BM25 Fundamentals (COMPLETED)
cd src/02-bm25-fundamentals
go run -tags fts5 ./bm25-fundamentals --help
go run -tags fts5 ./bm25-fundamentals corpus --help

# See comprehensive usage examples in each phase's README
```

## Key Learning Concepts

- **FTS5 Virtual Tables**: Full-text search indexing and querying
- **BM25 Scoring**: Relevance ranking algorithm and SQLite's implementation
- **Query Patterns**: Boolean operators, phrase matching, proximity search
- **Performance Optimization**: Index tuning and query optimization
- **Go Integration**: Building FTS5-enabled applications with proper build tags
- **CLI Architecture**: Hierarchical command patterns and error handling strategies

## Important Notes

### Build Tags Required

All Go programs in this repository require the FTS5 build tag:

```bash
# Correct - using subdirectory structure
go run -tags fts5 ./setup-validation [command]
go run -tags fts5 ./fts5-foundation [command]
go run -tags fts5 ./bm25-fundamentals [command]

# Wrong - will fail without FTS5 tags
go run ./setup-validation [command]
```

### Phase Isolation

Each learning phase is completely isolated with its own:

- `go.mod` file for dependencies at the phase root
- CLI application in a subdirectory (e.g., `setup-validation/`, `fts5-foundation/`, `bm25-fundamentals/`)
- `main.go` CLI entry point using Cobra/Viper with CommandGroup pattern
- `README.md` with phase-specific instructions and learning objectives
- Self-contained utilities (no shared dependencies between phases)

## Development Philosophy

- **Incremental Learning**: Start simple, add complexity gradually
- **Hands-on Experimentation**: Working code over theoretical concepts
- **Phase Independence**: Each phase stands alone and can be run separately
- **Native Development**: Local SQLite installation for consistency

## Getting Help

- Check individual phase README files for specific instructions
- Review `_context/sqlite-fts5-reference.md` for SQLite FTS5 setup
- See `CLAUDE.md` for detailed project guidelines and development workflow

Start with phase 00 (setup validation) and progress through the numbered phases in order to build your SQLite FTS5 expertise systematically.
