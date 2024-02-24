<h1 align="center">sq</h1>
<h3 align="center">Convert and query JSON, JSONL, CSV, and SQLite with ease!</h3>

`sq` is a simple yet powerful command line tool for query and converting from and to SQLite, CSV, JSON, and JSONL files, heavily inspired by [ `jq` ](https://jqlang.github.io/jq/).

## Features

* Convert between SQLite, CSV, JSON, and JSONL files.
* Query CSV, JSON, and JSONL using real SQLite queries.
* Allows for rapid scripting and conversion of data.
* Pipe in data or read from files.

## Installation

Download a prebuilt binary from the [releases page](https://github.com/chand1012/sq/releases) and add it to your PATH, or use `go install` to install the latest version.

```bash
go install github.com/chand1012/sq@latest
```

## Examples

Query some orders from a CSV file. Column names for CSV files are converted to lower case and spaces are replaced with underscores.

```bash
$ sq -r orders.csv 'select country from sq where seller_amount > 20 not null;'
United States of America
United States of America
United States of America
United States of America
United States of America
United States of America
United Kingdom
Canada
United States of America
United States of America
United States of America
United States of America
United States of America
Canada
United States of America
Canada
Canada
Australia
United States of America
...
```

Download and query some JSONL datasets.

```bash
$ curl https://raw.githubusercontent.com/TimeSurgeLabs/llm-finetuning/4e934ce602f34f62f4d803c40cd1e7825d216192/data/fingpt-sentiment-1k.jsonl | sq 'select * from sq where output = "positive";' -f jsonl > positive.jsonl
```

You can even use it with `jq` !

```bash
$ curl https://api.gogopool.com/stakers | jq '.stakers' | sq -t stakers 'SELECT stakerAddr,avaxValidating FROM stakers WHERE avaxValidating > 0;' -f json > stakers.json
```
