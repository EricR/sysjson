# sys.json

Exposes server performance stats as an API.

It works, but it's a work in progress. Things are subject to change without notice, so please use
caution if you decide to use this in production. Assume there are bugs.

## Why?

There are plenty of other tools out there that can give you stats. There's SNMP, dozens of monitoring
solutions and stat collection agents, and more. So why this?

Many of these solutions are complicated to setup and get data from. The goal of sys.json is to
provide a simple API for querying a server about the stats you're interested in. No more, no less.
What you do with that data is up to you.

## Running It

```bash
go get
go build
./sys.json --listen 0.0.0.0:5374
```

## Example Requests

```bash
curl http://localhost:5374/
```

```js
{
  "modules_available": {
    "conntrack": {
      "name": "conntrack",
      "description": "Provides stats on current network connections (requires conntrack-tools)"
    },
    "disk": {
      "name": "disk",
      "description": "Provides stats on each disk"
    },
    "host": {
      "name": "host",
      "description": "Provides basic host info"
    },
    "load": {
      "name": "load",
      "description": "Provides load averages"
    },
    "mem": {
      "name": "mem",
      "description": "Provides system-wide memory stats"
    },
    "net": {
      "name": "net",
      "description": "Provides stats on each network interface"
    },
    "process": {
      "name": "process",
      "description": "Provides a process tree"
    },
    "uptime": {
      "name": "uptime",
      "description": "Provides time since startup"
    }
  }
}
```

```bash
curl http://localhost:5374/?modules=disk,load
```

```js
{
  "disk": {
    "vda": {
      "reads": {
        "completed": 21955,
        "sectors": 1258186,
        "merged": 11170,
        "total_ms": 8302
      },
      "writes": {
        "completed": 170359,
        "sectors": 8927472,
        "merged": 940398,
        "total_ms": 893929
      },
      "io": {
        "in_progress": 0,
        "total_ms": 0,
        "total_weighted": 97388
      }
    },
    "vda1": {
      "reads": {
        "completed": 21781,
        "sectors": 1256794,
        "merged": 11170,
        "total_ms": 8295
      },
      "writes": {
        "completed": 170359,
        "sectors": 8927472,
        "merged": 940398,
        "total_ms": 893929
      },
      "io": {
        "in_progress": 0,
        "total_ms": 0,
        "total_weighted": 97382
      }
    }
  },
  "load": {
    "15m": 0,
    "1m": 0,
    "5m": 0.02
  }
}
```

## Todo

* Write tests
* Standardize and document API
* More in-depth networking and disk stats
* Plugin engine for reading, parsing, and presenting data
