# Phase 2: BM25 Fundamentals Learning Tool

A comprehensive CLI tool for mastering BM25 scoring mechanics and interpretation in SQLite FTS5. This educational tool demonstrates advanced full-text search concepts through hands-on experimentation.

## Learning Objectives

By the end of this phase, you will understand:

- **BM25 Scoring System**: How SQLite FTS5 implements BM25 with negative scores (lower = better relevance)
- **Document Length Normalization**: Impact of document length on search relevance with parameters k1=1.2, b=0.75
- **Column Weighting**: How to tune multi-field search relevance for different use cases
- **Score Distribution Analysis**: Statistical patterns in search result scoring
- **Ranking Optimization**: Strategies for improving search result quality

## Quick Start

### Prerequisites
- Go 1.24+ with SQLite FTS5 support
- Build with FTS5 tags: `go run -tags "fts5" .`

### Basic Usage

```bash
# Generate sample corpus
go run -tags "fts5" . corpus generate --size 100 --database corpus.db

# Basic search
go run -tags "fts5" . search query --query "database optimization" --database corpus.db

# Search with column weights (emphasize title matches)
go run -tags "fts5" . search query --query "performance" --title-weight 3.0 --content-weight 1.0 --database corpus.db

# Compare different ranking strategies
go run -tags "fts5" . search compare --query "SQL" --compare-weights "title:2.0,content:1.0" --database corpus.db

# Visualize score distributions
go run -tags "fts5" . visualize distribution --query "performance" --database corpus.db
```

## Command Reference

### Corpus Management

#### `corpus generate`
Generate sample documents for experimentation.

```bash
go run -tags "fts5" . corpus generate --size 50 --database test.db
```

**Key Learning**: Understand how corpus size and diversity affect BM25 scoring patterns.

#### `corpus stats`
View corpus statistics including category distribution and document characteristics.

```bash
go run -tags "fts5" . corpus stats --database test.db
```

#### `corpus clear`
Remove all documents from the corpus.

```bash
go run -tags "fts5" . corpus clear --database test.db
```

### Search Operations

#### `search query`
Perform basic BM25 search with optional column weighting.

```bash
# Basic search
go run -tags "fts5" . search query --query "machine learning" --database test.db

# With column weights
go run -tags "fts5" . search query --query "SQL optimization" \
  --title-weight 2.0 --content-weight 1.0 --category-weight 0.5 --database test.db
```

**Key Learning**: See how BM25 scores are negative values where -1.5 ranks higher than -3.2.

#### `search stats`
Generate statistical analysis of search results.

```bash
go run -tags "fts5" . search stats --query "database" --database test.db
```

**Key Learning**: Understand score distribution patterns and percentile analysis.

#### `search explain`
Get detailed BM25 score explanations with term analysis.

```bash
go run -tags "fts5" . search explain --query "optimization performance" --database test.db
```

**Key Learning**: See how individual terms contribute to final BM25 scores.

#### `search compare`
Compare ranking results between default and custom weighted strategies.

```bash
go run -tags "fts5" . search compare --query "SQL database" \
  --compare-weights "title:3.0,content:1.0" --database test.db
```

**Key Learning**: Observe how column weighting changes document rankings.

### Visualization

#### `visualize distribution`
Display score distribution histogram.

```bash
go run -tags "fts5" . visualize distribution --query "performance" --database test.db
```

**Key Learning**: See patterns in how BM25 scores cluster across search results.

#### `visualize categories`
Compare score distributions across document categories.

```bash
go run -tags "fts5" . visualize categories --query "optimization" --database test.db
```

**Key Learning**: Understand how different content types perform for specific queries.

#### `visualize range`
Analyze score ranges with percentile breakdowns.

```bash
go run -tags "fts5" . visualize range --query "database" --database test.db
```

**Key Learning**: Identify outliers and understand score variance patterns.

### Output Formats

All commands support multiple output formats:

```bash
# JSON format (machine readable)
go run -tags "fts5" . search query --query "data" --format json --database test.db

# CSV format (spreadsheet compatible)
go run -tags "fts5" . search stats --query "SQL" --format csv --database test.db

# Text format (human readable, default)
go run -tags "fts5" . search query --query "optimization" --database test.db
```

## BM25 Fundamentals

### Understanding Negative Scores

SQLite FTS5 returns **negative BM25 scores** where:
- **Higher scores** (closer to 0) = **better relevance**: -1.5 ranks higher than -3.2
- **Lower scores** (more negative) = **worse relevance**: -5.8 ranks lower than -2.1

### BM25 Parameters in SQLite

SQLite FTS5 uses fixed parameters:
- **k1 = 1.2**: Controls term frequency saturation
- **b = 0.75**: Controls document length normalization effect

