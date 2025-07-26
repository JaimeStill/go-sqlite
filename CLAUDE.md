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
- **Testing Execution**: Use `go run -tags fts5` instead of `go build` to avoid creating binaries that need .gitignore management
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

## CLI Command Standards

- **Flag-based Arguments**: Use flags (`--size 100`) instead of positional arguments for all parameters
- **Error Handling**: Use `RunE` instead of `Run` to enable proper error returns without wrapper functions
- **Flag Access**: Use `cmd.Flags().GetString()` pattern directly in RunE functions for clean, readable code
- **Flag Utilities**: Create helper functions for common flag patterns to reduce boilerplate

## Naming Conventions

- **Config vs Options**: Reserve `Config` suffix for application configuration structures in the `config` package. Use `Options` suffix for operation-specific parameter structures (e.g., `CorpusOptions`, `SearchOptions`)
- **Default Functions**: Name default constructors as `DefaultXxxOptions()` for option structures
- **Variable Names**: Use descriptive names like `options` for option structures, `config.App` for application configuration

## Successful Architectural Patterns

### CommandGroup Hierarchical Structure

Use the CommandGroup pattern for scalable CLI organization:

```go
type CommandGroup struct {
    Command     *cobra.Command      // The Cobra command this group represents
    ChildGroups []*CommandGroup     // Child command groups (base commands)
    SubCommands []*cobra.Command    // Direct sub-commands
    FlagSetup   func()             // Flag registration function
}
```

**Benefits**:

- Prevents command naming conflicts
- Enables recursive command organization
- Simplifies command registration with `Init()` method

### Type-Safe Error Handling System

Implement centralized error handling with sentinel errors:

```go
// Define error types
var (
    ErrValidation = errors.New("validation failed")
    ErrDatabase   = errors.New("database operation failed")
    ErrFTS5       = errors.New("FTS5 operation failed")
    ErrNotFound   = errors.New("not found")
)

// Helper functions for creating typed errors
func Validationf(format string, args ...interface{}) error {
    return fmt.Errorf("%w: "+format, append([]interface{}{ErrValidation}, args...)...)
}

// Centralized display with automatic verbose handling
func DisplayError(err error) {
    if viper.GetBool("verbose") {
        displayVerbose(err)
    } else {
        displaySimple(err)
    }
}
```

**Benefits**:

- Consistent error presentation across commands
- Type-safe error checking with `errors.Is()`
- Automatic verbose/simple mode switching
- User-friendly error messages with helpful hints

### Structured Configuration Management

Use a structured configuration approach with validation:

```go
// Config represents the application configuration schema
type Config struct {
    // Global settings
    Database string `mapstructure:"database"`
    Verbose  bool   `mapstructure:"verbose"`
    Format   string `mapstructure:"format"`
    
    // Nested configuration sections
    Corpus        CorpusConfig        `mapstructure:"corpus"`
    Search        SearchConfig        `mapstructure:"search"`
    Visualization VisualizationConfig `mapstructure:"visualization"`
}

// Configuration methods
func NewConfig() *Config                    // Create with defaults
func (c *Config) SetDefaults()             // Apply defaults to viper
func (c *Config) Validate() error          // Validate settings
func (c *Config) Load() error              // Load from viper and validate
```

**Implementation Pattern**:

1. Define configuration as nested structs with `mapstructure` tags
2. Provide `NewConfig()` factory with sensible defaults
3. Implement `Validate()` for runtime validation
4. Use `Load()` to unmarshal from viper and validate
5. Expose global `App` instance following the same pattern as `commands.Root`
6. **Initialize config AFTER flag parsing using `PersistentPreRun` hook**

**Initialization Flow Pattern**:

```go
// In commands/root.go
var rootCmd = &cobra.Command{
    // ... other fields
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        // Initialize configuration after flags are parsed
        config.App.Init()
    },
}

// In main.go
func main() {
    // Initialize command structure first (sets up flags)
    commands.Root.Init()
    
    // Execute commands (triggers PersistentPreRun -> config.App.Init())
    err := commands.Root.Command.Execute()
    // ...
}

// In other packages
import "path/to/config"
db, err := handlers.NewDatabase(config.App.GetDatabasePath())
```

**Key Insight**: Configuration must be initialized AFTER cobra parses flags, not before. Use `PersistentPreRun` to ensure config initialization happens after flag parsing but before command execution.

**Benefits**:

- Type-safe configuration access throughout the application
- Centralized validation logic
- Clear configuration schema documentation
- Easy testing with struct instantiation
- Supports config files, env vars, and flags
- Consistent with other global patterns like `commands.Root`

### Global Instance Pattern with Centralized Initialization

Establish global instances for core application components that are initialized once and accessed throughout the application:

