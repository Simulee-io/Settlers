protoc --go_out=../server/proto/board/ --go_opt=paths=source_relative board.proto
protoc --go_opt=paths=source_relative --go_out=plugins=grpc:../server/proto/service/ settlers.proto