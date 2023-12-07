# ops_tools

### 介绍
运维常用工具开发(Develop daily ops tools)


### go编译构建

```shell
go build  .
```

1. Mac下编译linux、windows平台的包

```shell
$ go env -w CGO_ENABLED=0 GOOS=linux GOARCH=amd64
$ go env -w CGO_ENABLED=0 GOOS=windows GOARCH=amd64
```

2. Linux下编译Mac和Windows平台的包

```shell
$ go env -w CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 
$ go env -w CGO_ENABLED=0 GOOS=windows GOARCH=amd64 
```

3. Windows下编译Mac和Linux平台的包

```shell
$ go env -w CGO_ENABLED=0 GOOS=darwin3 GOARCH=amd64 
$ go env -w CGO_ENABLED=0 GOOS=linux GOARCH=amd64 
```

### go交叉编译构建

```shell
go get github.com/mitchellh/gox@latest
go install github.com/mitchellh/gox@latest
```

- 构建一个包和子包

```shell
gox ./...
```

- 构建不同的包

```shell
gox github.com/mitchellh/gox github.com/hashicorp/serf
```

- 只构建linux

```shell
 gox -os="linux"
```

- 只构建64位的linux

```shell
gox -osarch="linux/amd64"
```

