# confd

* `what I want`: Learning the src code by modifying them on my own way
* `target`: simplify, only etcd/etcdv3 is needed (without client security options)
* `todo`: first support backend etcd/etcdv3 only, then improve(config, test...)
* `status`: dev 


`confd` is a lightweight configuration management tool focused on:

* keeping local configuration files up-to-date using data stored in [etcd](https://github.com/coreos/etcd)
* reloading applications to pick up new config file changes



## Building

Go 1.8 is required to build confd, which uses the new vendor directory.

```
$ mkdir -p $GOPATH/src/github.com/opera443399
$ cd $GOPATH/src/github.com/opera443399
$ git clone https://github.com/opera443399/confd.git
$ cd confd
```

dep is needed
```
$ go get github.com/golang/dep/cmd/dep
$ make
```

You should now have confd in your `bin/` directory:

```
$ ls bin/
confd
```

