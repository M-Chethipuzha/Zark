# Zark - Next-Generation Version Control System

Zark is a next-generation version control system designed to address the complexities and performance limitations of existing solutions, particularly for new professionals and large-scale codebases. Zark aims to provide a more intuitive, powerful, and performant alternative to traditional Git, fostering better collaboration and reducing friction in software development workflows.

## Features

- **Repository Initialization:** Create a new Zark repository with `zark start`
- **Staging Area (Index):** Add files to a staging area before committing
- **Committing:** Save snapshots of your staged files with `zark save`
- **History:** View the commit history of a branch with `zark history`
- **Status:** Check the status of your working directory and staging area
- **Branching:** Create and list branches
- **Checkout:** Switch between branches or commits

## Getting Started

### Prerequisites

- Go version 1.21 or later

### Build

To build the zark executable, clone the repository and run the following command from the project's root directory:

```bash
go build -o zark ./cmd/zark
```

This will create an executable named `zark` (or `zark.exe` on Windows) in your current directory.

## Usage

Here are the basic commands currently implemented in Zark:

### Repository Management

**Initialize a repository**

```bash
./zark start
```

### Working with Files

**Check the status**

```bash
./zark status
```

**Add files to staging area**

```bash
./zark add <filename>
```

**Save changes (commit)**

```bash
./zark save
# You will be prompted to enter a commit message interactively
```

### Viewing History

**View commit history**

```bash
./zark history
```

### Branch Management

**Create a new branch**

```bash
./zark branch <branch-name>
```

**List all branches**

```bash
./zark branch
```

**Switch to a branch**

```bash
./zark checkout <branch-name>
```

## Example Workflow

Here's a typical workflow using the currently implemented commands:

```bash
# Initialize a new repository
./zark start

# Create a file and check status
echo "Hello, Zark!" > hello.txt
./zark status

# Add and save changes
./zark add hello.txt
./zark save
# Enter commit message when prompted

# View history
./zark history

# Create a new branch
./zark branch new-feature

# Switch to the new branch
./zark checkout new-feature

# Make changes and save
echo "New feature code" > feature.txt
./zark add feature.txt
./zark save

# Switch back to main
./zark checkout main
```

## Key Improvements Over Traditional Git

- **Simplified Commands:** Intuitive command names like `start`, `save`, and `history`
- **Interactive Prompts:** The `save` command provides interactive prompts for commit messages

## Future Implementation

The following features are planned for future releases:

### Enhanced Beginner-Friendly Features

- **Repository Cloning:** `zark get <repository-url>` (instead of `git clone`)
- **Smart Syncing:** `zark sync` with intelligent push/pull detection
- **Guided Workflows:** Interactive prompts for complex operations like `zark branch create <name>` and `zark merge <branch>`
- **Context-Sensitive Help:** `zark help <command>` system
- **Intelligent Error Messages:** Clear, actionable advice when commands fail
- **Visual Diff & Merge Tool Integration:** Built-in or recommended visual tools
- **Smart Autocompletion:** Predictive text for commands and file paths

## Running Tests

To run all the unit tests for the project, use the following command from the root directory:

```bash
go test -v ./...
```

## Architecture

Zark is built with:

- **Core Engine:** Written in Go for high performance and memory safety
- **Content-Addressable Storage:** Similar to Git but with planned optimizations
- **Extensible Design:** Plugin architecture for future enhancements

---

_Zark: Making version control intuitive, powerful, and performant for everyone._
