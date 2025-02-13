# Introduction to Git and GitHub

## What is Git?
Git is a **distributed version control system** that helps developers track and manage changes to their code. By using Git:
- You can keep a detailed history of every modification.
- Collaborate efficiently with others on the same project without overwriting each other’s work.
- Roll back to previous versions if something goes wrong.

## What is GitHub?
GitHub is a **web-based hosting service** for Git repositories. It provides:
- A place to store your repositories online.
- Collaboration tools (issue tracking, pull requests, code reviews).
- Hosting and deployment options for certain types of web projects.

## Why Use Git and GitHub?
- **Version Control**: Keep track of every change with clear commit messages.
- **Collaboration**: Work with team members on the same codebase without conflicts.
- **Backup and Sharing**: Store your code safely and allow others to contribute or review.
- **Portfolio**: Showcase your work to potential employers or the open-source community.

---

## Getting Started with Git and GitHub

### 1. Install Git
1. Download and install Git from the [official website](https://git-scm.com/downloads).
2. After installation, open a terminal (or command prompt).
3. Configure Git with your information:

    git config --global user.name "Your Name"
    git config --global user.email "you@example.com"

   This sets your identity for commits on your local machine.

### 2. Create a Local Repository
1. Open your terminal and navigate to the folder where you want your project to be.
2. Initialize the folder as a Git repository:

    git init

   This command creates a hidden `.git` folder inside your project, where Git stores its data.

3. Add files to your repository and commit:

    git add .
    git commit -m "Initial commit"

   `git add .` stages all changes, and `git commit -m "message"` creates a snapshot of your code.

### 3. Create a GitHub Account and Repository
1. Go to [GitHub.com](https://github.com/) and create an account if you don't have one.
2. Once logged in, click the “+” icon in the top-right corner and select **New repository**.
3. Provide a name for your repository and click **Create repository**.

### 4. Connect Your Local Repository to GitHub
1. Copy the repository URL from your newly created GitHub repository (e.g., `https://github.com/username/repository-name.git`).
2. In your terminal, set the remote URL:

    git remote add origin https://github.com/username/repository-name.git

3. Push your local commits to GitHub:

    git push -u origin master

   After the first push, you can simply use `git push` for subsequent commits.

---

## Basic Git Commands Cheat Sheet
- **git status**  
  Shows the status of your working tree (which files are staged, unstaged, or untracked).

- **git add <file>**  
  Stages a specific file. Use `git add .` to stage all changes.

- **git commit -m "message"**  
  Commits staged changes with a message describing the changes.

- **git log**  
  Shows the commit history.

- **git branch**  
  Lists all branches; add a name to create a new branch (for example, `git branch feature-xyz`).

- **git checkout <branch>**  
  Switches to another branch.

- **git merge <branch>**  
  Merges changes from one branch into your current branch.

- **git pull**  
  Updates your local repository with changes from the remote repository.

- **git push**  
  Uploads local commits to the remote repository on GitHub.

---

## Working with Others on GitHub
- **Fork** a repository you don’t own to create a personal copy on GitHub.
- **Clone** that repository locally using `git clone <URL>` to start working on it.
- **Push** changes back to your fork, and open a **Pull Request** to propose your changes to the original repository’s maintainers.

---

## Summary
Git and GitHub are essential tools for modern software development. Git allows you to record and track every change in your project, while GitHub lets you share your work and collaborate with others easily. By following the steps above, you can set up your environment, create your first project, and begin taking advantage of version control.
