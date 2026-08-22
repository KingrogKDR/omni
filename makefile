genProto:
	@protoc -I=. --go_out=. --go-grpc_out=. proto/*.proto
genBenchProto:
	@protoc -I=. --go_out=. --go-grpc_out=. bench/*.proto
compareWireFormat:
	@go test -run=TestBenchmarkComparison -v ./bench