# BM25 Quick Reference & Best Practices

## Algorithm Overview

### BM25 Formula

```
score(D,Q) = Σ IDF(qi) × [f(qi,D) × (k1 + 1)] / [f(qi,D) + k1 × (1 - b + b × |D|/avgdl)]
```

**Components:**

- **IDF(qi)**: Inverse Document Frequency of query term i
- **f(qi,D)**: Term frequency of qi in document D  
- **|D|**: Document length in tokens
- **avgdl**: Average document length in collection
- **k1**: Term frequency saturation parameter (default: 1.2)
- **b**: Document length normalization parameter (default: 0.75)

### SQLite FTS5 Implementation

```sql
-- SQLite returns NEGATIVE BM25 scores (lower = better)
SELECT *, bm25(table_name) as score 
FROM table_name 
WHERE table_name MATCH 'query'
ORDER BY bm25(table_name); -- Ascending = best first

-- Using rank column (equivalent to bm25)
SELECT * FROM table_name 
WHERE table_name MATCH 'query'
ORDER BY rank;
```

## Parameter Tuning

### k1 Parameter (Term Frequency Saturation)

**Range**: 0.0 - 3.0 (typical: 1.2 - 2.0)

```sql
-- SQLite FTS5: k1 is hardcoded to 1.2, cannot be changed
-- For reference, other systems allow:

-- Low k1 (0.5): Quick saturation, diminishing returns
-- Standard k1 (1.2): Balanced term frequency impact  
-- High k1 (2.0): Slower saturation, rewards repetition
```

**Effects:**

- **k1 → 0**: Approaches binary matching (term present/absent)
- **k1 = 1.2**: At average document length, tf=1.2 gives 50% of max score
- **k1 → ∞**: Approaches pure term frequency weighting

### b Parameter (Length Normalization)

**Range**: 0.0 - 1.0 (typical: 0.75)

```sql
-- SQLite FTS5: b is hardcoded to 0.75, cannot be changed
-- For reference, other systems allow:

-- b = 0: No length normalization
-- b = 0.75: Balanced length consideration
-- b = 1: Full length normalization
```

**Effects:**

- **b = 0**: Document length doesn't affect scoring
- **b = 0.75**: Moderate penalty for longer documents
- **b = 1**: Strong penalty for longer documents

### Column Weighting (SQLite FTS5)

```sql
-- Custom column weights (title=2x, content=1x, tags=3x)
SELECT *, bm25(docs, 2.0, 1.0, 3.0) as weighted_score
FROM docs 
WHERE docs MATCH 'query'
ORDER BY weighted_score;

-- Equal weighting (default behavior)
SELECT *, bm25(docs) as score
FROM docs 
WHERE docs MATCH 'query'
ORDER BY score;
```

## Score Interpretation

### SQLite FTS5 Scoring

```sql
-- SQLite inverts BM25 scores (multiplies by -1)
-- Lower (more negative) scores = better matches

-- Example scores:
-- Perfect match: -5.2
-- Good match: -3.1  
-- Fair match: -1.8
-- Poor match: -0.5
```

### Practical Score Ranges

- **Exact phrase matches**: -8.0 to -15.0
- **Multiple term matches**: -3.0 to -8.0  
- **Single term matches**: -0.5 to -3.0
- **Weak matches**: -0.1 to -0.5

### Score Analysis Queries

```sql
-- Score distribution analysis
SELECT 
    ROUND(bm25(docs), 1) as score_range,
    COUNT(*) as count
FROM docs 
WHERE docs MATCH 'query'
GROUP BY ROUND(bm25(docs), 1)
ORDER BY score_range;

-- Top scoring documents with context
SELECT 
    title,
    bm25(docs) as score,
    snippet(docs, 1, '[', ']', '...', 5) as context
FROM docs 
WHERE docs MATCH 'query'
ORDER BY rank
LIMIT 10;
```

## Best Practices

### Query Construction

```sql
-- Single term
SELECT * FROM docs WHERE docs MATCH 'sqlite';

-- Multiple terms (AND implied)
SELECT * FROM docs WHERE docs MATCH 'sqlite database';

-- Explicit boolean
SELECT * FROM docs WHERE docs MATCH 'sqlite AND database';

-- Phrase matching (higher scores)
SELECT * FROM docs WHERE docs MATCH '"sqlite database"';

-- Column targeting
SELECT * FROM docs WHERE docs MATCH 'title:sqlite';
```

### Score Optimization Strategies

#### 1. Document Structure

- **Shorter documents** generally score higher for term density
- **Front-load important terms** in title/summary fields
- **Consistent terminology** improves matching precision

#### 2. Query Enhancement

```sql
-- Boost exact phrases
SELECT * FROM docs 
WHERE docs MATCH '"exact phrase" OR (exact AND phrase)'
ORDER BY rank;

-- Field boosting through weights
SELECT *, bm25(docs, 3.0, 1.0) as score -- title=3x, content=1x
FROM docs WHERE docs MATCH 'query'
ORDER BY score;
```

