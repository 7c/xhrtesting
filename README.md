## XHR Testing 
golang/chi stack xhr testing server. Default server is listening `:8080`, you can use `-addr="..."` to define different address to listen. Supports http and https mode

server is designed to tests xhr clients and is available under https://xhrtest.com

## Live Server
[/ping](https://xhrtest.com/ping)- returns pong  
[/status/200](https://xhrtest.com/status/200) - returns statuscode: 200  
[/status/201](https://xhrtest.com/status/201) - returns statuscode: 201  
[/status/204](https://xhrtest.com/status/204) - returns statuscode: 204  
[/status/301](https://xhrtest.com/status/301) - forwards us to /to/301  
[/status/302](https://xhrtest.com/status/302) - forwards us to /to/302  
[/status/400](https://xhrtest.com/status/400) - returns statuscode: 400  
[/status/401](https://xhrtest.com/status/401) - returns statuscode: 401  
[/status/402](https://xhrtest.com/status/402) - returns statuscode: 402  
[/status/403](https://xhrtest.com/status/403) - returns statuscode: 403  
[/status/404](https://xhrtest.com/status/404) - returns statuscode: 404  
[/status/408](https://xhrtest.com/status/408) - returns statuscode: 408  
[/status/500](https://xhrtest.com/status/500) - returns statuscode: 500  
[/status/501](https://xhrtest.com/status/501) - returns statuscode: 501  
[/status/503](https://xhrtest.com/status/503) - returns statuscode: 503  
[/long/body/{number}](https://xhrtest.com/long/body/512) - returns a long body with given number of bytes   
[/json/random](https://xhrtest.com/json/random) - random json response  
[/cookie/random](https://xhrtest.com/cookie/random) - response with random cookie content  
[/cookie/random/{number}](https://xhrtest.com/cookie/random/10) - response with random cookie content at given number  


etc... 


## Compile
`go build -o bin/xhrtesting main.go`

## Running development server
`go run main.go [-addr=":8080"] [-tlsdomain="domain.com"]`

