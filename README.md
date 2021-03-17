# gitreefs

Virtual file system, mapping from a directory of clones to all of their possible contents

![icon](./gitreefs.png)

## Packages

### Executables
- [fuse](fuse) - Tool to run a virtual fs using FUSE, using [jacobsa/fuse](https://github.com/jacobsa/fuse).
- [nfs](nfs) - Tool to run a virtual fs using NFS server, using [willscott/go-nfs](https://github.com/willscott/go-nfs).
### Core
- [git](core/git) - Layer to access git data, using [go-git](https://github.com/go-git/go-git).
- [bfs](core/virtualfs/bfs) - Implementation of virtual git fs over [go-billy](https://github.com/go-git/go-billy).
- [inodefs](core/virtualfs/inodefs) - Implementation of virtual git fs using inodes abstraction, as suiting `jacobsa/fuse`.

```
/disk/git/               ===>  /mnt/git/
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

## Tests

```bash
go test -v ./...
```

## NFS solution

```bash
go run gitreefs/nfs --help
go run gitreefs/nfs --log-level INFO /var/git
# then
mkdir -p /tmp/git
mount -o port=2049,mountport=2049 -t nfs localhost:/ /tmp/git
umount /tmp/git
```

```bash
NAME:
   gitreefs-nfs - NFS server providing access to a forest of git trees as a virtual file system

USAGE:
   gitreefs-nfs [Options] clones-path storage-path [port]

ARGS:
    clones-path   path to a directory containing git clones (with .git in them)
    storage-path  path to a directory in which to keep persistent storage (file handler mapping)
    port          (optional) to serve the server at, defaults to 2049

OPTIONS:
   --log-file value   Output logs file path format. (default: "logs/gitreefs-%v-%v.log")
   --log-level value  Set log level. (default: "DEBUG")
   --help, -h         show help
   --version, -v      print the version
```

### Benchmark

Currently not that good

https://github.com/apiirolab/EVO-Exchange-BE-2019

```
+--------------------------------------------+------------+-------------+
|                 Operation                  | Time (sec) | Degradation |
+--------------------------------------------+------------+-------------+
| Walk physical clone                        |      0.002 |             |
| Walk physical clone + read all files       |       0.02 |             |
| Walk virtual #1 iteration                  |       0.14 | x70         |
| Walk virtual #2 iteration                  |       0.09 | x45         |
| Walk virtual + read all files #1 iteration |       0.28 | x14         |
| Walk virtual + read all files #2 iteration |       0.27 | x13.5       |
+--------------------------------------------+------------+-------------+
```

https://github.com/apiirolab/elasticsearch

```
+--------------------------------------------+------------+-------------+
|                 Operation                  | Time (sec) | Degradation |
+--------------------------------------------+------------+-------------+
| Walk physical clone                        |        0.4 |             |
| Walk physical clone + read all files       |        3.6 |             |
| Walk virtual #1 iteration                  |       44.85 | x112       |
| Walk virtual #2 iteration                  |       28.73 | x72        |
| Walk virtual + read all files #1 iteration |      109.96 | x30.5      |
| Walk virtual + read all files #2 iteration |      108.41 | x30.1      |
+--------------------------------------------+------------+-------------+
```

## FUSE solution

```bash
go run gitreefs/fuse --help
go run gitreefs/fuse --log-level INFO /var/git /mnt/git
```

```bash
NAME:
   gitreefs-fuse - Mount a forest of git trees as a virtual file system

USAGE:
   gitreefs-fuse [Options] clones-path mount-point

ARGS:
    clones-path  path to a directory containing git clones (with .git in them)
    mount-point  path to target location to mount the virtual fuseserver at

OPTIONS:
   --log-file value   Output logs file path format. (default: "logs/gitreefs-%v-%v.log")
   --log-level value  Set log level. (default: "DEBUG")
   --help, -h         show help
   --version, -v      print the version
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
