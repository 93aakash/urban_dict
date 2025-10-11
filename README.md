# urban_dict

A tiny CLI tool to look up terms from **Urban Dictionary**.  
It stores results in a local **SQLite3** database for faster future lookups.

---

## Features
- Lightweight
- Cross-platform (Linux, Windows, MacOS)  
- Uses SQLite3 for caching

---

## Installation

Building from the source requires [Go](https://go.dev/).

Build command:

```bash
go install -ldflags="-s -w" github.com/93aakash/urban_dict@latest
```

---

## Usage
Run urban_dict -help to see command line options.

```bash
NAME:
   urban_dict - Look up terms using Urban Dictionary

USAGE:
   urban_dict [global options] term

VERSION:
   1.0.0

GLOBAL OPTIONS:
   --delete, -d   (default: false)
   --debug, -x    (default: false)
   --help, -h     show help
   --version, -v  print the version
```

#### Example

```bash
urban_dict yeet
urban_dict -d yeet
```