```go
// config/config.go
type Config struct {
    DatabasePath string
    Verbose      bool
    Format       string
    // Nested configuration sections
    Corpus CorpusConfig `mapstructure:"corpus"`
    Search SearchConfig `mapstructure:"search"`
}

// App is the global configuration instance
var App *Config

func (c *Config) Init() error {
    // Initialize from viper after flags are parsed
    App = &Config{
        DatabasePath: viper.GetString("database"),
        Verbose:      viper.GetBool("verbose"),
        Format:       viper.GetString("format"),
    }
    // Apply defaults and validate
    return App.validate()
}

// database/database.go
type Database struct {
    db   *sql.DB
    path string
}

// Instance is the global database instance
var Instance *Database

func Init(dataSourceName string) error {
    db, err := NewDatabase(dataSourceName)
    if err != nil {
        return fmt.Errorf("initializing database: %w", err)
    }
    Instance = db
    return nil
}
```

**Benefits**:

- Single point of initialization in root command's PersistentPreRun
- Eliminates dependency injection complexity
- Global accessibility without circular dependencies
- Clear initialization order and error handling

### Factory Function Command Pattern

Use private factory functions to create command instances while exposing public variables for registration:

```go
// commands/corpus.go
// Corpus is the public corpus command group instance
var Corpus = newCorpusGroup()

// newCorpusGroup creates the corpus command group with all its subcommands
func newCorpusGroup() *CommandGroup {
    // All command variables scoped within this function
    corpusCmd := &cobra.Command{
        Use:   "corpus",
        Short: "Manage document corpus for BM25 experiments",
    }
    
    generateCmd := &cobra.Command{
        Use:   "generate",
        Short: "Generate synthetic documents",
        RunE:  handlers.Corpus.HandleGenerate,
    }
    
    statsCmd := &cobra.Command{
        Use:   "stats", 
        Short: "Display corpus statistics",
        RunE:  handlers.Corpus.HandleStats,
    }
    
    // Setup flags
    setupFlags := func() {
        generateCmd.Flags().IntP("size", "s", 0, "number of documents")
        // ... other flags
    }
    
    return &CommandGroup{
        Command:     corpusCmd,
        SubCommands: []*cobra.Command{generateCmd, statsCmd},
        FlagSetup:   setupFlags,
    }
}
```

**Benefits**:

- Prevents command variable naming conflicts between files
- Encapsulates command creation logic in private functions
- Public variables enable clean registration in root command
- Command variables are scoped within factory functions

### Stateless Handler Pattern with Global Access

Implement handlers as empty structs that access global instances directly:

```go
// handlers/corpus.go
// Corpus is the global corpus handler instance
var Corpus CorpusHandler

// CorpusHandler manages corpus operations (stateless - accesses global instances)
type CorpusHandler struct{}

func (h *CorpusHandler) HandleGenerate(cmd *cobra.Command, args []string) error {
    // Extract flags
    size, _ := cmd.Flags().GetInt("size")
    
    // Access global config instance directly
    if size == 0 {
        size = config.App.Corpus.Size
    }
    
    // Access global database instance directly
    if err := database.Instance.InitSchema(ctx); err != nil {
        return err
    }
    
    // Use config for verbose output
    if config.App.Verbose {
        fmt.Printf("Generating %d documents...\n", size)
    }
    
    return h.generateCorpus(ctx, size)
}

func (h *CorpusHandler) HandleStats(cmd *cobra.Command, args []string) error {
    // Direct access to global instances
    stats, err := h.getCorpusStats(ctx)
    if err != nil {
        return err
    }
    
    // Use global config for output format
    return h.displayStats(stats, config.App.Format)
}
```

**Benefits**:

- No initialization or dependency management needed
- Direct access to global state without parameter passing
- Simplified handler registration and testing
- Clear separation between command handling and business logic

### Centralized Initialization Flow

Establish a clear initialization sequence in the root command's PersistentPreRun:

```go
// commands/root.go
var rootCmd = &cobra.Command{
    Use:   "bm25-fundamentals",
    Short: "SQLite FTS5 BM25 fundamentals exploration",
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        // 1. Initialize configuration after flags are parsed
        config.App.Init()
        
        // 2. Initialize database connection using config
        if err := database.Init(config.App.GetDatabasePath()); err != nil {
            fmt.Fprintf(os.Stderr, "Database initialization error: %v\n", err)
            os.Exit(1)
        }
        
        // 3. Handlers are stateless - no initialization needed
    },
}

// Root represents the root command group
var Root = &CommandGroup{
    Command: rootCmd,
    ChildGroups: []*CommandGroup{
        Corpus,  // Public instance from corpus.go
        Search,  // Public instance from search.go
    },
    FlagSetup: setupGlobalFlags,
}
```

**Benefits**:

- Single initialization point for entire application
- Clear dependency order: config → database → handlers
- Centralized error handling for startup failures
- Consistent initialization across all commands

### Database Package Centralization

Centralize all database operations in a dedicated package with global access:

