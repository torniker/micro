module github.com/micro/micro/test/service/stream

go 1.13

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/micro/go-micro/v3 v3.0.0-beta.2.0.20200917131714-7750f542b4c2
	github.com/micro/micro/v3 v3.0.0-beta.3.0.20200910154135-e222e73e9a5c
)

replace github.com/micro/micro/v3 => ../../..