### Column Weighting Syntax

```bash
# Equal weights (default)
--title-weight 1.0 --content-weight 1.0 --category-weight 1.0

# Emphasize title matches
--title-weight 3.0 --content-weight 1.0

# De-emphasize categories
--title-weight 1.0 --content-weight 1.0 --category-weight 0.5

# Compare syntax (for search compare command)
--compare-weights "title:2.0,content:1.0,category:0.5"
```

## Learning Experiments

### Experiment 1: Document Length Impact

```bash
# Generate corpus and search for common terms
go run -tags "fts5" . corpus generate --size 100 --database exp1.db
go run -tags "fts5" . search explain --query "performance" --database exp1.db

# Observe: Longer documents typically score lower due to length normalization
```

### Experiment 2: Column Weight Optimization

```bash
# Test different weighting strategies
go run -tags "fts5" . search compare --query "SQL database" --compare-weights "title:1.0,content:1.0" --database exp1.db
go run -tags "fts5" . search compare --query "SQL database" --compare-weights "title:3.0,content:1.0" --database exp1.db
go run -tags "fts5" . search compare --query "SQL database" --compare-weights "title:0.5,content:2.0" --database exp1.db

# Observe: How title vs content emphasis changes rankings
```

### Experiment 3: Score Distribution Patterns

```bash
# Analyze different query types
go run -tags "fts5" . visualize distribution --query "database" --database exp1.db
go run -tags "fts5" . visualize distribution --query "optimization" --database exp1.db
go run -tags "fts5" . visualize distribution --query "machine learning" --database exp1.db

# Observe: How query specificity affects score distributions
```

### Experiment 4: Category Performance Analysis

```bash
# Compare how different categories perform
go run -tags "fts5" . visualize categories --query "performance" --database exp1.db
go run -tags "fts5" . visualize categories --query "development" --database exp1.db

# Observe: Which content categories rank best for different query types
```

## Advanced Usage

### Persistent vs In-Memory Databases

```bash
# Persistent database (recommended for experiments)
--database corpus.db

# In-memory database (faster, data lost on exit)
--database ":memory:"
```

### Verbose Mode

```bash
# Enable detailed explanations and debugging info
--verbose
```

### Large Corpus Experiments

```bash
# Generate larger corpus for more realistic distributions
go run -tags "fts5" . corpus generate --size 500 --database large.db

# Analyze with higher result limits
go run -tags "fts5" . search stats --query "data" --max-results 100 --database large.db
```

## Key Concepts Demonstrated

### 1. **Inverse Document Frequency (IDF)**
- Rare terms score higher than common terms
- Observe in `search explain` output

### 2. **Term Frequency (TF)**
- Multiple occurrences of query terms boost relevance
- Balanced by document length normalization

### 3. **Length Normalization**
- Longer documents are penalized to prevent bias
- Controlled by parameter b=0.75 in SQLite

### 4. **Multi-Field Relevance**
- Different fields can have different importance
- Title matches often more valuable than content matches

### 5. **Score Interpretation**
- Statistical distributions help understand result quality
- Percentile analysis identifies top-performing results

## Educational Outcomes

After completing experiments with this tool, you should be able to:

1. **Interpret BM25 scores** correctly (negative values, ordering)
2. **Optimize column weights** for specific use cases
3. **Analyze score distributions** to understand search quality
4. **Predict ranking behavior** based on query and document characteristics
5. **Design effective search strategies** using BM25 insights

## Technical Implementation

- **Database**: SQLite with FTS5 virtual tables
- **Tokenization**: Unicode61 tokenizer with Porter stemming
- **Architecture**: Clean separation of commands, handlers, and models
- **Visualization**: ASCII charts using asciigraph library
- **Error Handling**: Type-safe error system with educational messages

## Next Steps

This tool provides foundation for:
- Advanced ranking algorithms
- Machine learning-based relevance tuning  
- Multi-language search optimization
- Real-time search system design

## Troubleshooting

### FTS5 Not Available
```bash
# Verify FTS5 support
sqlite3 :memory: "PRAGMA compile_options;" | grep FTS5

# Ensure build tags
go run -tags "fts5" . corpus generate --size 10 --database test.db
```

### Empty Search Results
```bash
# Check corpus content
go run -tags "fts5" . corpus stats --database test.db

# Try broader queries
go run -tags "fts5" . search query --query "data" --database test.db
```

### Score Analysis Issues
```bash
# Use verbose mode for debugging
go run -tags "fts5" . search explain --query "test" --verbose --database test.db
```

---

**Phase 2 Complete**: You now have hands-on experience with BM25 fundamentals and are ready for advanced contextual memory system development.