```go
// database/database.go
type Database struct {
    db   *sql.DB
    path string
}

// Instance is the global database instance
var Instance *Database

func Init(dataSourceName string) error {
    Instance = &Database{path: dataSourceName}
    return Instance.connect()
}

func (d *Database) InitSchema(ctx context.Context) error {
    // Create FTS5 virtual table
    query := `
        CREATE VIRTUAL TABLE IF NOT EXISTS documents_fts USING fts5(
            title, content, category
        );
        
        CREATE TABLE IF NOT EXISTS documents (
            id INTEGER PRIMARY KEY,
            title TEXT NOT NULL,
            content TEXT NOT NULL,
            category TEXT NOT NULL,
            length INTEGER NOT NULL,
            created TEXT NOT NULL
        );
    `
    _, err := d.db.ExecContext(ctx, query)
    return err
}

func (d *Database) SearchBM25(ctx context.Context, query string, options SearchOptions) ([]*SearchResult, error) {
    // Build FTS5 query with BM25 scoring
    sqlQuery := `
        SELECT d.id, d.title, d.content, d.category, d.length, d.created,
               bm25(documents_fts) as score
        FROM documents d
        JOIN documents_fts fts ON d.id = fts.rowid
        WHERE documents_fts MATCH ?
        ORDER BY rank
        LIMIT ?
    `
    // Implementation...
}
```

**Benefits**:

- Single source of truth for database operations
- Global accessibility without circular dependencies
- Encapsulation of SQLite/FTS5 specifics
- Consistent error handling across database operations

### Layered Architecture Pattern

Maintain clear separation of concerns:

- **Commands Layer**: CLI interface, flag handling, error display
- **Handlers Layer**: Business logic, database operations, validation
- **Models Layer**: Data structures, type definitions
- **Errors Layer**: Error types, display functions, type checking

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

### Guide Generation

- Generate instructional GUIDE.md at phase completion using `_prompts/generate-guide.md`
- Execute via: `Execute _prompts/generate-guide.md for src/##-[phase-name]/`
- Include thorough description of all learning objectives with analogical illustrations
- Explain core project infrastructure directly associated with learning objectives
- Provide walkthroughs of command execution and outputs showing how commands function and relate to learning objectives
- Focus on educational value and instructional clarity for future learners

### Phase Template Structure

Each phase directory contains:

```
src/XX-phase-name/
├── README.md                    # Learning objectives, concepts, usage, reflections  
├── go.mod                       # Isolated dependencies for the phase
├── go.sum                       # Dependency lock file
└── [phase-name]/               # Root package directory (e.g., fts5-foundation, bm25-fundamentals)
    ├── main.go                 # CLI entry point using Cobra/Viper
    ├── config/                 # Configuration management package
    │   └── config.go          # Viper integration with global App instance
    ├── database/              # Database operations package
    │   └── database.go        # SQLite/FTS5 operations with global Instance
    ├── commands/              # Hierarchical CLI command definitions
    │   ├── command_group.go   # CommandGroup pattern for hierarchical structure
    │   ├── root.go           # Root command and centralized initialization
    │   └── [context].go      # Context-specific commands using factory pattern
    ├── handlers/             # Business logic layer (stateless)
    │   └── [context].go      # Stateless handlers accessing global instances
    ├── models/               # Data structures and types
    │   └── [context].go      # Context-specific models
    └── errors/               # Type-safe error handling system
        └── errors.go         # Sentinel errors and display functions
```

**Important**: The phase directory should contain README.md, go.mod, and go.sum at its root, with all Go code organized within a subdirectory named after the phase (using the same naming convention). This keeps the phase root clean and allows for better module organization.

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
- **Complex Dependency Injection**: Avoid complex DI patterns when global instances suffice for CLI applications
- **Stateful Handlers**: Keep handlers stateless to simplify initialization and testing
- **Early Config Initialization**: Don't initialize config before flag parsing; use PersistentPreRun
- **Command Naming Conflicts**: Avoid naming command variables the same across different command files
- **Circular Dependencies**: Structure packages to avoid import cycles (config ← database ← handlers)

## Useful Commands & Snippets

- **Enable FTS5**: `PRAGMA compile_options;` to verify FTS5 is available
- **Debug Queries**: Use `EXPLAIN QUERY PLAN` to understand FTS5 query execution
- **Check Index**: `.schema` to see generated FTS5 backing tables
- **Performance**: Use `PRAGMA table_info(fts_table)` to inspect virtual table structure

## Session Management Best Practices

### Context Optimization for Complex Phases

- **Chunked Development**: Break complex phases into 2-3 focused sessions
- **Progressive Complexity**: Start with basic structure, add sophistication iteratively  
- **Checkpoint Documentation**: Document architecture decisions at major transition points
- **Subagent Utilization**: Use specialized agents (go-integration-agent) for design reviews

### Effective Pair Programming Patterns

- **Architecture First**: Establish patterns early, then apply consistently
- **Incremental Validation**: Test each component as built rather than at end
- **User Feedback Integration**: Regular check-ins for course corrections
- **Pattern Documentation**: Capture successful approaches for reuse

### Context Handoff Between Sessions

- **Status Documentation**: Clear TODO lists with detailed context
- **Pattern Summary**: Document architectural decisions and successful patterns
- **File Organization**: Keep related changes in focused commits
- **Validation State**: Ensure working state before session transitions

## Project Evolution Strategy

- **Phase Gates**: Complete each syllabus phase before advancing
- **Knowledge Validation**: Build working examples that demonstrate understanding
- **Complexity Scaling**: Add features only after mastering prerequisites
- **Architecture Reuse**: Apply successful patterns from previous phases
