# gitreefs

Virtual file system, mapping from a directory of clones to all of their possible contents:

```bash
NAME:
   gitreefs - Mount a forest of git trees as a virtual file system

USAGE:
   gitreefs [global Options] clones mountpoint
GLOBAL OPTIONS:
   --log-file value   Output logs file path format. (default: "logs/gitreefs-%v-%v.log")
   --log-level value  Set log level. (default: "DEBUG")
   --help, -h         show help
   --version, -v      print the version
```

![icon](./gitreefs.png)

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

<div>Icons made by <a href="https://www.freepik.com" title="Freepik">Freepik</a> from <a href="https://www.flaticon.com/" title="Flaticon">www.flaticon.com</a></div>
