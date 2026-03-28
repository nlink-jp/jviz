# jviz

Visualize JSON data in your browser from the command line.

Pipe a JSON array to `jviz` and it opens a local web page with interactive charts. When the data updates, the chart updates live.

## Install

```sh
go install github.com/nlink-jp/jviz@latest
```

Or download a pre-built binary from the [releases page](https://github.com/nlink-jp/jviz/releases).

## Usage

```sh
# Pipe JSON from stdin
cat data.json | jviz

# Stream: every new JSON array refreshes the chart
while true; do generate-stats.sh; sleep 5; done | jviz

# Watch a file for changes
jviz --watch data.json

# Change port (default: 8765)
jviz --port 9000 < data.json

# Do not auto-open browser
jviz --no-open < data.json
```

### Input format

`jviz` expects a JSON array of objects:

```json
[
  {"label": "Jan", "sales": 120},
  {"label": "Feb", "sales": 95},
  {"label": "Mar", "sales": 140}
]
```

### Chart types

| Type  | Description               |
|-------|---------------------------|
| Bar   | Vertical bar chart        |
| Line  | Line chart with fill      |
| Pie   | Pie / donut chart         |
| Table | Raw data table            |

Use the **X / Label** and **Y / Value** selectors to choose which columns to plot.

## Works well with

- [jstats](https://github.com/nlink-jp/jstats) — aggregate JSON by field
- [csv-to-json](https://github.com/nlink-jp/csv-to-json) — convert CSV to JSON
- [json-filter](https://github.com/nlink-jp/json-filter) — filter JSON arrays

## Build

```sh
make build        # current platform → dist/jviz
make build-all    # all 5 platforms  → dist/
make test
```

## License

MIT — see [LICENSE](LICENSE).

Chart.js is bundled under its own MIT license (© 2023 Chart.js Contributors).
