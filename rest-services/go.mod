module github.com/pinas/rest-services

go 1.26.1

require (
	github.com/labstack/echo/v5 v5.2.1
	github.com/pinas/common-structs v0.0.0-00010101000000-000000000000
	github.com/r3labs/sse/v2 v2.10.0
)

replace github.com/pinas/common-structs => ../common-structs

require (
	golang.org/x/net v0.56.0 // indirect
	golang.org/x/time v0.15.0 // indirect
	gopkg.in/cenkalti/backoff.v1 v1.1.0 // indirect
)
