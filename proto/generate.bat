protoc --go_out=../server/proto/ --go_opt=paths=source_relative board.proto
protoc --go_opt=paths=source_relative --go_out=plugins=grpc:../server/proto/ settlers.proto