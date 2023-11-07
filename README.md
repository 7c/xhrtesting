## XHR Testing 
golang/chi stack xhr testing server. Default server is listening `:8080`, you can use `-addr="..."` to define different address to listen. Supports http and https mode

server is designed to tests xhr clients and is available under https://xhrtest.com

## Compile
`go build -o bin/xhrtesting main.go`

## Running development server
`go run main.go [-addr=":8080"] [-tlsdomain="domain.com"]`

