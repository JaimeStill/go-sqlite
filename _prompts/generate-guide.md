# Generate Phase Learning Guide

## Prompt Overview

Generate a comprehensive GUIDE.md for the specified phase directory that serves as an instructional walkthrough of the learning objectives. The guide should provide educational value for future learners working through the SQLite FTS5 & BM25 learning project.

## Usage

Execute this prompt with a phase directory argument:

```
Execute _prompts/generate-guide.md for src/##-[phase-name]/
```

## Requirements

### Target Audience

- Developers learning SQLite FTS5 and BM25 concepts
- Students working through the go-sqlite learning project
- Engineers building contextual memory systems

### Content Structure

The generated GUIDE.md must include:

#### 1. Learning Objectives Overview

- **Thorough description** of all learning objectives for the phase
- **Clear articulation** of what skills and knowledge will be gained
- **Contextual placement** within the broader learning project roadmap
- **Prerequisites** and dependencies from previous phases

#### 2. Conceptual Foundations

- **Analogical illustrations** of each core concept associated with the learning objectives
- **Real-world examples** that make abstract concepts concrete
- **Progressive complexity** from basic to advanced concepts
- **Common misconceptions** and clarifications

#### 3. Project Infrastructure

- **Explanation of core infrastructure** directly associated with learning objectives
- **Architectural patterns** and why they were chosen
- **Code organization** and its educational benefits
- **Key files and their roles** in demonstrating concepts

#### 4. Interactive Learning

- **Command execution walkthroughs** with step-by-step instructions
- **Expected outputs** with explanations of what they demonstrate
- **Relationship between commands and learning objectives**
- **Troubleshooting guidance** for common issues

#### 5. Hands-On Exercises

- **Guided experiments** that reinforce learning objectives
- **Progressive challenges** from basic to advanced
- **Self-assessment questions** to verify understanding
- **Extension activities** for deeper exploration

### Writing Guidelines

#### Educational Focus

- **Learning-first approach**: Every section should clearly advance understanding
- **Conceptual clarity**: Complex topics broken into digestible pieces  
- **Practical application**: Theory connected to hands-on experience
- **Knowledge retention**: Repetition and reinforcement of key concepts

#### Technical Accuracy

- **Precise explanations** of SQLite FTS5 and BM25 behavior
- **Accurate command examples** that work as documented
- **Correct interpretation** of outputs and error messages
- **Up-to-date references** to tools and techniques

#### Accessibility

- **Progressive disclosure**: Simple concepts before complex ones
- **Clear language**: Avoid unnecessary jargon, explain when needed
- **Visual organization**: Use headers, lists, and code blocks effectively
- **Multiple learning styles**: Text, examples, exercises, and experiments

## Instructions for Claude

1. **Use the instructor-agent for optimal educational content creation**:

   ```
   @instructor-agent Generate comprehensive GUIDE.md for [phase-directory]
   ```

2. **The instructor-agent will**:
   - Analyze the specified phase directory thoroughly
   - Read README.md to understand the phase objectives
   - Examine the codebase to understand the architecture and key concepts
   - Review command structures and available functionality
   - Identify the specific learning objectives from ROADMAP.md

3. **Generate comprehensive GUIDE.md** following the structure above:
   - Start with a compelling introduction that motivates the learning
   - Provide thorough explanations of learning objectives with concrete examples
   - Include analogical illustrations that make complex concepts accessible
   - Document the project infrastructure and its educational design
   - Create detailed command walkthroughs with expected outputs
   - Design hands-on exercises that reinforce the learning objectives

4. **Ensure educational quality**:
   - Validate that all commands work as documented
   - Verify that explanations accurately reflect the code behavior
   - Confirm that learning objectives are thoroughly addressed
   - Include practical tips and best practices

5. **Structure for future reference**:
   - Create clear sections that can be referenced independently
   - Include a table of contents for navigation
   - Provide summary sections for quick review
   - Link to relevant external resources where appropriate

## Expected Deliverable

A comprehensive GUIDE.md file in the specified phase directory that:

- Serves as a complete instructional walkthrough for the phase
- Demonstrates deep understanding of the learning objectives
- Provides practical, executable examples
- Includes educational explanations of infrastructure and concepts
- Offers hands-on exercises for skill development
- Functions as a standalone learning resource

The guide should enable a motivated learner to work through the phase independently while gaining deep understanding of the SQLite FTS5 and BM25 concepts being taught.
