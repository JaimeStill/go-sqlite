# Phase Development Session Initialization

## Context Loading

Load the following context files to understand the project scope and requirements:

1. **Project Foundation**:
   - Read `/home/jaime/research/go-sqlite/CLAUDE.md` - Core project directives and development philosophy
   - Read `/home/jaime/research/go-sqlite/ROADMAP.md` - Phase objectives and deliverables

2. **Technical References**:
   - Read `/home/jaime/research/go-sqlite/_context/sqlite-fts5-reference.md` - SQLite FTS5 implementation details
   - Read `/home/jaime/research/go-sqlite/_context/bm25-reference.md` - BM25 scoring mechanics
   - Read `/home/jaime/research/go-sqlite/_context/go-sqlite-reference.md` - Go driver configuration

3. **Phase-Specific Context**:
   - Read current phase directory (`src/XX-phase-name/`) to understand existing progress
   - Check for existing execution plan in `_context/plans/XX-phase-name.md`

## Planning Phase Execution

### 1. Requirements Analysis

- Extract specific deliverables and success criteria from ROADMAP.md for current phase
- Identify key learning outcomes that must be demonstrated
- Note any dependencies on previous phases or shared utilities

### 2. Subagent Engagement Strategy

Based on the phase requirements, proactively engage relevant subagents:

- **`sqlite-schema-agent`**: For FTS5 virtual table design and schema optimization
- **`bm25-research-agent`**: For BM25 scoring analysis, ranking interpretation, and performance studies
- **`go-integration-agent`**: For Go implementation patterns, error handling, and CLI design
- **`testing-validation-agent`**: For test design, validation, and functionality verification
- **`performance-analysis-agent`**: For profiling, benchmarking, and optimization analysis

### 3. Execution Plan Development

Create a detailed task breakdown following these principles:

**Task Granularity**:

- Each task should be completable in a single focused session
- Tasks must have clear acceptance criteria
- Dependencies between tasks must be explicit

**Phase Template Adherence**:

- Ensure `main.go` uses Cobra/Viper CLI patterns
- Plan for isolated `go.mod` with proper FTS5 build tags
- Include sample data generation for experiments
- Structure for comprehensive `README.md` at completion

**Code Quality Standards**:

- Explicit SQLite error handling in every function
- Resource cleanup with defer patterns
- Descriptive variable names indicating FTS5 context
- Build verification with `go build -tags "fts5"`

### 4. Documentation Requirements

Write execution plan to `_context/plans/XX-phase-name.md` including:

```markdown
# Phase XX: [Phase Name] Execution Plan

## Learning Objectives
[Specific skills and understanding to be demonstrated]

## Technical Requirements
[Build configuration, dependencies, environment setup]

## Task Breakdown
- [ ] Task 1: [Description with acceptance criteria]
- [ ] Task 2: [Description with acceptance criteria]
[Continue with all planned tasks]

## Subagent Coordination
[Which agents will be used for which tasks]

## Validation Strategy
[How functionality will be tested and verified]

## Success Metrics
[Measurable outcomes that indicate phase completion]
```

## Pair Programming Execution Mode

### Development Standards

- Work in normal mode (not auto-accept) for review opportunities
- Single-step progression through task list
- Allow for clarifications and course corrections
- Use TodoWrite tool to track progress throughout execution

### Code Implementation Patterns

- Always handle SQLite errors explicitly with meaningful messages
- Use in-memory databases (`:memory:`) for experiments unless persistence needed
- Follow FTS5 best practices: MATCH operators, proper virtual table syntax
- Implement CLI with Cobra/Viper for consistency across phases

### Integration Points

- Leverage `src/shared/` utilities when appropriate
- Ensure phase can operate standalone via CLI
- Maintain data compatibility with shared `data/` directory
- Document any new shared utilities created

## Phase Completion Protocol

### Functionality Validation

- Engage `testing-validation-agent` to verify all deliverables
- Run comprehensive tests including edge cases
- Validate CLI functionality and error handling
- Confirm FTS5 build process works correctly

### Documentation Generation

- Create comprehensive `README.md` capturing:
  - Learning objectives and how they were achieved
  - Key concepts demonstrated with code examples  
  - Usage instructions for CLI program
  - Reflections on challenges and insights gained
  - Context for next phase development

### Knowledge Transfer Preparation

- Document key learnings for future phase context
- Identify reusable patterns and utilities
- Note any architectural decisions that impact future phases
- Ensure project state is clean for next phase initialization

## Session Management

### Context Optimization

- Clear unnecessary context between major task transitions
- Focus on current phase scope to avoid context drift
- Use subagents for specialized tasks to preserve main context
- Regular progress checkpoints with TodoWrite tool
- **Complex phases may require multiple sessions**: Plan for context handoff between sessions
- **Architecture decisions early**: Establish CLI patterns and code organization upfront to avoid major refactoring

### Effective Subagent Utilization

- **Design consultation first**: Engage relevant subagents to analyze requirements and propose approaches before implementation
- **Domain-specific expertise**: Use subagents for their specialized knowledge (schema design, BM25 analysis, Go patterns, testing strategies)
- **Output-driven execution**: Take subagent recommendations and implement them in main session to maintain control
- **Iterative feedback loop**: Use subagent analysis to validate approaches and catch issues early
- **Context preservation**: Delegate research and analysis to subagents while keeping implementation decisions in main context

### Quality Assurance

- Validate understanding through working code demonstrations
- Test all CLI functionality before marking tasks complete
- Ensure reproducible builds with documented dependencies
- Maintain educational focus over feature complexity
- **Iterative refinement expected**: Allow for multiple rounds of user feedback and architectural improvements

---

**Initialization Checklist**:

- [ ] All context files loaded and understood
- [ ] Phase requirements extracted and analyzed  
- [ ] Relevant subagents identified for engagement
- [ ] Execution plan written to `_context/plans/`
- [ ] TodoWrite tool initialized with task breakdown
- [ ] Development environment verified (Go 1.24, FTS5 support)
- [ ] Ready to begin pair programming execution
