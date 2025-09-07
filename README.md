# Scruffy

> [!CAUTION]
> This repository contains unstable code, if a `rm -rf /` occurs, it's not my problem, use with caution.

---

General utility, but at the moment, this is just a go version of [mateusjdev/rename-files-to-hash](https://github.com/mateusjdev/rename-files-to-hash).

## 🛠️ Building

### 📦 Prerequisites

Make sure you have the following installed:

- [Go](https://golang.org/doc/install) (version ≥ 1.21)
- For the build process [Taskfile](https://taskfile.dev/#/installation) and [Git](https://git-scm.com/downloads) are highly recommended.

### 🏗️ Building

1. Clone the repository to your local machine:

```shell
git clone https://github.com/mateusjdev/scruffy
cd scruffy
```

2. Build the project:

```shell
go build -o ./build/scruffy
# or using Taskfile (recommended)
task build
```
