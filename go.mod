module gitreefs

go 1.14

require (
	github.com/dgraph-io/badger v1.6.2
	github.com/go-git/go-billy/v5 v5.0.0
	github.com/go-git/go-git/v5 v5.2.0
	github.com/google/uuid v1.2.0
	github.com/jacobsa/fuse v0.0.0-20201216155545-e0296dec955f
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/orcaman/concurrent-map v0.0.0-20210106121528-16402b402231
	github.com/stretchr/testify v1.4.0
	github.com/urfave/cli v1.22.5
	github.com/willscott/go-nfs v0.0.0-20210308004034-50941b6e35e1
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110
	golang.org/x/sys v0.0.0-20210305230114-8fe3ee5dd75b // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

replace github.com/willscott/go-nfs v0.0.0-20210308004034-50941b6e35e1 => github.com/apiiro/go-nfs v0.0.0-20210317140244-72fc6f7c68d7

//replace github.com/willscott/go-nfs v0.0.0-20210308004034-50941b6e35e1 => ../go-nfs
