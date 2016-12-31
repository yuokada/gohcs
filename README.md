# gohealthcheck


## How to use

### run server

```
$ ./gohealthcheck  -checklist /path/to/checklist.json
```

## configfile

server checklist

```
# checklist.json
[
  # [<network>, <address> ]
  ["tcp", "localhost:8081"],
  ["tcp", "127.0.0.1:8083"],
  ["tcp", "127.0.0.1:8000"]
]
```

## LICENSE

MIT licensed. See the LICENSE file for details.
