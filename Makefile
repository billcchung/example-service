PROTO_PATH := ./proto

proto-go:
	for FILE in $(PROTO_PATH)/*.proto; do \
		echo "[`date`] Generate pb.go for $${FILE}"; \
		protoc --proto_path=$(PROTO_PATH) --go_out=plugins=grpc:$(PROTO_PATH) $(PROTO_PATH)/*.proto; \
	done
