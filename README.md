# github-blob-sender [![Build Status](https://travis-ci.org/ArthurHlt/github-blob-sender.svg?branch=master)](https://travis-ci.org/ArthurHlt/github-blob-sender)

Send a file into github blob in command line (written in golang).
The main goal is to store and restore file (large or medium) in github.
Those files will not be visible on your github repo (instead of github lfs) and will can be only downloaded if you know the link.

github-blob-sender can upload file from your computer into a github repo and keep trace of these upload into a file
(A `.github-blob-sender` file) and could restore uploaded file with a checksum verification to be reliable.

## Installation

### On *nix system

You can install this via the command-line with either `curl` or `wget`.

#### via curl

```bash
$ sh -c "$(curl -fsSL https://raw.github.com/ArthurHlt/github-blob-sender/master/bin/install.sh)"
```

#### via wget

```bash
$ sh -c "$(wget https://raw.github.com/ArthurHlt/github-blob-sender/master/bin/install.sh -O -)"
```

### On windows

You can install it by downloading the `.exe` corresponding to your cpu from releases page: https://github.com/ArthurHlt/github-blob-sender/releases .
Alternatively, if you have terminal interpreting shell you can also use command line script above, it will download file in your current working dir.

### From go command line

Simply run in terminal:

```bash
$ go get github.com/ArthurHlt/github-blob-sender
```

If compile failed you can use godep to restore dependencies:


## Usage

### Global command

```
NAME:
   github-blob-sender - Store and restore file from github blob api

USAGE:
   github-blob-sender [global options] command [command options] [arguments...]

COMMANDS:
   upload, u		Upload file to github blob
   cat, c		Cat file from github blob
   download, d		Download file from github blob
   download-all, a	Download all registered files in folder from github blob
   list, l		List registered files
   help, h		Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h		show help
   --version, -v	print the version
```

### `upload` command

```
NAME:
   github-blob-sender upload - Upload file to github blob

USAGE:
   github-blob-sender upload [command options] [file1 file2 ...]

OPTIONS:
   --github-token, --gt 	Set your github token (optional if already set or GITHUB_TOKEN env var set)
   --owner, -o 			Which org or user you own to send file
   --repo, -r 			Which repo you own to send file
```

### `cat` command

```
NAME:
   github-blob-sender cat - Cat file from github blob

USAGE:
   github-blob-sender cat [command options] [file-name] (Note: file-name can be listed with list command)

OPTIONS:
   --github-token, --gt 	Set your github token (optional if already set or GITHUB_TOKEN env var set)
   --owner, -o 			Which org or user you own to send file (optional)
   --repo, -r 			Which repo you own to send file (optional)
```

**TIP**: `owner` and `repo` flags are optional if only one file with this name is registered in `.github-blob-sender` file

### `download` command

```
NAME:
   github-blob-sender download - Download file from github blob

USAGE:
   github-blob-sender download [command options] [file-name] (Note: file-name can be listed with list command)

OPTIONS:
   --github-token, --gt 	Set your github token (optional if already set or GITHUB_TOKEN env var set)
   --owner, -o 			Which org or user you own to send file (optional)
   --repo, -r 			Which repo you own to send file (optional)
   --output 			Set where to write the downloaded file
```

**TIP**: `owner` and `repo` flags are optional if only one file with this name is registered in `.github-blob-sender` file

### `download-all` command

```
NAME:
   github-blob-sender download-all - Download all registered files in folder from github blob

USAGE:
   github-blob-sender download-all [command options] [folder/to/put/downloaded/file]

OPTIONS:
   --create-folders, --create-folder, -c	Create folders if not exist
```
### `list` command

```
NAME:
   github-blob-sender list - List registered files

USAGE:
   github-blob-sender list [command options] [arguments...]

OPTIONS:
   --show-github-sha1, -g	Show github checksum (in sha1)
   --show-registered-sha1, -r	Show registered checksum (in sha1)
   --show-link, -l		Show github link
```

**Example output**:

```
+-----------+-----------+--------------------+
|   NAME    |   OWNER   |        REPO        |
+-----------+-----------+--------------------+
| FILE1     | ArthurHlt | github-blob-sender |
| FILE2     | ArthurHlt | github-blob-sender |
+-----------+-----------+--------------------+
```

