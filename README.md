<div align="center">
  
# ghfetch

:octocat: ghfetch is a CLI tool to fetch GitHub user information and show like neofetch
  
<a href="https://opensource.org/licenses/MIT">
<img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="MIT License badge">
</a>
<a href="https://pkg.go.dev/github.com/orangekame3/stree">
<img src="https://github.com/orangekame3/ghfetch/actions/workflows/release.yml/badge.svg" alt="Release workflow status badge">
</a>
<a href="https://github.com/orangekame3/ghfetch/actions/workflows/tagpr.yml">
<img src="https://github.com/orangekame3/ghfetch/actions/workflows/tagpr.yml/badge.svg" alt="Tag PR workflow status badge">
</a>
</div>

## Demo

<p align="center">
<img src="img/demo.gif" alt="Demonstration of ghfetch" height="auto" width="auto"/>
</p>

## Install

### Go

```shell
go install github.com/orangekame3/ghfetch@latest
```

### Homebrew

```shell
brew install orangekame3/tap/ghfetch
```

## Usage

```shell
❯ ghfetch -h
Fetch GitHub user's profile

Usage:
  ghfetch [flags]

Flags:
  -c, --color string   Highlight color for text (default "blue")
  -h, --help           help for ghfetch
  -t, --toggle         Help message for toggle
  -u, --user string    GitHub username
  -v, --version        version for ghfetch
```

## To Do

- [ ] Add GitHub stats
