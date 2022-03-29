<h1 align="center">
<br />
SQLite/HTTP Server Performance Benchmark
<br /><br />
</h1>

Forked from [philippta/sqlitebench](github.com/philippta/sqlitebench)

This benchmark provides a overview of the different SQLite driver performances available in Go. For benchmarking a simple HTTP server is used to perform random read queries on the database.

For benchmarking the [hey](https://github.com/rakyll/hey) load generator is used to call the HTTP server (with 50 concurrent requests).

### Driver Overview

Package | Uses CGo | Is driver for `database/sql`
:------ |:--------:| :-----:
crawshaw.io/sqlite |   Yes    | No
github.com/mattn/go-sqlite3 |   Yes    | Yes
modernc.org/sqlite |    No    | Yes
github.com/tailscale/tailscale |   Yes    | Yes
zombiezen.com/go/sqlite |    No    | No

## Implementation

The implementation consists of a simple HTTP server that runs a single select query on the SQLite database for each request.
```sql
SELECT * FROM foo WHERE rowid = ? -- where ? is a random number between 1 and 10000
```

The SQLite database has the following schema and contains 100000 rows with random values.
```sql
CREATE TABLE foo (id integer, value integer);
```

See any of the subfolder/main.go files for more details.

## Benchmark

The benchmark was run on a MacBook Pro 2020 with a 2.3 GHz Quad-Core Intel Core i7 and 32 GB of RAM.

The server is started with a configurable number of "connections" to the SQLite database, here called _poolsize_. Once the server is running [hey](https://github.com/rakyll/hey) is used to run the HTTP load test. See `runbenchmark.go` for details.

All reports of hey can be found in the `result_*.txt` files.

## Result

```

```

## Observations

- The performance results between the packages do not differ that much.
- Limiting the number of connections/poolsize to the SQLite database  to roughly the number of CPU cores on the machine gives best performance results.
