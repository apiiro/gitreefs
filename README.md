# gitreefs

Virtual file system, mapping from a directory of clones to all of their possible contents

![icon](./gitreefs.png)

## Tests

```bash
go test -v ./...
```

## FUSE solution

```bash
go run gitreefs/gitree-fuse --help
go run gitreefs/gitree-fuse --log-level INFO /var/git /mnt/git
```

```bash
NAME:
   gitreefs - Mount a forest of git trees as a virtual file system

USAGE:
   gitreefs [Options] clones-path mount-point

ARGS:
    clones-path  path to a directory containing git clones (with .git in them)
    mount-point  path to target location to mount the virtual fuseserver at

OPTIONS:
   --log-file value   Output logs file path format. (default: "logs/gitreefs-%v-%v.log")
   --log-level value  Set log level. (default: "DEBUG")
   --help, -h         show help
   --version, -v      print the version
```

```
/disk/git/                =>    /mnt/git/
          | clone1/.git                 | clone1/
                                             | commit1/
                                                       | src/...
                                             | commit2/
                                                       | src/...
                                             | commit3/
                                                       | src/...
          | clone2/.git
          | clone3/.git
          | ...
```

### Open Issues

- Performance - can add caching, either in memory of physical fs based
- Memory usage - currently nothing allocated will ever be released. Can add interval clean up to swipe away unused roots (in repository or commitish level).

### Benchmark

Currently not that good

https://github.com/apiirolab/EVO-Exchange-BE-2019

```
+--------------------------------------------+------------+-------------+
|                 Operation                  | Time (sec) | Degradation |
+--------------------------------------------+------------+-------------+
| Walk physical clone                        |      0.002 |             |
| Walk physical clone + read all files       |       0.02 |             |
| Walk virtual #1 iteration                  |       0.07 | x35         |
| Walk virtual #2 iteration                  |       0.03 | x15         |
| Walk virtual + read all files #1 iteration |       0.18 | x9          |
| Walk virtual + read all files #2 iteration |       0.11 | x5.5        |
+--------------------------------------------+------------+-------------+
```

https://github.com/apiirolab/elasticsearch

```
+--------------------------------------------+------------+-------------+
|                 Operation                  | Time (sec) | Degradation |
+--------------------------------------------+------------+-------------+
| Walk physical clone                        |        0.4 |             |
| Walk physical clone + read all files       |        4.4 |             |
| Walk virtual #1 iteration                  |       35.3 | x88         |
| Walk virtual #2 iteration                  |       24.4 | x61         |
| Walk virtual + read all files #1 iteration |      161.3 | x36.6       |
| Walk virtual + read all files #2 iteration |      136.2 | x31         |
+--------------------------------------------+------------+-------------+
```

## Credits
<div>Icons made by <a href="https://www.freepik.com" title="Freepik">Freepik</a> from <a href="https://www.flaticon.com/" title="Flaticon">www.flaticon.com</a></div>
