# Phase 1: SQLite FTS5 Foundation Learning Guide

## Table of Contents

1. [Learning Objectives Overview](#learning-objectives-overview)
2. [Conceptual Foundations](#conceptual-foundations)
3. [Project Infrastructure](#project-infrastructure)
4. [Interactive Learning](#interactive-learning)
5. [Hands-On Exercises](#hands-on-exercises)
6. [Self-Assessment](#self-assessment)
7. [Troubleshooting Guide](#troubleshooting-guide)
8. [Next Steps](#next-steps)

## Introduction

Welcome to the foundational phase of your SQLite FTS5 and BM25 learning journey! This phase establishes the core concepts and practical skills needed to work with SQLite's Full-Text Search (FTS5) capabilities. Through hands-on experimentation with a purpose-built CLI tool, you'll develop a solid understanding of how modern search engines index and retrieve documents.

Think of this phase as learning to use a powerful telescope before exploring the cosmos. FTS5 is your telescope for searching through vast amounts of text data, and BM25 is the sophisticated optical system that determines which documents are most relevant to your queries.

## Learning Objectives Overview

### Core Learning Objectives

By completing this phase, you will master five essential competencies that form the foundation for all advanced FTS5 work:

#### 1. **FTS5 Virtual Tables**: Create and Configure Search Infrastructure

**What You'll Learn**: How to create SQLite virtual tables specifically designed for full-text search, understanding the fundamental difference between regular SQL tables and FTS5's specialized indexing structures.

**Why It Matters**: Virtual tables are the gateway to FTS5's powerful search capabilities. Unlike traditional tables that store data row-by-row, FTS5 virtual tables automatically create inverted indexes that enable lightning-fast text searches across thousands or millions of documents.

**Practical Skills**:
- Create FTS5 virtual tables with the proper syntax
- Configure tokenizers (unicode61 with diacritics removal)
- Understand the difference between data storage and search indexes
- Design table schemas optimized for different search patterns

#### 2. **BM25 Scoring**: Understand Modern Relevance Ranking

**What You'll Learn**: How SQLite implements the BM25 relevance scoring algorithm, including its unique negative score approach and the mathematical principles that determine document ranking.

**Why It Matters**: BM25 is the gold standard for text relevance scoring used by major search engines. Understanding how it evaluates term frequency, document length, and collection statistics is crucial for building effective search systems.

**Practical Skills**:
- Interpret BM25 scores correctly (negative values, lower = better)
- Understand how term frequency and document length affect scoring
- Use BM25 scores to rank search results meaningfully
- Recognize SQLite's specific BM25 implementation characteristics

#### 3. **Search Patterns**: Master FTS5 Query Techniques

**What You'll Learn**: The various ways to query FTS5 tables, from simple text searches to sophisticated field-specific and category-filtered queries using the MATCH operator.

**Why It Matters**: Different search patterns serve different use cases. Basic searches find general matches, field-specific searches target titles or specific content areas, and category filters narrow results to specific document types.

**Practical Skills**:
- Construct basic MATCH queries for general text search
- Use field-specific syntax (e.g., `title:programming`)
- Combine categories with text searches for precise filtering
- Apply search operators for complex query logic

#### 4. **Index Management**: Understand Automatic FTS5 Indexing

**What You'll Learn**: How FTS5 automatically maintains sophisticated indexes during Create, Read, Update, and Delete (CRUD) operations, ensuring search performance remains optimal without manual intervention.

**Why It Matters**: Unlike traditional databases where you manually create and maintain indexes, FTS5 handles this complexity transparently. Understanding this automation helps you design efficient applications and troubleshoot performance issues.

**Practical Skills**:
- Perform CRUD operations on FTS5 tables
- Understand how inserts/updates affect the search index
- Recognize when index rebuilding occurs
- Design workflows that leverage automatic index maintenance

#### 5. **Error Handling**: Build Robust Search Applications

**What You'll Learn**: Type-safe error handling patterns specific to FTS5 applications, including validation errors, database connection issues, and FTS5-specific configuration problems.

**Why It Matters**: Search systems must gracefully handle various failure modes. Proper error handling ensures your applications provide meaningful feedback to users and administrators when issues occur.

**Practical Skills**:
- Implement sentinel error patterns for different error categories
- Provide user-friendly error messages with helpful guidance
- Handle FTS5-specific errors (like missing build tags)
- Design robust error recovery strategies

### Contextual Placement in Learning Project

This foundation phase is the first of six phases in a comprehensive SQLite FTS5 learning journey:

```
Phase 1: Foundation (YOU ARE HERE)
    ↓ Establishes basic FTS5 competency
Phase 2: BM25 Fundamentals
    ↓ Deepens understanding of relevance scoring
Phase 3: Query Operations
    ↓ Advanced search patterns and techniques
Phase 4: Ranking & Relevance
    ↓ Custom ranking strategies
Phase 5: Advanced Features
    ↓ FTS5 auxiliary functions and specialized tools
Phase 6: Integration Patterns
    ↓ Production-ready implementation patterns
```

**Prerequisites**: Basic Go programming knowledge, SQL familiarity, and command-line comfort. No prior FTS5 experience required.

**What Comes Next**: Phase 2 builds on this foundation by exploring advanced BM25 concepts, corpus analysis, and relevance tuning techniques.

## Conceptual Foundations

### Understanding FTS5 Virtual Tables: The Library Analogy

Imagine you're organizing a massive library with millions of books. A traditional approach would be like a regular SQL table - you store each book in a specific location and create a simple catalog with basic information (title, author, shelf number). To find books about "machine learning," you'd need to manually scan through every catalog entry.

FTS5 virtual tables work like having a team of expert librarians who create specialized finding aids:

**Traditional Table (Regular Library Catalog)**:
```sql
CREATE TABLE books (
    id INTEGER PRIMARY KEY,
    title TEXT,
    content TEXT,
    category TEXT
);
-- Finding books requires scanning every row
SELECT * FROM books WHERE content LIKE '%machine learning%';
```

**FTS5 Virtual Table (Expert Librarian System)**:
```sql
CREATE VIRTUAL TABLE books_fts USING fts5(
    title,
    content, 
    category,
    tokenize='unicode61 remove_diacritics 1'
);
-- Finding books uses sophisticated indexes
SELECT * FROM books_fts WHERE books_fts MATCH 'machine learning';
```

The FTS5 "librarians" automatically:
- **Tokenize**: Break down every book's content into individual words
- **Index**: Create reverse lookups (word → list of books containing it)
- **Normalize**: Handle variations (Machine = machine, café = cafe)
- **Score**: Rate how relevant each book is to your query

### BM25 Scoring: The Expertise Rating System

BM25 is like having an expert reviewer who rates how well each book matches your search based on three key factors:

#### 1. **Term Frequency (TF)**: How often does your search term appear?

Think of a book about "python programming" vs a book that mentions "python" once in passing. The book about python programming will have a higher relevance score because "python" appears frequently throughout the text.

```
Book A: "python" appears 50 times → Higher relevance
Book B: "python" appears 2 times → Lower relevance  
```

#### 2. **Document Length**: How long is the document relative to others?

A 10-page article that mentions "python" 5 times is more focused than a 500-page book that mentions "python" 5 times. BM25 adjusts scores based on document length to favor focused content.

```
Short Article: 5 mentions in 10 pages → Higher density, better score
Long Book: 5 mentions in 500 pages → Lower density, reduced score
```

#### 3. **Inverse Document Frequency (IDF)**: How rare is the term?

If "the" appears in 90% of documents, finding it doesn't tell us much. If "quantum computing" appears in only 2% of documents, finding it is highly significant.

```
Common word "the": Low IDF → Less impact on relevance
Rare phrase "quantum computing": High IDF → Major impact on relevance
```

#### SQLite's Unique BM25 Implementation

SQLite uses **negative scores** where **lower numbers indicate better matches**:

```
Score: -0.5  → Excellent match (better)
Score: -2.3  → Good match 
Score: -5.1  → Poor match (worse)
```

This is the opposite of many search engines, so remember: **lower = better** in SQLite FTS5.

### The MATCH Operator: Your Search Language

The MATCH operator is your interface to FTS5's power. Unlike SQL's LIKE operator which does simple pattern matching, MATCH leverages the full FTS5 index system:

#### Basic Search
```sql
-- LIKE: Slow, no relevance scoring
SELECT * FROM books WHERE content LIKE '%python programming%';

-- MATCH: Fast, relevance-scored results
SELECT * FROM books_fts WHERE books_fts MATCH 'python programming';
```

#### Field-Specific Search
```sql
-- Search only in titles
SELECT * FROM books_fts WHERE books_fts MATCH 'title:introduction';

-- Search only in content
SELECT * FROM books_fts WHERE books_fts MATCH 'content:advanced techniques';
```

#### Category Filtering
```sql
-- Combine category filter with text search
SELECT * FROM books_fts WHERE books_fts MATCH 'category:programming AND python';
```

### Common Misconceptions and Clarifications

#### Misconception 1: "FTS5 tables store data differently"
**Reality**: FTS5 virtual tables store the same data as regular tables, but they automatically maintain additional index structures for fast searching.

#### Misconception 2: "Higher BM25 scores are always better"
**Reality**: In SQLite FTS5, lower (more negative) scores indicate better matches. This is opposite to many other search systems.

#### Misconception 3: "LIKE and MATCH do the same thing"
**Reality**: LIKE does pattern matching without understanding word boundaries or relevance. MATCH uses linguistic analysis and relevance scoring.

#### Misconception 4: "FTS5 requires manual index maintenance"
**Reality**: FTS5 automatically maintains all indexes during insert, update, and delete operations. No manual maintenance required.

## Project Infrastructure

The `fts5-foundation` project demonstrates professional Go application architecture while keeping the focus on FTS5 learning objectives. Understanding this structure helps you see how FTS5 concepts apply in real applications.

### Architectural Overview

The project uses a **layered architecture** that separates concerns and makes the FTS5 concepts clear:

```
📁 fts5-foundation/
├── main.go                    # Application entry point
├── commands/                  # CLI interface layer
│   ├── command_group.go      # Hierarchical command pattern
│   ├── root.go              # Global configuration and initialization
│   └── document.go          # Document-specific commands
├── config/                   # Configuration management
│   └── config.go           # Global config with Viper integration
├── database/                 # FTS5 operations layer
│   └── database.go         # Global database instance and FTS5 operations
├── handlers/                # Business logic layer
│   └── document.go        # Stateless handlers for document operations
├── models/                  # Data structures
│   └── document.go       # Document, SearchResult, and DocumentInfo models
└── errors/                  # Error handling system
    └── errors.go         # Type-safe errors and user-friendly display
```

### Key Architectural Patterns and Their Educational Benefits

#### 1. **CommandGroup Pattern**: Hierarchical CLI Structure

**Why This Pattern**: Real-world search applications often have complex command structures. This pattern demonstrates how to organize FTS5 operations logically.

**Educational Value**: Shows how different FTS5 operations (create, search, update) can be grouped and presented to users in an intuitive way.

```go
// From commands/command_group.go
type CommandGroup struct {
    Command     *cobra.Command      // The Cobra command this group represents
    ChildGroups []*CommandGroup     // Child command groups
    SubCommands []*cobra.Command    // Direct sub-commands
    FlagSetup   func()             // Flag registration function
}
```

The resulting CLI structure makes FTS5 concepts discoverable:
```
fts5-foundation
└── document
    ├── create-table      # FTS5 virtual table creation
    ├── insert           # Single document insertion
    ├── batch-insert     # Multiple document insertion
    ├── search           # Basic MATCH queries
    ├── search-category  # Category-filtered searches
    ├── search-field     # Field-specific searches
    ├── list            # Document listing
    ├── update          # Document updates (reindexing)
    └── delete          # Document deletion (index cleanup)
```

#### 2. **Global Instance Pattern**: Simplified State Management

**Why This Pattern**: For CLI applications focused on learning, complex dependency injection adds cognitive overhead without educational value.

**Educational Value**: Allows learners to focus on FTS5 concepts rather than Go architecture complexity.

```go
// From database/database.go
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

This pattern enables clear, educational code in handlers:

```go
// From handlers/document.go
func (h *DocumentHandler) HandleSearch(cmd *cobra.Command, args []string) error {
    query := args[0]
    
    // Direct access to global database instance
    results, err := database.Instance.SearchDocuments(ctx, query, options)
    if err != nil {
        return err
    }
    
    // Display results with BM25 scores
    return h.displayResults(results, options)
}
```

#### 3. **Type-Safe Error Handling**: Educational Error System

**Why This Pattern**: FTS5 applications encounter specific types of errors (validation, database, FTS5-specific). Clear error categorization helps learners understand what went wrong and why.

**Educational Value**: Demonstrates proper error handling for search applications while providing helpful learning hints.

```go
// From errors/errors.go
var (
    ErrValidation = errors.New("validation failed")
    ErrDatabase   = errors.New("database operation failed")
    ErrFTS5       = errors.New("fts5 operation failed")
    ErrNotFound   = errors.New("not found")
)

// Helper functions for creating typed errors
func FTS5f(format string, args ...interface{}) error {
    return fmt.Errorf("%w: "+format, append([]interface{}{ErrFTS5}, args...)...)
}
```

The error system provides educational feedback:

```
FTS5 Error: fts5 operation failed: SQLite not compiled with FTS5 support
Hint: Ensure SQLite is compiled with FTS5 support (go build -tags fts5)
```

### Key Files and Their Roles in Demonstrating Concepts

#### `/database/database.go`: FTS5 Operations Center

This file demonstrates all core FTS5 concepts:

**Virtual Table Creation**:
```go
func (d *Database) CreateTable(ctx context.Context) error {
    query := `
        CREATE VIRTUAL TABLE IF NOT EXISTS documents USING fts5(
            title,
            content,
            category,
            tokenize='unicode61 remove_diacritics 1'
        );`
    _, err := d.db.ExecContext(ctx, query)
    return err
}
```

**BM25 Search Implementation**:
```go
func (d *Database) SearchDocuments(ctx context.Context, query string, options SearchOptions) ([]*models.SearchResult, error) {
    sqlQuery := `
        SELECT rowid, title, content, category, bm25(documents) as score
        FROM documents 
        WHERE documents MATCH ? 
        ORDER BY rank 
        LIMIT ?`
    // Implementation shows proper BM25 score handling
}
```

#### `/models/document.go`: Data Structure Design

Shows how to structure data for FTS5 operations:

```go
type SearchResult struct {
    RowID    int64   // FTS5 rowid for document identification
    Title    string  // Searchable title field
    Content  string  // Main searchable content
    Category string  // Filterable category field
    Score    float64 // BM25 relevance score (negative values)
}
```

#### `/handlers/document.go`: Business Logic Examples

Demonstrates how to use FTS5 operations in application logic:

- Input validation for search queries
- Proper error handling for FTS5 operations  
- Result formatting and presentation
- Flag-based option processing

### Educational Design Decisions

#### Decision 1: In-Memory Default Database
**Reasoning**: Eliminates file system concerns, focuses attention on FTS5 concepts
**Learning Value**: Students can experiment freely without file cleanup

#### Decision 2: Verbose Error Messages
**Reasoning**: Learning applications should explain what went wrong and how to fix it
**Learning Value**: Students learn to recognize and resolve common FTS5 issues

#### Decision 3: Command-Line Interface
**Reasoning**: CLI commands map clearly to FTS5 operations
**Learning Value**: Students can experiment iteratively, building understanding step by step

#### Decision 4: Explicit BM25 Score Display
**Reasoning**: Makes relevance scoring visible and understandable
**Learning Value**: Students see how different queries produce different relevance scores

## Interactive Learning

This section provides step-by-step walkthroughs of the key FTS5 operations, showing you exactly what commands to run and what outputs to expect. Each walkthrough connects directly to the learning objectives.

### Walkthrough 1: FTS5 Virtual Table Creation

**Learning Objective**: Understand FTS5 virtual table creation and configuration

**Command**:
```bash
go run -tags fts5 ./fts5-foundation document create-table --database learning.db
```

**Expected Output**:
```
✓ FTS5 documents table created successfully
```

**What This Demonstrates**:
- The `--database` flag specifies persistent storage (vs default in-memory)
- FTS5 virtual table creation happens instantly (no data migration needed)
- The unicode61 tokenizer is configured with diacritics removal
- Success indicates SQLite was compiled with FTS5 support

**Behind the Scenes**: This command executes:
```sql
CREATE VIRTUAL TABLE IF NOT EXISTS documents USING fts5(
    title,
    content,
    category,
    tokenize='unicode61 remove_diacritics 1'
);
```

**Try This Variation**: 
```bash
# Use verbose flag to see detailed information
go run -tags fts5 ./fts5-foundation document create-table --database learning.db --verbose
```

### Walkthrough 2: Document Insertion and Automatic Indexing

**Learning Objective**: Understand how FTS5 automatically indexes content during insertion

**Commands**:
```bash
# Insert sample documents for experimentation
go run -tags fts5 ./fts5-foundation document batch-insert --database learning.db

# Insert a custom document
go run -tags fts5 ./fts5-foundation document insert \
  --title "Custom Document Title" \
  --content "This document contains specific content about SQLite FTS5 and BM25 scoring algorithms." \
  --category "custom" \
  --database learning.db
```

**Expected Outputs**:
```
Inserting 4 example documents...
✓ Successfully inserted 4 documents

✓ Document inserted successfully with rowid: 5
```

**What This Demonstrates**:
- FTS5 automatically indexes all text content during insertion
- Each document gets a unique `rowid` for identification
- No manual index rebuilding required
- The tokenizer immediately processes the text for searching

**Sample Documents Inserted**:
1. **SQLite FTS5 Full-Text Search** (database category)
2. **Database Indexing Fundamentals** (database category) 
3. **Go Programming Best Practices** (programming category)
4. **Introduction to Algorithms** (programming category)

### Walkthrough 3: Basic Text Search with BM25 Scoring

**Learning Objective**: Understand MATCH queries and BM25 relevance scoring

**Command**:
```bash
go run -tags fts5 ./fts5-foundation document search "database" --database learning.db --scores
```

**Expected Output**:
```
Found 2 document(s) matching: database

--- Result #1 ---
Title: Database Indexing Fundamentals
Category: database
Content: Database indexes are data structures that improve the speed of data retrieval operations...
BM25 Score: -0.0000 (lower is better)

--- Result #2 ---
Title: SQLite FTS5 Full-Text Search  
Category: database
Content: SQLite FTS5 is an SQLite virtual table module that provides full-text search functionality...
BM25 Score: -0.0000 (lower is better)
```

**What This Demonstrates**:
- The MATCH operator finds documents containing "database"
- Results are ranked by BM25 relevance (ORDER BY rank)
- Lower scores indicate better matches (SQLite's negative scoring)
- Content preview shows matching context

**Try These Variations**:
```bash
# Search without scores to see cleaner output
go run -tags fts5 ./fts5-foundation document search "programming" --database learning.db

# Limit results to see top matches only
go run -tags fts5 ./fts5-foundation document search "sqlite" --database learning.db --limit 1

# Search for multi-word phrases
go run -tags fts5 ./fts5-foundation document search "full text search" --database learning.db --scores
```

### Walkthrough 4: Field-Specific Searches

**Learning Objective**: Understand how to search within specific fields using FTS5 syntax

**Command**:
```bash
go run -tags fts5 ./fts5-foundation document search-field "Introduction" "title" --database learning.db
```

**Expected Output**:
```
Found 1 document(s) matching 'Introduction' in field 'title'

--- Result #1 ---
Title: Introduction to Algorithms
Category: programming
Content: Algorithms are step-by-step procedures for solving problems...
```

**What This Demonstrates**:
- Field-specific search using `title:Introduction` MATCH syntax
- Only titles are searched, not content or category
- Useful for finding documents by specific metadata

**Try These Variations**:
```bash
# Search in content field only
go run -tags fts5 ./fts5-foundation document search-field "algorithms" "content" --database learning.db

# Search in category field only  
go run -tags fts5 ./fts5-foundation document search-field "database" "category" --database learning.db
```

### Walkthrough 5: Category-Filtered Searches

**Learning Objective**: Understand how to combine category filtering with text search

**Command**:
```bash
go run -tags fts5 ./fts5-foundation document search-category "programming" "algorithms" --database learning.db
```

**Expected Output**:
```
Found 1 document(s) in category 'programming' matching: algorithms

--- Result #1 ---
Title: Introduction to Algorithms
Category: programming
Content: Algorithms are step-by-step procedures for solving problems...
```

**What This Demonstrates**:
- Combination of category filter with text search
- Uses `category:programming AND algorithms` MATCH syntax
- Narrows results to specific document types

**Try These Variations**:
```bash
# Find database-related documents mentioning indexing
go run -tags fts5 ./fts5-foundation document search-category "database" "indexing" --database learning.db

# Find programming documents mentioning best practices
go run -tags fts5 ./fts5-foundation document search-category "programming" "best practices" --database learning.db
```

### Walkthrough 6: Document Management (CRUD Operations)

**Learning Objective**: Understand how FTS5 handles updates and maintains indexes automatically

**Commands**:
```bash
# List all documents first
go run -tags fts5 ./fts5-foundation document list --database learning.db

# Update a document (this will reindex automatically)
go run -tags fts5 ./fts5-foundation document update 1 \
  --title "Updated: SQLite FTS5 Advanced Features" \
  --content "Updated content about advanced FTS5 features including auxiliary functions." \
  --database learning.db

# Search to see the updated content
go run -tags fts5 ./fts5-foundation document search "auxiliary" --database learning.db

# Delete a document (removes from index automatically)
go run -tags fts5 ./fts5-foundation document delete 1 --database learning.db

# Verify deletion
go run -tags fts5 ./fts5-foundation document list --database learning.db
```

**Expected Outputs**:
```
Documents in FTS5 table:

ID: 1  | Title: SQLite FTS5 Full-Text Search | Category: database
Preview: SQLite FTS5 is an SQLite virtual table module...

ID: 2  | Title: Database Indexing Fundamentals | Category: database
Preview: Database indexes are data structures...

[Additional documents...]

✓ Document updated successfully

Found 1 document(s) matching: auxiliary
--- Result #1 ---
Title: Updated: SQLite FTS5 Advanced Features
[Updated content shown...]

✓ Document deleted successfully

[List shows remaining documents without deleted entry]
```

**What This Demonstrates**:
- FTS5 automatically reindexes content during updates
- Deleted documents are immediately removed from search results
- No manual index maintenance required
- CRUD operations work seamlessly with search functionality

### Understanding Command Output Formats

The CLI provides different output formats to support various learning needs:

#### Text Format (Default)
Human-readable output perfect for learning and experimentation:
```bash
go run -tags fts5 ./fts5-foundation document search "database" --database learning.db
```

#### JSON Format  
Structured output useful for programmatic processing:
```bash
go run -tags fts5 ./fts5-foundation document search "database" --database learning.db --format json
```

#### Verbose Mode
Detailed error information and operation details:
```bash
go run -tags fts5 ./fts5-foundation document search "nonexistent" --database learning.db --verbose
```

## Hands-On Exercises

These exercises are designed to reinforce your understanding through progressively challenging tasks. Work through them in order, as each builds on the previous one.

### Exercise 1: Basic FTS5 Setup (Beginner)

**Objective**: Practice virtual table creation and understand the foundation concepts.

**Tasks**:
1. Create a new database file named `exercise1.db`
2. Create the FTS5 virtual table using the CLI
3. Verify the table exists by attempting to insert a document
4. Use verbose mode to see detailed operation information

**Commands to Try**:
```bash
# Step 1: Create table
go run -tags fts5 ./fts5-foundation document create-table --database exercise1.db --verbose

# Step 2: Insert a test document
go run -tags fts5 ./fts5-foundation document insert \
  --title "My First FTS5 Document" \
  --content "This is my first document in an FTS5 virtual table. It contains searchable text." \
  --category "test" \
  --database exercise1.db

# Step 3: Verify with a search
go run -tags fts5 ./fts5-foundation document search "searchable" --database exercise1.db --scores
```

**Success Criteria**:
- Table creation succeeds without errors
- Document insertion returns a rowid
- Search finds your inserted document
- BM25 score is displayed (negative value)

**Self-Assessment Questions**:
1. What happens if you try to create the table twice?
2. What does the rowid represent in FTS5 context?
3. Why is the BM25 score negative?

### Exercise 2: Search Pattern Exploration (Intermediate)

**Objective**: Master different search patterns and understand their use cases.

**Setup**: Use the batch-insert command to populate your database:
```bash
go run -tags fts5 ./fts5-foundation document batch-insert --database exercise2.db
```

**Tasks**:
1. Perform basic searches for different terms
2. Compare basic search vs field-specific search results
3. Use category filtering to narrow results
4. Experiment with multi-word queries

**Specific Challenges**:

**Challenge 2.1 - Search Comparison**:
```bash
# Basic search for "programming"
go run -tags fts5 ./fts5-foundation document search "programming" --database exercise2.db

# Title-only search for "programming"  
go run -tags fts5 ./fts5-foundation document search-field "programming" "title" --database exercise2.db

# Compare the results. Why are they different?
```

**Challenge 2.2 - Category Filtering**:
```bash
# Find all database-related documents
go run -tags fts5 ./fts5-foundation document search-category "database" "" --database exercise2.db

# Find database documents mentioning "indexing"
go run -tags fts5 ./fts5-foundation document search-category "database" "indexing" --database exercise2.db

# How do the result sets differ?
```

**Challenge 2.3 - Multi-word Queries**:
```bash
# Search for "best practices" (phrase)
go run -tags fts5 ./fts5-foundation document search "best practices" --database exercise2.db --scores

# Compare with individual words
go run -tags fts5 ./fts5-foundation document search "best" --database exercise2.db --scores
go run -tags fts5 ./fts5-foundation document search "practices" --database exercise2.db --scores

# Analyze the scoring differences
```

**Success Criteria**:
- You understand when to use each search pattern
- You can explain why results differ between search types
- You recognize how BM25 scores change with query complexity

**Self-Assessment Questions**:
1. When would you use field-specific search vs basic search?
2. How does category filtering affect BM25 scoring?
3. What happens to relevance scores with multi-word queries?

### Exercise 3: Document Management and Index Behavior (Advanced)

**Objective**: Understand how FTS5 maintains indexes during CRUD operations.

**Setup**: Start with a fresh database and sample documents:
```bash
go run -tags fts5 ./fts5-foundation document create-table --database exercise3.db
go run -tags fts5 ./fts5-foundation document batch-insert --database exercise3.db
```

**Tasks**:

**Task 3.1 - Update Impact Analysis**:
1. Search for "SQLite" and note the results and scores
2. Update one of the matching documents to add more instances of "SQLite"
3. Perform the same search and compare results
4. Analyze how the update affected BM25 scoring

```bash
# Initial search
go run -tags fts5 ./fts5-foundation document search "SQLite" --database exercise3.db --scores

# List documents to find a rowid
go run -tags fts5 ./fts5-foundation document list --database exercise3.db

# Update document (replace ID with actual rowid)
go run -tags fts5 ./fts5-foundation document update [ID] \
  --content "SQLite is amazing. SQLite FTS5 makes SQLite even better for search. SQLite with FTS5 enables powerful text search capabilities." \
  --database exercise3.db

# Search again and compare scores
go run -tags fts5 ./fts5-foundation document search "SQLite" --database exercise3.db --scores
```

**Task 3.2 - Deletion and Search Results**:
1. Perform a search that returns multiple results
2. Delete one of the matching documents  
3. Repeat the search to verify the document is gone
4. Confirm the indexes are automatically cleaned up

```bash
# Search to find documents
go run -tags fts5 ./fts5-foundation document search "programming" --database exercise3.db

# Delete one result (note its rowid first)
go run -tags fts5 ./fts5-foundation document delete [ROWID] --database exercise3.db

# Search again to confirm removal
go run -tags fts5 ./fts5-foundation document search "programming" --database exercise3.db
```

**Task 3.3 - Content Analysis Impact**:
1. Insert documents of varying lengths with the same keywords
2. Observe how document length affects BM25 scoring
3. Experiment with keyword density vs document length

```bash
# Insert short document
go run -tags fts5 ./fts5-foundation document insert \
  --title "Short Doc" \
  --content "Python programming" \
  --category "test" \
  --database exercise3.db

# Insert long document with same keywords
go run -tags fts5 ./fts5-foundation document insert \
  --title "Long Doc" \
  --content "Python programming is a vast field with many aspects. Programming languages like Python offer great flexibility. When programming in Python, developers can build various applications. Python programming encompasses web development, data science, automation, and much more." \
  --category "test" \
  --database exercise3.db

# Compare BM25 scores
go run -tags fts5 ./fts5-foundation document search "Python programming" --database exercise3.db --scores
```

**Success Criteria**:
- You can predict how document changes affect search results
- You understand the relationship between document length and BM25 scoring
- You recognize that FTS5 automatically maintains indexes

**Advanced Self-Assessment Questions**:
1. How does term frequency in updates affect relevance ranking?
2. What happens to search performance as you add/remove documents?
3. Can you predict which document will rank higher based on content analysis?

### Exercise 4: Error Handling and Troubleshooting (Expert)

**Objective**: Experience common error scenarios and learn to resolve them effectively.

**Tasks**:

**Task 4.1 - Validation Error Handling**:
```bash
# Try to insert a document with missing required fields
go run -tags fts5 ./fts5-foundation document insert \
  --title "" \
  --content "Content without title" \
  --database exercise4.db

# Try to search with empty query
go run -tags fts5 ./fts5-foundation document search "" --database exercise4.db --verbose

# Analyze the error messages and their educational value
```

**Task 4.2 - Database Error Scenarios**:
```bash
# Try to search before creating table
go run -tags fts5 ./fts5-foundation document search "test" --database new_exercise.db

# Try to access a read-only location (Linux/Mac)
go run -tags fts5 ./fts5-foundation document create-table --database /root/readonly.db

# Use verbose mode to see detailed error information
```

**Task 4.3 - FTS5 Configuration Testing**:
```bash
# Test without FTS5 build tags (this should fail)
go run ./fts5-foundation document create-table --database test.db

# Compare with proper FTS5 tags
go run -tags fts5 ./fts5-foundation document create-table --database test.db
```

**Success Criteria**:
- You can identify different error categories
- You understand how to resolve common FTS5 issues
- You can use verbose mode effectively for debugging

### Extension Activities for Deeper Exploration

**Activity 1: Custom Document Collections**
Create themed document collections (technical papers, news articles, recipes) and analyze how BM25 scoring behaves with different content types.

**Activity 2: Performance Analysis**
Insert varying numbers of documents (10, 100, 1000) and measure search response time differences.

**Activity 3: Query Complexity Study**
Design increasingly complex queries and study how they affect both result relevance and performance.

**Activity 4: Real-World Integration**
Import actual text files or web content into your FTS5 database and experiment with real search scenarios.

## Self-Assessment

Use these questions to verify your understanding before moving to Phase 2.

### Knowledge Check: Fundamental Concepts

**Question 1**: Virtual Table Creation
```
You need to create an FTS5 table for a recipe database with fields for recipe_name, ingredients, and cuisine_type. Write the CREATE VIRTUAL TABLE statement.

Expected considerations:
- Proper FTS5 syntax
- Unicode tokenizer configuration
- Appropriate field naming
```

**Question 2**: BM25 Score Interpretation
```
You search for "chocolate cake" and get these results:
- Document A: Score -1.2
- Document B: Score -0.8  
- Document C: Score -2.1

Which document is most relevant and why? What factors might cause these score differences?
```

**Question 3**: Search Pattern Selection
```
You have a database of technical documentation and need to:
a) Find all documents with "API" in the title
b) Find database-related documents mentioning "performance"
c) Find any document mentioning "optimization"

Which search pattern would you use for each scenario?
```

### Practical Skills Assessment

**Skill Test 1**: Complete Workflow
Demonstrate a complete FTS5 workflow by:
1. Creating a new database with FTS5 table
2. Inserting 3 documents with different categories
3. Performing 3 different types of searches
4. Updating one document and showing the search impact
5. Properly cleaning up when done

**Skill Test 2**: Error Resolution
Show you can handle these scenarios:
1. Missing build tags error
2. Empty search query validation
3. Database permission issues
4. Non-existent document updates

**Skill Test 3**: Results Analysis
Given search results with BM25 scores, explain:
1. Why the ranking order makes sense
2. How document length affects scoring
3. How term frequency impacts relevance
4. When you might want different ranking

### Advanced Understanding Questions

**Question 4**: Architecture Benefits
Explain why the project uses:
- Global database instances vs dependency injection
- Typed error handling vs generic errors
- Command groups vs flat command structure

**Question 5**: FTS5 vs Traditional Search
Compare FTS5 MATCH queries with traditional SQL LIKE queries:
- Performance differences
- Capability differences  
- When you'd choose each approach

**Question 6**: Index Management
Describe what happens "behind the scenes" when:
- You insert a new document
- You update an existing document's content
- You delete a document
- You perform a search query

### Answer Guidelines

**Scoring Rubric**:
- **Beginner (Foundation)**: Can create tables, insert documents, perform basic searches
- **Intermediate (Proficient)**: Understands BM25 scoring, uses appropriate search patterns
- **Advanced (Expert)**: Explains index behavior, handles errors effectively, designs good queries

**Ready for Phase 2?**
You should be able to:
- ✅ Create and configure FTS5 virtual tables confidently
- ✅ Understand and interpret BM25 scores correctly
- ✅ Choose appropriate search patterns for different scenarios  
- ✅ Handle common errors and troubleshoot issues
- ✅ Explain how FTS5 indexes work conceptually

## Troubleshooting Guide

### Common Error Scenarios and Solutions

#### Error 1: FTS5 Not Available

**Symptom**:
```
FTS5 Error: fts5 operation failed: SQLite not compiled with FTS5 support
Hint: Ensure SQLite is compiled with FTS5 support (go build -tags fts5)
```

**Cause**: Missing `-tags fts5` build flag

**Solution**:
```bash
# Wrong
go run ./fts5-foundation document create-table

# Correct  
go run -tags fts5 ./fts5-foundation document create-table
```

**Prevention**: Always include `-tags fts5` in your run commands.

#### Error 2: Table Not Found

**Symptom**:
```
Database Error: database operation failed: no such table: documents
```

**Cause**: Attempting operations before creating the FTS5 table

**Solution**:
```bash
# First create the table
go run -tags fts5 ./fts5-foundation document create-table --database mydb.db

# Then perform other operations
go run -tags fts5 ./fts5-foundation document search "test" --database mydb.db
```

**Prevention**: Always run `create-table` command first for new databases.

#### Error 3: Empty Search Query

**Symptom**:
```
Validation Error: validation failed: search query cannot be empty
```

**Cause**: Providing empty string as search query

**Solution**:
```bash
# Wrong
go run -tags fts5 ./fts5-foundation document search "" --database mydb.db

# Correct
go run -tags fts5 ./fts5-foundation document search "your query" --database mydb.db
```

**Prevention**: Always provide meaningful search terms.

#### Error 4: Database Permission Issues

**Symptom**:
```
Database Error: database operation failed: unable to open database file
```

**Cause**: Insufficient permissions or invalid file path

**Solutions**:
1. Check file permissions: `ls -la your-database.db`
2. Ensure directory exists: `mkdir -p /path/to/database/`  
3. Use absolute paths: `--database /full/path/to/database.db`
4. Use writable directory: `--database ./writable/path/database.db`

#### Error 5: Invalid Document ID

**Symptom**:
```
Not Found: not found: document with rowid 999
```

**Cause**: Attempting to update/delete non-existent document

**Solution**:
```bash
# List documents to see valid IDs
go run -tags fts5 ./fts5-foundation document list --database mydb.db

# Use valid rowid from the list
go run -tags fts5 ./fts5-foundation document update 1 --title "New Title" --database mydb.db
```

### Performance Troubleshooting

#### Slow Search Performance

**Investigation Steps**:
1. Check database size: Large databases may require optimization
2. Verify FTS5 is being used: Ensure you're using MATCH, not LIKE
3. Monitor query complexity: Very complex queries may be slower

**Solutions**:
- Use field-specific searches when possible
- Limit result count with `--limit` flag
- Consider breaking complex queries into simpler ones

#### Memory Usage with In-Memory Databases

**Symptom**: Application using excessive memory

**Cause**: In-memory databases store all data in RAM

**Solution**: Use file-based databases for large datasets:
```bash
# Instead of default in-memory
go run -tags fts5 ./fts5-foundation document search "query" --database persistent.db
```

### Development Environment Issues

#### Go Module Problems

**Symptom**: Import errors or dependency issues

**Solutions**:
```bash
# Clean module cache
go clean -modcache

# Tidy dependencies
go mod tidy

# Re-download dependencies
go mod download
```

#### SQLite Version Issues

**Check SQLite Version**:
```bash
sqlite3 --version
```

**Update SQLite** (if needed):
- macOS: `brew update && brew upgrade sqlite`
- Ubuntu: `sudo apt update && sudo apt upgrade sqlite3`
- Windows: Download from https://sqlite.org/download.html

### Getting Help

#### Enable Verbose Mode
For any error, try adding `--verbose` to see detailed information:
```bash
go run -tags fts5 ./fts5-foundation document search "query" --database mydb.db --verbose
```

#### Check Application Help
```bash
# General help
go run -tags fts5 ./fts5-foundation --help

# Command-specific help
go run -tags fts5 ./fts5-foundation document --help
go run -tags fts5 ./fts5-foundation document search --help
```

#### Validate Environment
Run a simple test to ensure everything works:
```bash
# Quick validation test
go run -tags fts5 ./fts5-foundation document create-table
go run -tags fts5 ./fts5-foundation document insert --title "Test" --content "Test content" --category "test"
go run -tags fts5 ./fts5-foundation document search "test" --scores
```

If this works, your environment is correctly configured.

## Next Steps

Congratulations! You've completed the SQLite FTS5 Foundation phase. Here's how to continue your learning journey:

### Phase 2: BM25 Fundamentals

**What You'll Learn**: 
- Advanced BM25 scoring concepts
- Corpus analysis and document statistics
- Relevance tuning and optimization
- Comparative scoring analysis

**How to Proceed**:
```bash
cd ../02-bm25-fundamentals
go run -tags fts5 ./bm25-fundamentals --help
```

### Recommended Learning Path

1. **Review**: Ensure you can complete all exercises in this guide
2. **Practice**: Create your own document collections and experiment
3. **Research**: Read about BM25 algorithm details (linked in Resources)
4. **Advance**: Move to Phase 2 when you're comfortable with all Foundation concepts

### Key Concepts to Remember

- **FTS5 Virtual Tables**: Foundation for all text search operations
- **BM25 Scoring**: Lower (more negative) scores = better matches in SQLite
- **MATCH Operator**: Your gateway to FTS5's powerful search capabilities
- **Automatic Indexing**: FTS5 maintains indexes transparently during CRUD operations
- **Search Patterns**: Choose the right pattern for your use case

### Resources for Continued Learning

**Official Documentation**:
- [SQLite FTS5 Documentation](https://www.sqlite.org/fts5.html)
- [BM25 Algorithm Details](https://en.wikipedia.org/wiki/Okapi_BM25)

**Advanced Topics to Explore**:
- Custom tokenizers for different languages
- FTS5 auxiliary functions (highlight, snippet)
- Performance optimization techniques
- Integration patterns with web applications

### Building on This Foundation

The concepts you've learned form the foundation for advanced search applications:

- **Contextual Memory Systems**: Using FTS5 for AI knowledge retrieval
- **Document Management Systems**: Enterprise search solutions
- **Content Discovery**: Recommendation and search systems
- **Data Analysis**: Text mining and information retrieval

The solid foundation you've built here will serve you well as you tackle increasingly sophisticated search challenges in the phases ahead.

---

**Completed Phase 1: SQLite FTS5 Foundation** ✅

Ready to explore advanced BM25 concepts in Phase 2? Your journey into powerful, intelligent text search continues!