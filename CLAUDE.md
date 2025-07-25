# SQLite FTS5 & BM25 Contextual Memory System Project

## Project Overview

- **Goal**: Build understanding of SQLite FTS5 and BM25 scoring through iterative experimental projects in Go
- **Approach**: Start with lowest-level features, expand complexity incrementally
- **Learning Focus**: Hands-on experimentation rather than theoretical depth upfront
- **Future Target**: Contextual memory system implementation

## Repository Restrictions

- Any folder prefixed with a `_` is read-only for you. Unless explicitly directed by me or one of your directives in this document, you are only allowed to read these files.
- Any folder prefixed with `.` is private and not accessible by you. Unless explicitly directed by me or one of your directives in this document, you are not allowed to access or modify these files.

## Development Philosophy

- **Incremental Learning**: Build small, focused experiments that demonstrate single concepts
- **Minimal Viable Examples**: Prefer working code over comprehensive features initially
- **Context Scoping**: Keep each experiment bounded to avoid context drift
- **Feature-First**: Use latest supported SQLite FTS5 features from the start

## Go Environment Setup

- **Go Version**: Go 1.24 (current installation)
- **SQLite Driver**: `github.com/mattn/go-sqlite3` with FTS5 enabled
- **Build Tags**: Always compile with `go build -tags "fts5"` for FTS5 support
- **Database Mode**: Prefer in-memory databases (`:memory:`) for experiments unless persistence needed
- **Native SQLite**: Ensure SQLite with FTS5 support is installed locally (see `_context/sqlite-fts5-reference.md`)

## SQLite FTS5 Standards

- **Table Creation**: Always use `CREATE VIRTUAL TABLE name USING fts5(columns)` syntax
- **BM25 Usage**: Remember SQLite FTS5 BM25 returns negative scores (lower = better match)
- **Scoring**: Use `ORDER BY rank` or `ORDER BY bm25(table_name)` for relevance sorting
- **Query Syntax**: Follow FTS5 MATCH operator patterns, not LIKE patterns
- **Tokenizers**: Default to 'unicode61' unless specific language requirements exist

## Project Structure

- **Phase Isolation**: Each phase lives in `src/XX-phase-name/` as a standalone project
- **CLI Applications**: All phases use Cobra/Viper for consistent command-line interfaces
- **Shared Utilities**: Common functionality available in `src/shared/` package
- **Native Development**: Local SQLite FTS5 installation for consistent development environment
- **Data Persistence**: Shared `data/` directory for cross-phase experiments

## Code Structure Preferences

- **Error Handling**: Always handle SQLite errors explicitly, never ignore them
- **Resource Management**: Use defer for database cleanup in every function
- **CLI Design**: Use Cobra/Viper patterns for all phase projects
- **Naming**: Use descriptive variable names that indicate FTS5 context (e.g., `ftsDB`, `bm25Score`)

## Phase Development Workflow

### Planning Phase
- Load phase requirements and context artifacts into session
- Engage relevant subagents from `.claude/agents/` for specialized tasks
- Establish execution plan as small, reviewable tasks
- Write execution plan to `_context/plans/XX-phase-name.md`

### Pair Programming Execution
- Work through task list together in normal mode (not auto-accept)
- Single-step processes for easy review and approval
- Allow for clarifications and pivots throughout execution
- Use subagents for schema design, BM25 analysis, Go implementation, etc.

### Phase Completion & Validation
- Use `testing-validation-agent` to validate project functionality
- Generate comprehensive README.md capturing learning objectives and reflections
- Ensure phase works as standalone CLI program
- Document key learnings for context handoff to next phase

### Phase Template Structure
Each phase directory contains:
```
src/XX-phase-name/
├── main.go              # CLI entry point using Cobra/Viper
├── go.mod               # Isolated dependencies
├── README.md            # Learning objectives, concepts, usage, reflections
└── examples/            # Sample data files for experiments
```

## Experimental Project Guidelines

- **Single Responsibility**: Each phase focuses on specific FTS5/BM25 concepts
- **Standalone Operation**: Projects must run independently via CLI
- **Reproducible**: Include sample data generation in each experiment
- **Educational Focus**: Emphasize learning and understanding over production features

## BM25 Implementation Notes

- **Constants**: SQLite FTS5 uses k1=1.2, b=0.75 (hardcoded, non-configurable)
- **Scoring**: Negative scores where -1.5 ranks higher than -3.2
- **Column Weighting**: Use `bm25(table_name, weight1, weight2, ...)` for custom column importance
- **Performance**: BM25 requires columnsize backing table (enabled by default)

## Subagent Coordination

- **File Conflicts**: Subagents must coordinate file modifications to prevent overwrites
- **Responsibility Boundaries**: Each subagent owns specific file patterns or directories
- **Communication**: Subagents report findings to main session before making changes
- **Testing**: Subagents run tests before committing changes to validate functionality

## Common Patterns to Avoid

- **LIKE Queries**: Don't use LIKE with FTS5 tables, use MATCH operator
- **Missing Indexes**: Don't forget FTS5 automatically creates indexes
- **Score Interpretation**: Don't assume positive BM25 scores (SQLite inverts them)
- **Column Types**: Don't specify column types in FTS5 CREATE VIRTUAL TABLE statements

## Useful Commands & Snippets

- **Enable FTS5**: `PRAGMA compile_options;` to verify FTS5 is available
- **Debug Queries**: Use `EXPLAIN QUERY PLAN` to understand FTS5 query execution
- **Check Index**: `.schema` to see generated FTS5 backing tables
- **Performance**: Use `PRAGMA table_info(fts_table)` to inspect virtual table structure

## Project Evolution Strategy

- **Phase Gates**: Complete each syllabus phase before advancing
- **Knowledge Validation**: Build working examples that demonstrate understanding
- **Complexity Scaling**: Add features only after mastering prerequisites
- **Context Management**: Clear session context between major phase transitions
