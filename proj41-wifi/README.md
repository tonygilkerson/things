# Read Me

Run a local HTTP server. One is provided under server directory:

```sh
# start local server
go run ./test-server

# in a different terminal 
tinygo flash -target=pico -stack-size=8kb -monitor main.go
```
