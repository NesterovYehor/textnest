module github.com/NesterovYehor/TextNest/services/api_service

go 1.23.1

replace github.com/NesterovYehor/TextNest/pkg => ../../pkg

require (
	github.com/NesterovYehor/TextNest/pkg v0.0.0-20241030211018-6cf5f307cfce
	google.golang.org/grpc v1.67.1
)

require (
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)
