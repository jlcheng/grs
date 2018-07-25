# Dev (gopath refactoring)

Assume the root of the project is in your $GOPATH

```
grs
 +-- cmd/grs/grs.go
 +-- docs
 +-- src
      +- src/jcheng/grs
      +- src/jcheng/grs/test
 +-- out (build artifacts)

```

```
$ export GOPATH=$PWD
$ cd src/jcheng/grs && dep ensure # install dependencies
$ make test # test only
$ make all  # test and produce artifcats in out
```

 
## Testing
Run the entire test suite before committing any changes

    $ go test -v ./test
    ?       jcheng/grs/cmd/grsp     [no test files]
    <snipped> 
    PASS
    ok      jcheng/grs/test 0.259s


Run a single test

    $ go test -v ./test -run TestMockCommandMapOk
