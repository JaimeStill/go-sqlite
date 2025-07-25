# Claude Code Subagent Profiles

## Project Architecture

### Subagent Overview

Custom sub agents in Claude Code are specialized AI assistants that can be invoked to handle specific types of tasks. They enable more efficient problem-solving by providing task-specific configurations with customized system prompts, tools and a separate context window.

### Key Benefits for This Project

- **Context preservation**: Each subagent operates independently, preventing main conversation pollution
- **Specialized expertise**: Domain-specific instructions for SQLite FTS5 and BM25 development
- **Reusability**: Share subagents across learning phases and team members
- **Flexible permissions**: Limit powerful tools to appropriate subagent types

---

## Usage Instructions

### Creating Project Subagents

```bash
# Navigate to project root
cd your-fts5-project

# Create agents directory
mkdir -p .claude/agents

# Copy agent files to project
cp sqlite-schema-agent.md .claude/agents/
cp bm25-research-agent.md .claude/agents/
cp go-integration-agent.md .claude/agents/
cp testing-validation-agent.md .claude/agents/
cp performance-analysis-agent.md .claude/agents/
```

### Managing Subagents

```bash
# View all available subagents
/agents

# Create new subagent interactively
/agents

# Edit existing subagent
/agents
```

### Invocation Patterns

#### Automatic Delegation

Claude Code proactively delegates tasks based on the task description in your request and the description field in sub agent configurations.

Example requests that trigger automatic delegation:

```
"Design an FTS5 virtual table for document search" 
→ Automatically invokes sqlite-schema-agent

"Analyze BM25 scoring distribution for these search results"
→ Automatically invokes bm25-research-agent

"Implement Go database connection pooling for SQLite"
→ Automatically invokes go-integration-agent
```

#### Explicit Invocation

```
"@sqlite-schema-agent create external content table for articles"
"@bm25-research-agent evaluate column weighting strategies"
"@testing-validation-agent design benchmark suite for search performance"
```

### Multi-Agent Workflows

#### Schema Design Session

```
1. "Design FTS5 schema for contextual memory system"
   → sqlite-schema-agent proposes initial design
   
2. "@performance-analysis-agent review proposed schema for bottlenecks"
   → performance-analysis-agent analyzes design
   
3. "@go-integration-agent implement connection management for this schema"
   → go-integration-agent creates Go implementation
```

#### Performance Optimization Sprint

```
1. "@performance-analysis-agent profile current search performance"
   → Identifies bottlenecks through profiling
   
2. "@bm25-research-agent analyze scoring algorithm efficiency"
   → Analyzes ranking algorithm performance
   
3. "@sqlite-schema-agent propose index optimizations"
   → Suggests schema improvements
```

### Best Practices

#### Subagent Design

- **Single responsibility**: Each agent focuses on one domain area
- **Detailed descriptions**: Include specific trigger phrases and use cases
- **Tool limitation**: Grant only necessary tools for security and focus
- **Action-oriented descriptions**: Use phrases like "MUST BE USED" and "PROACTIVELY"

#### Project Integration

- **Version control**: Check `.claude/agents/` into git for team collaboration
- **Documentation**: Include subagent usage in project README
- **Iteration**: Start with generated agents, then customize for project needs
- **Testing**: Validate subagent behavior with test queries

### Performance Considerations

- **Context efficiency**: Subagents preserve main context for longer sessions
- **Startup latency**: Each invocation starts with clean context
- **Tool access**: Limited tool access improves focus and security
- **Parallel execution**: Up to 10 concurrent subagents supported
