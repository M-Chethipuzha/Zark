# Zark - A Beginner-Friendly Version Control System

Zark is a simple, user-friendly version control system built in Go that makes it easy to track changes in your code. If you've ever been frustrated by Git's complexity or just want to learn how version control works, Zark is perfect for you!

Think of Zark as a time machine for your code - it helps you save snapshots of your work, try new features without fear, and collaborate with others safely.

## Why Zark?

- **Easy to Learn:** Simple commands that make sense (like `start` instead of `init`)
- **Interactive Help:** Guided prompts walk you through complex operations
- **Powerful Features:** Advanced tools that grow with your skills
- **Git-Inspired:** Based on proven concepts, so skills transfer to other tools

## What Can Zark Do?

### Essential Features (Perfect for Beginners)

- **Save Your Work:** Take snapshots of your code at any point
- **Track Changes:** See what files you've modified
- **Work on Features:** Create separate branches to try new ideas safely
- **Go Back in Time:** View your project's history and revert mistakes
- **Stay Organized:** Keep your repository clean and optimized

### Advanced Features (For Growing Developers)

- **Search Your History:** Find specific commits by author or message
- **Handle Large Files:** Track big files (like images or videos) efficiently
- **Security First:** Scan for secrets and sign your commits
- **Team Collaboration:** Tools to work better with others

## Getting Started

### What You Need

- Go version 1.21 or later installed on your computer

### Installation

1. Download or clone the Zark repository
2. Open your terminal and navigate to the project folder
3. Build Zark by running:

```bash
go build -o zark ./cmd/zark
```

You'll now have a `zark` executable ready to use!

## Your First Zark Repository

Let's create your first project and learn the basics:

### Step 1: Start a New Project

```bash
./zark start
```

This creates a new Zark repository in your current folder. Think of it as telling Zark "start watching this folder for changes."

### Step 2: Create Some Files

```bash
echo "Hello, World!" > hello.txt
echo "This is my first Zark project" > README.md
```

### Step 3: Check What's Changed

```bash
./zark status
```

This shows you what files are new or modified. It's like asking "what's different since my last save?"

### Step 4: Stage Your Changes

```bash
./zark add hello.txt
./zark add README.md
```

This tells Zark "I want to include these files in my next snapshot."

### Step 5: Save Your Work

```bash
./zark save -m "My first commit - added hello.txt and README.md"
```

This creates a permanent snapshot of your work with a description.

### Step 6: See Your History

```bash
./zark history
```

This shows all the snapshots you've saved. Each snapshot is called a "commit."

## Essential Commands

### Working with Files

```bash
# Check what's changed
./zark status

# Add files to your next commit
./zark add filename.txt
./zark add .  # Add all changed files

# Save your work with a message
./zark save -m "Describe what you changed"

# Save and sign your commit (more secure)
./zark save -m "Important change" -s
```

### Viewing History

```bash
# See all your commits
./zark history

# Search for commits by a specific author
./zark search --author "your-name"

# Search for commits with specific words
./zark search --message "bug fix"
```

### Working with Branches

Branches let you work on different features without affecting your main code:

```bash
# Create a new branch (with helpful prompts)
./zark branch create

# See all your branches
./zark branch

# Switch to a different branch
./zark checkout branch-name

# Go back to your main branch
./zark checkout main
```

## Intermediate Features

### Keep Your Repository Clean

```bash
# Optimize your repository (makes it smaller and faster)
./zark gc
```

### Handle Large Files

```bash
# Track large files efficiently (like images, videos, zip files)
./zark lfs track "*.jpg"
./zark lfs track "*.mp4"
./zark lfs track "*.zip"
```

### Work Securely

Zark helps prevent common mistakes:

- **Secret Scanning:** Automatically warns you before committing passwords or API keys
- **Code Signing:** Sign important commits to prove they're from you

## Common Workflows

### Starting a New Feature

```bash
# Create a branch for your new feature
./zark branch create
# (Follow the prompts to name your branch)

# Work on your feature...
echo "new feature code" > feature.txt
./zark add feature.txt
./zark save -m "Add new feature"

# Switch back to main when done
./zark checkout main
```

### Finding Old Work

```bash
# Look for commits you made
./zark search --author "your-name"

# Find when you fixed a bug
./zark search --message "fix"

# See everything you've done
./zark history
```

### Keeping Things Organized

```bash
# Clean up your repository regularly
./zark gc

# Track large files properly
./zark lfs track "*.pdf"
./zark lfs track "*.zip"
```

## Tips for Success

1. **Commit Often:** Save your work frequently with descriptive messages
2. **Use Branches:** Try new ideas in separate branches to keep your main code safe
3. **Write Good Messages:** Describe what you changed and why
4. **Clean Up Regularly:** Run `./zark gc` occasionally to keep things fast
5. **Track Large Files:** Use LFS for files bigger than a few MB

## Example: Building a Simple Website

Let's walk through a real example:

```bash
# Start your project
./zark start

# Create your first files
echo "<h1>My Website</h1>" > index.html
echo "body { font-family: Arial; }" > style.css

# Save your initial version
./zark add index.html style.css
./zark save -m "Initial website with basic HTML and CSS"

# Create a branch for a new feature
./zark branch create
# Name it "add-navigation"

# Add navigation
echo "<nav><a href='#'>Home</a> <a href='#'>About</a></nav>" >> index.html
./zark add index.html
./zark save -m "Add navigation menu"

# Go back to main
./zark checkout main

# See your history
./zark history

# Search for navigation-related commits
./zark search --message "navigation"
```

## Testing Your Installation

Run the test suite to make sure everything works:

```bash
go test -v ./...
```

## What's Coming Next?

Zark is actively being developed! Future features include:

- **Smart Merging:** Better tools for combining different branches
- **Visual Tools:** Graphical interfaces for complex operations
- **Performance Boosts:** Even faster operations on large projects

## Need Help?

- **New to Version Control?** Start with the basic commands above
- **Coming from Git?** Most concepts are similar, just with friendlier names
- **Want to Learn More?** Experiment with branches and different workflows
- **Found a Bug?** Check the test suite and report issues

Remember: version control is like learning to ride a bike - it seems complex at first, but once you get it, you can't imagine coding without it!

---

_Zark: Making version control intuitive, powerful, and performant for everyone._
