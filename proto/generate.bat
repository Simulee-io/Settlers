protoc --go_out=../server/datafile/ --go_opt=paths=source_relative board.proto
protoc --go_opt=paths=source_relative --go_out=plugins=grpc:../server/service/ settlers.proto