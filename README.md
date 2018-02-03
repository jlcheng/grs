# Dev

```
$GOPATH
  \
   +-- src/jcheng/grs/.git
         \
          +-- cmd/grsp/prototype.go
          +-- grs/basic.go
```

```
$ go build ./...
$ go test ./...
```

 
## Testing
Run the entire test suite before committing any changes

    $ go test -v ./...
    ?       jcheng/grs/cmd/grsp     [no test files]
    <snipped> 
    PASS
    ok      jcheng/grs/test 0.259s


Run a single test

    $ go test -v ./... -run TestMockCommandMapOk