#### 3. Collection Optimization

- **Remove stop words** during indexing to improve IDF calculations
- **Stemming** consolidates term variations
- **Document segmentation** can improve relevance for long texts

### Performance Considerations

#### Index Size Impact

```sql
-- Check FTS index size
SELECT 
    name,
    COUNT(*) as segments,
    SUM(pgno) as pages
FROM dbstat 
WHERE name LIKE '%_fts%'
GROUP BY name;
```

#### Query Performance

- **BM25 calculation** adds minimal overhead to FTS5 queries
- **Column weighting** has negligible performance impact
- **Complex boolean queries** can be slower than simple term matching

## Advanced Techniques

### Custom Scoring Functions

```sql
-- Simulate custom ranking with computed scores
SELECT *,
    bm25(docs) as base_score,
    -- Custom boost for recent documents
    (bm25(docs) + (julianday('now') - julianday(created)) * -0.1) as time_boosted_score
FROM docs 
WHERE docs MATCH 'query'
ORDER BY time_boosted_score;
```

### Score Normalization

```sql
-- Normalize scores to 0-1 range within result set
WITH scored_results AS (
    SELECT *,
        bm25(docs) as raw_score
    FROM docs 
    WHERE docs MATCH 'query'
),
score_bounds AS (
    SELECT 
        MIN(raw_score) as min_score,
        MAX(raw_score) as max_score
    FROM scored_results
)
SELECT 
    title,
    raw_score,
    (raw_score - min_score) / (max_score - min_score) as normalized_score
FROM scored_results, score_bounds
ORDER BY raw_score;
```

### A/B Testing Framework

```sql
-- Compare different weighting schemes
WITH weight_test AS (
    SELECT 
        title,
        bm25(docs) as default_score,
        bm25(docs, 2.0, 1.0) as title_boosted,
        bm25(docs, 1.0, 3.0) as content_boosted
    FROM docs 
    WHERE docs MATCH 'test query'
)
SELECT 
    title,
    default_score,
    title_boosted,
    content_boosted,
    -- Rank changes
    ROW_NUMBER() OVER (ORDER BY default_score) as default_rank,
    ROW_NUMBER() OVER (ORDER BY title_boosted) as title_rank,
    ROW_NUMBER() OVER (ORDER BY content_boosted) as content_rank
FROM weight_test
ORDER BY default_score
LIMIT 20;
```

## Limitations and Alternatives

### BM25 Limitations

- **Lexical matching only**: No semantic understanding
- **Fixed parameters**: Cannot tune k1/b in SQLite FTS5
- **Document independence**: Doesn't consider document relationships
- **No learning**: Static algorithm, no adaptation to user behavior

### When to Consider Alternatives

- **Semantic search needed**: Use embedding-based retrieval
- **Parameter tuning required**: Consider external search engines
- **Complex document relationships**: Graph-based or neural approaches
- **User personalization**: Machine learning ranking models

### Hybrid Approaches

```sql
-- Combine BM25 with additional signals
SELECT 
    title,
    bm25(docs) as relevance_score,
    view_count,
    recency_score,
    -- Weighted combination
    (bm25(docs) * 0.7 + 
     LOG(view_count + 1) * 0.2 + 
     recency_score * 0.1) as final_score
FROM docs 
WHERE docs MATCH 'query'
ORDER BY final_score;
```

## Debugging and Analysis

### Query Explanation

```sql
-- Understand query execution
EXPLAIN QUERY PLAN 
SELECT * FROM docs 
WHERE docs MATCH 'complex AND query'
ORDER BY rank;
```

### Score Debugging

```sql
-- Detailed score analysis
SELECT 
    title,
    bm25(docs) as score,
    -- Term frequency indicators
    LENGTH(content) - LENGTH(REPLACE(LOWER(content), 'searchterm', '')) as term_freq,
    LENGTH(content) as doc_length,
    snippet(docs, -1, '[', ']', '...', 3) as matched_terms
FROM docs 
WHERE docs MATCH 'searchterm'
ORDER BY score;
```

### Performance Monitoring

```sql
-- Query performance analysis
.timer on
SELECT COUNT(*) FROM docs WHERE docs MATCH 'performance test';
.timer off

-- Index statistics
PRAGMA table_info(docs_fts);
```

## Integration Patterns

### Result Processing

```sql
-- Pagination with consistent scoring
SELECT 
    title,
    bm25(docs) as score,
    ROW_NUMBER() OVER (ORDER BY rank) as position
FROM docs 
WHERE docs MATCH 'query'
ORDER BY rank
LIMIT 20 OFFSET 40;
```

### Faceted Search

```sql
-- Category-aware scoring
SELECT 
    category,
    COUNT(*) as count,
    AVG(bm25(docs)) as avg_score,
    MIN(bm25(docs)) as best_score
FROM docs 
WHERE docs MATCH 'query'
GROUP BY category
ORDER BY best_score;
```
