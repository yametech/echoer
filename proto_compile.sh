protoc -I api api/api.proto --gofast_out=plugins=grpc:api
protoc -I proto proto/echoer.proto --gofast_out=plugins=grpc:proto