package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/alexsuriano/clean-architecture/configs"
	"github.com/alexsuriano/clean-architecture/internal/event/handler"
	"github.com/alexsuriano/clean-architecture/internal/infra/graph"
	"github.com/alexsuriano/clean-architecture/internal/infra/grpc/pb"
	"github.com/alexsuriano/clean-architecture/internal/infra/grpc/service"
	"github.com/alexsuriano/clean-architecture/internal/infra/web/webserver"
	"github.com/alexsuriano/clean-architecture/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg, err := configs.LoadConfig("./configs")
	if err != nil {
		panic(err)
	}

	strConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	dbConn, err := sql.Open(cfg.DBDriver, strConn)
	if err != nil {
		panic(err)
	}
	defer dbConn.Close()

	rabbitMQChannel := getRabbitMQChannel(cfg.RabbitMqUser, cfg.RabbitMqPass, cfg.RabbitMqHost, cfg.RabbitMqPort)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	webserver := webserver.NewWebServer(cfg.WebServerPort)
	webOrderHandler := NewWebOrderHandler(dbConn, eventDispatcher)
	webserver.AddHandler("/order", webOrderHandler.Create)
	webserver.AddHandler("/orders", webOrderHandler.GetAll)
	log.Println("Stating web server on port: ", cfg.WebServerPort)
	go webserver.Start()

	createOrderUseCase := NewCreateOrderUseCase(dbConn, eventDispatcher)
	listOrdersUseCase := NewListOrdersUseCase(dbConn)

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)
	log.Println("Starting gRPC server on port: ", cfg.GRPCServerPort)
	lis, err := net.Listen("tcp", cfg.GRPCServerPort)
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrderUseCase:   *listOrdersUseCase,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Println("Starting GraphQL server on port", cfg.GraphQLServerPort)
	err = http.ListenAndServe(cfg.GraphQLServerPort, nil)
	if err != nil {
		panic(err)
	}

}

func getRabbitMQChannel(user, pass, host, port string) *amqp.Channel {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port))
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
