# BM25 Fundamentals Learning Guide

## Table of Contents

1. [Learning Objectives Overview](#learning-objectives-overview)
2. [Conceptual Foundations](#conceptual-foundations)
3. [Project Infrastructure](#project-infrastructure)
4. [Interactive Learning](#interactive-learning)
5. [Hands-On Exercises](#hands-on-exercises)
6. [Advanced Concepts](#advanced-concepts)
7. [Educational Outcomes](#educational-outcomes)
8. [Troubleshooting Guide](#troubleshooting-guide)

## Learning Objectives Overview

### Primary Learning Goals

This phase transforms your understanding from basic FTS5 operations to sophisticated BM25 scoring mastery. By completing this phase, you will achieve:

**1. BM25 Scoring System Mastery**
- Deep understanding of SQLite FTS5's unique negative scoring approach
- Ability to interpret and compare BM25 scores correctly
- Knowledge of how SQLite's implementation differs from standard BM25

**2. Document Length Normalization Effects**
- Understanding how the 'b' parameter (0.75) affects scoring
- Recognition of length bias and its mitigation strategies
- Practical experience with short vs. long document performance

**3. Column Weighting Optimization**
- Multi-field search relevance tuning capabilities
- Strategic weighting decisions for different content types
- Comparative analysis of weighting strategies

**4. Score Distribution Analysis**
- Statistical pattern recognition in search results
- Percentile-based quality assessment
- Outlier identification and relevance threshold setting

**5. Ranking Optimization Strategies**
- Evidence-based ranking improvements
- A/B testing methodologies for search relevance
- Performance vs. relevance trade-off analysis

### Prerequisites from Phase 1

Before starting this phase, ensure you have mastered:
- Basic FTS5 virtual table creation and operations
- MATCH vs LIKE operator differences
- CommandGroup pattern and CLI architecture
- Type-safe error handling patterns
- Go SQLite driver configuration with FTS5 enabled

### Position in Learning Roadmap

**Phase 1 Foundation** → **Phase 2 BM25 Fundamentals** → **Phase 3 Query Operations**

This phase bridges the gap between basic FTS5 operations and advanced query patterns. You'll use BM25 knowledge throughout the remaining phases for relevance optimization and ranking strategy development.

## Conceptual Foundations

### Understanding BM25: The Restaurant Review Analogy

Imagine you're reading restaurant reviews to find the best dining experience. BM25 works similarly to how you'd naturally evaluate reviews:

**Term Frequency (TF)**: If a review mentions "delicious" multiple times, it's probably more positive than one mentioning it once. However, mentioning "delicious" 10 times isn't necessarily 10 times better than mentioning it twice—there's diminishing returns.

**Inverse Document Frequency (IDF)**: The word "restaurant" appears in every review, so it's not very helpful for ranking. But "molecular gastronomy" appears rarely, making it highly significant when present.

**Document Length Normalization**: A 50-word review mentioning "excellent" twice is more significant than a 500-word review mentioning it twice. BM25 compensates for this natural dilution effect.

### SQLite's Negative Scoring System

Unlike standard BM25 implementations that return positive scores, SQLite FTS5 returns negative values:

```
Standard BM25:  Higher positive scores = Better relevance
SQLite FTS5:    Higher negative scores = Better relevance (closer to 0)

Example Rankings:
-1.2 ranks HIGHER than -3.5
-0.8 ranks HIGHER than -1.2
-5.1 ranks LOWER than -2.3
```

**Why Negative Scores?**
SQLite inverts BM25 scores to optimize internal sorting algorithms. This is purely an implementation detail—the relative rankings remain mathematically correct.

### BM25 Parameters in SQLite

SQLite FTS5 uses hardcoded parameters that have been empirically validated:

**k1 = 1.2 (Term Frequency Saturation)**
- Controls how quickly additional term occurrences provide diminishing returns
- Higher k1 = More weight to repeated terms
- Lower k1 = Quicker saturation of term frequency benefits

**b = 0.75 (Length Normalization)**
- Controls how much document length affects scoring
- b = 0: No length normalization (longer docs have advantage)
- b = 1: Full length normalization (proportional scoring)
- b = 0.75: Balanced approach (SQLite's choice)

### Multi-Field Document Relevance

Consider searching for "database performance" across these document fields:

```
Document A:
Title: "Database Performance Optimization"
Content: "Various techniques for improving query speed..."
Category: "Performance"

Document B:
Title: "Software Development Best Practices"
Content: "Database performance is crucial for scalable applications..."
Category: "Development"
```

Without column weighting, both documents score equally if they contain the search terms the same number of times. Column weighting lets you specify that title matches are more important:

```bash
# Equal weighting (default)
--title-weight 1.0 --content-weight 1.0 --category-weight 1.0

# Emphasize title matches
--title-weight 3.0 --content-weight 1.0 --category-weight 1.0
```

This mathematical weighting is applied to the BM25 calculation, effectively multiplying the relevance contribution of title terms by 3.

### Score Distribution Patterns

Real-world search results typically follow predictable patterns:

**Steep Drop-off**: Most results cluster in poor relevance ranges, with few highly relevant results
**Long Tail**: Many marginally relevant results create extended distributions
**Outliers**: Exceptional matches or irrelevant results appear at distribution extremes

Understanding these patterns helps you:
- Set relevance thresholds for result filtering
- Identify when search quality is poor
- Optimize ranking strategies based on distribution characteristics

## Project Infrastructure

### Educational Architecture Design

The BM25 fundamentals tool follows sophisticated patterns designed to support learning:

**Global Instance Pattern**
```go
// database/database.go
var Instance *Database  // Global database access

// config/config.go  
var App *Config        // Global configuration access
```

This pattern eliminates complex dependency injection while maintaining clean separation of concerns—ideal for educational CLI tools where simplicity aids understanding.

**Stateless Handler Architecture**
```go
// handlers/search.go
var Search SearchHandler

type SearchHandler struct{} // No state stored

func (h *SearchHandler) HandleQuery(cmd *cobra.Command, args []string) error {
    // Access global instances directly
    results, err := database.Instance.SearchBM25(ctx, options)
    if config.App.Verbose {
        // Show detailed explanations
    }
}
```

Handlers remain stateless and testable while accessing global configuration and database instances. This pattern scales well for CLI applications and simplifies educational examples.

**Factory Pattern for Commands**
```go
// commands/search.go
var Search = newSearchGroup() // Public instance

func newSearchGroup() *CommandGroup {
    // Command creation logic encapsulated
    // Prevents naming conflicts between files
}
```

### Key Educational Components

**1. Statistical Analysis Engine** (`handlers/search.go`)
Demonstrates advanced Go patterns for data analysis:
- Score distribution calculations
- Percentile analysis algorithms
- Statistical summary generation

**2. Visualization System** (`handlers/visualize.go`)
Shows practical data visualization techniques:
- ASCII chart generation using asciigraph
- Histogram creation and display
- Multi-category comparison charts

**3. Type-Safe Error Handling** (`errors/errors.go`)
Educational error management patterns:
- Sentinel error types for different failure modes
- Contextual error messages with hints
- Automatic verbose/simple mode switching

**4. Configuration Management** (`config/config.go`)
Demonstrates professional configuration practices:
- Viper integration with struct mapping
- Default value management
- Validation and error handling

### File Organization for Learning

```
bm25-fundamentals/
├── main.go                 # Clean entry point
├── commands/              # CLI interface layer
│   ├── command_group.go   # Hierarchical organization pattern
│   ├── root.go           # Global flag setup and initialization
│   ├── corpus.go         # Corpus management commands
│   ├── search.go         # Search operation commands
│   └── visualize.go      # Visualization commands
├── handlers/             # Business logic layer (stateless)
│   ├── corpus.go         # Corpus generation and management
│   ├── search.go         # BM25 search operations
│   └── visualize.go      # Data visualization
├── models/               # Data structures
│   ├── corpus.go         # Corpus-related types
│   ├── document.go       # Document representation
│   ├── search.go         # Search result types
│   └── analysis.go       # Statistical analysis types
├── database/             # Data persistence layer
│   └── database.go       # SQLite FTS5 operations
├── config/               # Configuration management
│   └── config.go         # Application configuration
└── errors/               # Error handling system
    └── errors.go         # Type-safe error definitions
```

This organization demonstrates clean architecture principles while keeping related educational concepts together.

## Interactive Learning

### Command Execution Walkthroughs

#### Walkthrough 1: Basic BM25 Search

**Step 1: Generate Sample Corpus**
```bash
go run -tags "fts5" . corpus generate --size 50 --database tutorial.db
```

**Expected Output:**
```
Generating 50 documents in categories: [Database Web Development Performance Security Machine Learning]
✓ Generated 50 documents with average length 156 words
✓ FTS5 indexes created for full-text search
Database: tutorial.db
```

**Learning Points:**
- Corpus generation creates realistic document distributions
- FTS5 indexes are automatically maintained
- Document length varies naturally (affects BM25 scoring)

**Step 2: Perform Basic Search**
```bash
go run -tags "fts5" . search query --query "database optimization" --database tutorial.db
```

**Expected Output:**
```
Search Results for: "database optimization"
Found 8 results in 2.3ms

[1] Score: -1.23 | Database Performance Tuning Guide
    Category: Database | Length: 145 words
    Content: Complete guide to optimizing database queries...

[2] Score: -2.18 | SQL Query Optimization Techniques  
    Category: Database | Length: 203 words
    Content: Advanced strategies for improving query performance...

[3] Score: -3.45 | Web Application Database Design
    Category: Web Development | Length: 178 words
    Content: Best practices for database schema optimization...
```

**Learning Points:**
- Negative scores where -1.23 ranks higher than -2.18
- Document length shown to understand normalization effects
- Category information reveals content context
- Execution time demonstrates FTS5 performance

#### Walkthrough 2: Column Weighting Analysis

**Step 1: Search with Default Weights**
```bash
go run -tags "fts5" . search query --query "performance" --database tutorial.db --max-results 5
```

**Step 2: Search with Title Emphasis**
```bash
go run -tags "fts5" . search query --query "performance" --title-weight 3.0 --content-weight 1.0 --database tutorial.db --max-results 5
```

**Learning Points:**
- Compare ranking differences between the two results
- Documents with "performance" in titles will rank higher in step 2
- Score values change but negative ordering remains consistent
- Column weighting affects relevance calculation, not just sorting

#### Walkthrough 3: Score Distribution Visualization

**Step 1: Generate Score Distribution**
```bash
go run -tags "fts5" . visualize distribution --query "database" --database tutorial.db
```

**Expected Output:**
```
Score Distribution for: "database"
Query returned 15 results

 3.00 ┤        ╭─╮
 2.50 ┤        │ │
 2.00 ┤     ╭──╯ ╰─╮
 1.50 ┤     │      │
 1.00 ┤  ╭──╯      ╰─╮
 0.50 ┤  │          │
 0.00 ┼──╯          ╰────
     -6 -5 -4 -3 -2 -1 -0

Distribution Statistics:
  • Best Score: -0.82 (highest relevance)
  • Worst Score: -5.43 (lowest relevance)
  • Mean: -2.67, Median: -2.34
  • Standard Deviation: 1.23
```

**Learning Points:**
- Histogram shows score clustering patterns
- Most results cluster around mean values
- Few highly relevant results (scores near 0)
- Wide distribution indicates diverse relevance levels

#### Walkthrough 4: Statistical Analysis

**Step 1: Generate Search Statistics**
```bash
go run -tags "fts5" . search stats --query "optimization" --database tutorial.db
```

**Expected Output:**
```
Search Statistics for: "optimization"

Query Performance:
  • Total Results: 12
  • Execution Time: 1.8ms
  • Average Processing: 0.15ms per result

Score Analysis:
  • Score Range: -0.94 to -4.67 (spread: 3.73)
  • Mean: -2.31, Median: -2.18
  • Standard Deviation: 1.05
  • Coefficient of Variation: 45.5%

Percentile Distribution:
  • 90th percentile: -1.23 (top 10% of results)
  • 75th percentile: -1.67 (top quartile)
  • 50th percentile: -2.18 (median result)
  • 25th percentile: -3.45 (bottom quartile)

Category Breakdown:
  • Database: 5 results (41.7%)
  • Performance: 3 results (25.0%)
  • Web Development: 2 results (16.7%)
  • Security: 2 results (16.7%)
```

**Learning Points:**
- Statistical summary provides search quality assessment
- Percentile analysis identifies relevance thresholds
- Category breakdown shows content distribution
- Coefficient of variation indicates result diversity

### Troubleshooting Common Issues

#### Issue 1: No Search Results

**Symptoms:**
```bash
go run -tags "fts5" . search query --query "machine learning" --database tutorial.db
# Returns: No results found for query: "machine learning"
```

**Diagnosis Steps:**
1. Check corpus size: `go run -tags "fts5" . corpus stats --database tutorial.db`
2. Try broader queries: `go run -tags "fts5" . search query --query "machine" --database tutorial.db`
3. Enable verbose mode: `go run -tags "fts5" . search query --query "machine learning" --verbose --database tutorial.db`

**Learning Points:**
- FTS5 requires exact token matches (not partial)
- Unicode61 tokenizer affects how terms are indexed
- Verbose mode provides debugging information

#### Issue 2: Unexpected Score Rankings

**Symptoms:**
A shorter document ranks lower than expected despite containing more query terms.

**Diagnosis Approach:**
```bash
# Use explain command to see scoring details
go run -tags "fts5" . search explain --query "database performance" --database tutorial.db
```

**Learning Points:**
- BM25 balances term frequency, document length, and term rarity
- Multiple factors contribute to final score
- Explain output shows individual term contributions

#### Issue 3: Visualization Display Issues

**Symptoms:**
Charts appear malformed or cut off in terminal output.

**Solutions:**
1. Ensure terminal width is adequate (80+ characters)
2. Use smaller bucket counts: `--buckets 10`
3. Try different output formats: `--format json`

## Hands-On Exercises

### Exercise 1: Document Length Impact Investigation

**Objective**: Understand how document length affects BM25 scoring.

**Preparation:**
```bash
# Generate a corpus for testing
go run -tags "fts5" . corpus generate --size 100 --database length_study.db
```

**Exercise Steps:**

1. **Search for a common term:**
   ```bash
   go run -tags "fts5" . search query --query "database" --database length_study.db --max-results 10
   ```

2. **Analyze the results:**
   - Record the document lengths and scores
   - Identify patterns between length and score

3. **Generate detailed statistics:**
   ```bash
   go run -tags "fts5" . search stats --query "database" --database length_study.db
   ```

4. **Use the explain command:**
   ```bash
   go run -tags "fts5" . search explain --query "database" --database length_study.db --max-results 5
   ```

**Self-Assessment Questions:**
- Do longer documents consistently score lower than shorter ones?
- How does term frequency interact with document length?
- What happens when a long document has many term occurrences?

**Expected Learning Outcomes:**
- Understanding of length normalization effects
- Recognition of BM25's balanced approach to length bias
- Ability to predict ranking behavior based on document characteristics

### Exercise 2: Column Weighting Optimization

**Objective**: Master multi-field relevance tuning for different use cases.

**Scenario**: You're building a technical documentation search where title matches should be significantly more important than content matches, but category matches should be de-emphasized.

**Exercise Steps:**

1. **Establish baseline performance:**
   ```bash
   go run -tags "fts5" . search query --query "optimization" --database length_study.db --max-results 10
   ```

2. **Test title emphasis strategy:**
   ```bash
   go run -tags "fts5" . search query --query "optimization" --title-weight 3.0 --content-weight 1.0 --category-weight 0.5 --database length_study.db --max-results 10
   ```

3. **Compare strategies side-by-side:**
   ```bash
   go run -tags "fts5" . search compare --query "optimization" --compare-weights "title:3.0,content:1.0,category:0.5" --database length_study.db
   ```

4. **Experiment with different weightings:**
   - Try extreme title weighting: `--title-weight 5.0 --content-weight 1.0`
   - Try content emphasis: `--title-weight 1.0 --content-weight 2.0`
   - Try balanced approach: `--title-weight 2.0 --content-weight 1.5 --category-weight 1.0`

**Progressive Challenges:**

**Beginner**: Find optimal weights for a blog search (titles important, categories relevant)
**Intermediate**: Optimize for an e-commerce search (titles crucial, content moderate, categories minimal)
**Advanced**: Design weights for a code search (titles moderate, content high, categories ignored)

**Self-Assessment Questions:**
- How do different weightings change the top 5 results?
- Which weighting strategy produces the most relevant results for your use case?
- What trade-offs exist between title and content emphasis?

### Exercise 3: Score Distribution Analysis

**Objective**: Learn to assess search quality through statistical analysis.

**Exercise Steps:**

1. **Generate distributions for different query types:**
   ```bash
   # Specific technical term
   go run -tags "fts5" . visualize distribution --query "SQL" --database length_study.db

   # Common general term  
   go run -tags "fts5" . visualize distribution --query "performance" --database length_study.db

   # Multi-word query
   go run -tags "fts5" . visualize distribution --query "database optimization" --database length_study.db
   ```

2. **Compare category performance:**
   ```bash
   go run -tags "fts5" . visualize categories --query "performance" --database length_study.db
   ```

3. **Analyze score ranges:**
   ```bash
   go run -tags "fts5" . visualize range --query "optimization" --database length_study.db
   ```

**Analysis Questions:**
- Which query type produces the tightest score distribution?
- Which categories perform best for different query types?
- How do multi-word queries change distribution patterns?

**Advanced Extension:**
Export results to CSV and create detailed statistical analysis:
```bash
go run -tags "fts5" . search stats --query "performance" --format csv --database length_study.db > performance_stats.csv
```

### Exercise 4: Real-World Optimization Scenario

**Objective**: Apply BM25 knowledge to solve a practical search relevance problem.

**Scenario**: You're optimizing search for a technical blog with these content types:
- **Tutorials**: Step-by-step guides (should rank high for how-to queries)
- **Reference**: API documentation (should rank high for specific technical terms)
- **Opinion**: Analysis articles (should rank moderately for general topics)

**Challenge Steps:**

1. **Generate a specialized corpus:**
   ```bash
   go run -tags "fts5" . corpus generate --size 150 --database blog_optimization.db
   ```

2. **Test different query scenarios:**
   ```bash
   # How-to query (should favor tutorials)
   go run -tags "fts5" . search query --query "how to optimize" --database blog_optimization.db

   # Technical reference query (should favor reference docs)  
   go run -tags "fts5" . search query --query "API documentation" --database blog_optimization.db

   # General topic query (mixed results expected)
   go run -tags "fts5" . search query --query "performance" --database blog_optimization.db
   ```

3. **Design and test weighting strategies:**
   - Create weights for tutorial emphasis
   - Create weights for technical reference emphasis
   - Create balanced weights for general browsing

4. **Validate your strategies:**
   Use the compare command to evaluate different approaches and visualize results to confirm improvements.

**Success Criteria:**
- Tutorials rank higher for procedural queries
- Reference documents rank higher for technical terms
- Overall search quality improves as measured by score distribution analysis

## Advanced Concepts

### BM25 Algorithm Deep Dive

The complete BM25 formula reveals why certain ranking behaviors occur:

```
BM25(D,Q) = Σ IDF(qi) × f(qi,D) × (k1 + 1) / (f(qi,D) + k1 × (1 - b + b × |D|/avgdl))

Where:
- f(qi,D) = frequency of term qi in document D
- |D| = length of document D in words
- avgdl = average document length in collection
- k1 = 1.2 (controls term frequency saturation)
- b = 0.75 (controls length normalization)
- IDF(qi) = log((N - df(qi) + 0.5) / (df(qi) + 0.5))
```

**Key Insights:**
- Term frequency has diminishing returns due to the saturation factor
- Document length normalization prevents bias toward longer documents
- IDF ensures rare terms have higher impact than common ones
- SQLite inverts the final score for internal optimization

### Performance Characteristics

Understanding BM25 performance helps with system design:

**Index Structure**: FTS5 maintains inverted indexes for each term
**Memory Usage**: Columnsize data adds storage overhead but enables BM25
**Query Complexity**: Multi-term queries require intersection of posting lists
**Update Performance**: Document modifications require index updates

**Optimization Strategies:**
1. **Batch Operations**: Group document updates for efficiency
2. **Selective Indexing**: Only index fields that need full-text search
3. **Query Design**: Prefer specific terms over common ones for better performance

### Column Weighting Mathematics

When you specify column weights, SQLite applies them during BM25 calculation:

```sql
-- Default equal weighting
SELECT bm25(documents_fts) FROM documents_fts WHERE documents_fts MATCH 'query'

-- Custom column weighting  
SELECT bm25(documents_fts, 3.0, 1.0, 0.5) FROM documents_fts WHERE documents_fts MATCH 'query'
```

The weights multiply the contribution of each column's terms:
- Title matches get 3x relevance boost
- Content matches get standard relevance
- Category matches get 0.5x relevance reduction

### Integration with Larger Systems

BM25 fundamentals prepare you for production scenarios:

**Caching Strategies**: Pre-compute common query results
**Result Filtering**: Use BM25 scores to set relevance thresholds
**Multi-Index Operations**: Combine FTS5 with other SQLite indexes
**Real-Time Updates**: Handle document modifications efficiently

## Educational Outcomes

### Competency Assessment

After completing this phase, you should demonstrate:

**1. BM25 Score Interpretation**
- Correctly identify higher-ranking results from negative scores
- Explain why specific documents rank where they do
- Predict ranking changes based on document modifications

**2. Optimization Decision Making**
- Choose appropriate column weights for different use cases
- Analyze score distributions to assess search quality
- Design experiments to validate ranking improvements

**3. Technical Implementation**
- Configure SQLite FTS5 for BM25 scoring
- Implement multi-field search with custom weights
- Generate statistical analysis of search results

**4. Problem-Solving Abilities**
- Diagnose poor search relevance issues
- Design A/B tests for ranking strategies
- Optimize search performance for specific requirements

### Knowledge Validation Exercises

**Quick Assessment Questions:**

1. If Document A scores -1.8 and Document B scores -2.3, which ranks higher and why?
2. How would increasing the title weight from 1.0 to 3.0 affect a document with query terms in the title?
3. What statistical measures help determine if search results have good relevance diversity?
4. Why might a shorter document with fewer term matches outrank a longer document with more matches?

**Practical Validation:**

1. **Explain Exercise**: Use the explain command to describe why the top result ranks #1
2. **Optimization Challenge**: Improve search relevance for a specific query using column weights
3. **Analysis Task**: Identify whether a set of search results has good or poor quality based on score distribution

### Preparation for Next Phases

This BM25 foundation enables advanced learning objectives:

**Phase 3 - Query Operations**: Understanding BM25 helps evaluate complex query performance
**Phase 4 - Ranking and Relevance**: BM25 knowledge is essential for custom ranking strategies  
**Phase 5 - Advanced Features**: BM25 insights inform optimization decisions
**Phase 6 - Integration Patterns**: BM25 competency enables production system design

### Common Misconceptions Clarified

**Misconception**: "Higher BM25 scores always mean better relevance"
**Reality**: SQLite uses negative scores; -1.0 ranks higher than -3.0

**Misconception**: "Document length doesn't matter in modern search"
**Reality**: BM25 length normalization is crucial for fair ranking across different document sizes

**Misconception**: "More term occurrences always improve ranking"
**Reality**: BM25 has saturation—additional occurrences provide diminishing returns

**Misconception**: "Column weighting is just a multiplier applied after scoring"
**Reality**: Weights are integrated into the BM25 calculation itself, affecting the mathematical relevance computation

## Troubleshooting Guide

### Common Issues and Solutions

#### FTS5 Build Problems

**Symptom**: `sqlite3: SQL logic error: no such module: fts5`

**Solutions:**
1. Verify build tags: `go run -tags "fts5" .`
2. Check SQLite compilation: `sqlite3 :memory: "PRAGMA compile_options;" | grep FTS5`
3. Rebuild with proper flags: `CGO_ENABLED=1 go build -tags "fts5"`

#### Search Quality Issues

**Symptom**: Search results seem irrelevant or poorly ranked

**Diagnosis Approach:**
1. Check corpus diversity: `go run -tags "fts5" . corpus stats --database your.db`
2. Analyze score distributions: `go run -tags "fts5" . visualize distribution --query "your_query" --database your.db`
3. Use explain mode: `go run -tags "fts5" . search explain --query "your_query" --database your.db`

**Common Fixes:**
- Increase corpus size for better IDF calculations
- Adjust column weights based on content structure
- Use more specific query terms

#### Performance Issues

**Symptom**: Slow search execution times

**Investigation Steps:**
1. Enable verbose timing: `go run -tags "fts5" . search query --query "test" --verbose --database your.db`
2. Check database size: `ls -lh your.db`
3. Profile query complexity: Use EXPLAIN QUERY PLAN in SQLite

**Optimization Techniques:**
- Reduce result limits for testing
- Use more selective query terms
- Consider in-memory databases for development

#### Visualization Problems

**Symptom**: Charts don't display correctly or appear truncated

**Solutions:**
1. Ensure terminal width ≥ 80 characters
2. Reduce bucket count: `--buckets 8`
3. Use alternative formats: `--format json` or `--format csv`
4. Check for Unicode support in terminal

### Getting Help

**Built-in Help System:**
```bash
# General help
go run -tags "fts5" . --help

# Command-specific help
go run -tags "fts5" . search --help
go run -tags "fts5" . visualize distribution --help
```

**Verbose Mode for Debugging:**
```bash
# Add --verbose to any command for detailed explanations
go run -tags "fts5" . search query --query "test" --verbose --database your.db
```

**Export for External Analysis:**
```bash
# Export results in different formats
go run -tags "fts5" . search stats --query "test" --format csv --database your.db
go run -tags "fts5" . search query --query "test" --format json --database your.db
```

---

**Congratulations!** Upon completing this guide and exercises, you'll have mastered BM25 fundamentals and be prepared for advanced FTS5 query operations in Phase 3. Your deep understanding of scoring mechanics will serve as the foundation for all future search relevance work.