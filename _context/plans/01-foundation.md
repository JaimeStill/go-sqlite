# Phase 01: Foundation Execution Plan

## Learning Objectives

- Master basic FTS5 virtual table creation in Go using SQLite driver
- Implement fundamental CRUD operations with FTS5 MATCH queries
- Establish reliable build configuration with FTS5 support
- Create standard error handling patterns for SQLite operations
- Build CLI application using Cobra/Viper for consistent user interface

## Technical Requirements

**Build Configuration:**

- Go 1.24 with CGO_ENABLED=1
- Build with `-tags "fts5"` for FTS5 support
- `github.com/mattn/go-sqlite3` driver dependency
- In-memory database (`:memory:`) for experiments

**Environment Setup:**

- Verify local SQLite installation has FTS5 support
- Confirm PRAGMA compile_options includes ENABLE_FTS5
- Test Go build process with FTS5 tags

## Task Breakdown

- [ ] **Task 1: Environment Verification**
  - Verify Go 1.24 installation
  - Check SQLite FTS5 availability via PRAGMA compile_options
  - Test basic Go SQLite connection and FTS5 module loading
  - **Acceptance Criteria**: Can connect to SQLite and verify FTS5 is available

- [ ] **Task 2: Project Structure Setup**
  - Create `src/01-foundation/` directory structure
  - Initialize `go.mod` with FTS5-enabled dependencies
  - Set up `main.go` with Cobra/Viper CLI framework
  - **Acceptance Criteria**: Project builds successfully with `go build -tags fts5`

- [ ] **Task 3: Basic FTS5 Table Creation**
  - Implement function to create FTS5 virtual table
  - Add schema for simple document storage (id, title, content)
  - Test table creation with basic SQL verification
  - **Acceptance Criteria**: Can programmatically create FTS5 virtual tables

- [ ] **Task 4: Document Insert Operations**
  - Implement document insertion with proper error handling
  - Support both single and batch insert operations
  - Add input validation and SQLite error interpretation
  - **Acceptance Criteria**: Can reliably insert documents into FTS5 tables

- [ ] **Task 5: Basic Search Implementation**
  - Create search function using MATCH operator (not LIKE)
  - Implement result retrieval with proper scanning
  - Add query parameter sanitization to prevent syntax errors
  - **Acceptance Criteria**: Can perform basic full-text searches and retrieve results

- [ ] **Task 6: Error Handling System**
  - Define custom error types for common FTS5 issues
  - Implement proper SQLite error interpretation and wrapping
  - Add graceful handling of FTS5 unavailability
  - **Acceptance Criteria**: Provides meaningful error messages for all failure modes

- [ ] **Task 7: CLI Interface Implementation**
  - Add Cobra commands for create-table, insert, and search operations
  - Implement proper flag handling and input validation
  - Add help text and usage examples
  - **Acceptance Criteria**: All operations accessible via CLI with clear documentation

- [ ] **Task 8: Build and Documentation**
  - Document build process with proper FTS5 tags
  - Create examples/ directory with sample data
  - Add troubleshooting guide for common issues
  - **Acceptance Criteria**: Build process is documented and repeatable

## Subagent Coordination

**`sqlite-schema-agent`**:

- Design FTS5 virtual table schema for document storage
- Provide SQL creation statements and best practices
- Review schema design for Phase 1 requirements

**`go-integration-agent`**:

- Set up Go project structure with proper module configuration
- Implement Cobra/Viper CLI patterns
- Handle SQLite driver integration and build configuration
- Create error handling patterns specific to Go SQLite operations

**`testing-validation-agent`**:

- Validate CLI functionality through manual execution
- Verify build process and documentation accuracy
- Ensure all acceptance criteria are met through hands-on testing

## Validation Strategy

Verify CLI execution works for all intended scenarios including table creation, document insertion, search operations, and error handling.

## Success Metrics

**Functional Requirements:**

- [ ] FTS5 virtual tables can be created via CLI
- [ ] Documents can be inserted and searched through CLI commands
- [ ] Basic MATCH queries return expected results
- [ ] All SQLite errors produce meaningful CLI error messages

**Technical Requirements:**

- [ ] Project builds successfully with `go build -tags fts5`
- [ ] CLI interface is intuitive and well-documented
- [ ] Code follows Go best practices and project conventions
- [ ] All operations work with in-memory database

**Documentation Requirements:**

- [ ] README.md explains learning objectives and CLI usage
- [ ] Build process is clearly documented with troubleshooting
- [ ] CLI help text provides sufficient guidance for all commands
- [ ] Sample data and usage examples are provided

## Risk Mitigation

**Technical Risks:**

- FTS5 build issues: Validate environment setup early, provide troubleshooting guide
- SQLite version conflicts: Pin driver version, document compatibility requirements
- Cross-platform builds: Test on available platforms, document platform-specific requirements

**Learning Risks:**

- Scope creep: Focus strictly on Phase 1 deliverables, defer advanced features
- Complexity drift: Keep examples minimal and focused on fundamentals
- Context overload: Use subagents for specialized tasks to preserve main context

## Phase Completion Checklist

- [ ] All tasks completed with acceptance criteria met
- [ ] CLI functionality validated through manual execution
- [ ] Build process documented and verified
- [ ] README.md written with learning reflections and usage instructions
- [ ] Code reviewed for quality and consistency
- [ ] Sample data and examples provided
- [ ] Phase ready for knowledge handoff to Phase 2
