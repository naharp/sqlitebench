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

The benchmark was run on KVM Ubuntu 18.04 with 4 cores of Xeon CPU E7-8870 v4 @ 2.10GHz and 8 GB of RAM.

The server is started with a configurable number of "connections" to the SQLite database, here called _poolsize_. Once the server is running [hey](https://github.com/rakyll/hey) is used to run the HTTP load test. See `runbenchmark.go` for details.

## Result

```
package         rps @ C1        rps @ C4        rps @ C8        rps @ C50       rps @ C100
crawshaw        29626.7135      30867.1062      27580.4381      28274.1264      26701.9440
mattn           28835.1667      27950.9703      28475.7259      28521.1984      28147.3413
modernc         38659.9235      39049.2933      38437.6926      35588.3684      36337.5995
tailscale       28172.8906      27527.9857      27229.4827      27325.8529      27974.9314
zombiezen       32595.6653      35346.5005      35777.0249      36909.9302      40441.6229
```
