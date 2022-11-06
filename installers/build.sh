protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative protos/chat.proto
cd client && go build .
cd ../server && go build .
cd ../

