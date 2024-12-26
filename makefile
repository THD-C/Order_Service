PROTO_SRC_DIR=./Protocol/proto
PROTO_OUT_DIR=generated

run:
	go run main.go

proto:
	mkdir -p $(PROTO_OUT_DIR)
	find $(PROTO_SRC_DIR) -name "*_*.proto" -exec protoc -I=$(PROTO_SRC_DIR) \
	--go_out=$(PROTO_OUT_DIR) --go_opt=paths=source_relative \
	--go-grpc_out=$(PROTO_OUT_DIR) --go-grpc_opt=paths=source_relative {} \;
	find $(PROTO_SRC_DIR) -name "*.proto" ! -name "*_*.proto" -exec protoc -I=$(PROTO_SRC_DIR) \
	--go_out=$(PROTO_OUT_DIR) --go_opt=paths=source_relative \
	--go-grpc_out=$(PROTO_OUT_DIR) --go-grpc_opt=paths=source_relative {} \;

.PHONY: run proto
