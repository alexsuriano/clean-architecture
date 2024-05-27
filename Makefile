.PHONY: migrate proto

migrate-create:
	migrate create -ext=sql -dir=sql/migrations -seq init

migrate-up:
	migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3306)/orders" -verbose up

migrate-down:
	migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3306)/orders" -verbose down

proto-gen:
	protoc --go_out=. --go-grpc_out=.  internal/infra/grpc/protofiles/order.proto

graph-gen:
	go run github.com/99designs/gqlgen generate