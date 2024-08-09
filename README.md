
<h1 align="center">dsh</h1>

<p align="center">
  Set of database tools
</p>

```go
go install github.com/yznts/dsh/cmd/...
```

Main goal of the project is to provide a set of multiplatform tools
to work with databases in a unified way,
avoiding differences in UX between clients like `psql`, `sqlite3`, `mysql`, etc.

It tries to stick with the UNIX-like naming and approach,
where each tool does one thing and does it well.
List database tables, or table columns? Just use `dls`.
Get table contents? Use `dcat`.
Or, if you need just to execute an SQL query, `dsql` is here for you.
Want to get the output in JSON, JSONL, CSV?
No problem, just specify an according flag, like `-json` or `-csv`.

![example](.github/assets/example.png)

Now, utility set includes:
- `dls`   - lists database tables or table columns
- `dsql`  - executes SQL queries
- `dcat`  - outputs table data (in the not-so-dumb way)
- `dps`   - lists database processes (if supported by the database)
- `dkill` - kills database processes (if supported by the database)

May be used with:
- `sqlite`
- `postgresql`

And supports this output formats:
- `json` (partial support)
- `jsonl`
- `csv`
- `gloss` (default terminal output)

## Installation

You have multiple ways to install/use this utility set:
- Install in Go-way
- Build by yourself
- Download binaries
- Spin-up a Docker container

### Install in Go-way

This is the easiest way to install,
but you need to have Go installed on your machine,
including `GOBIN` in your `PATH`.

```bash
go install github.com/yznts/dsh/cmd/...
```

### Build by yourself

This way still requires Go to be installed on your machine,
but it's up to you to decide where to put the binaries.

```bash
mkdir -p /tmp/dsh && cd /tmp/dsh
git clone git@github.com:yznts/dsh.git .
make build

# You'll find binaries in the `bin` directory.
# Feel free to move them to the desired location, e.g. /usr/local/bin.
```

### Download binaries

Also you have an option to download the latest binaries from the [Releases](https://github.com/yznts/dsh/releases) page.
Please note, that darwin(macos) binaries are not signed!
If you know a simple way to handle this issue, please open issue or PR.

### Spin-up a Docker container

Docker way doesn't require Go to be installed on your machine
and it allows you to use the tooling in isolated way,
without polluting your system.

```bash
docker run --rm -it ghcr.io/yznts/dsh:latest
```
