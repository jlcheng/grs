# Dev

```
$GOPATH
  \
   +-- src/jcheng/grs/.git
         \
          +-- cmd/grsp/prototype.go
          +-- grs/basic.go
```
 
## Testing
```
~$ cd <project root>
.../grs$ go test -v ./..
?       jcheng/grs/cmd/grsp     [no test files]
<snipped> 
PASS
ok      jcheng/grs/test 0.259s
```