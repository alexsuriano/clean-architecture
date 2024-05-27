# Desafio Clean Architecture

### Objetivo: 
Para este desafio, você precisará criar o usecase de listagem das orders.
Esta listagem precisa ser feita com:
- Endpoint REST (GET /order)
- Service ListOrders com GRPC
- Query ListOrders GraphQL
Não esqueça de criar as migrações necessárias e o arquivo api.http com a request para criar e listar as orders.

Para a criação do banco de dados, utilize o Docker (Dockerfile / docker-compose.yaml), com isso ao rodar o comando docker compose up tudo deverá subir, preparando o banco de dados.
Inclua um README.md com os passos a serem executados no desafio e a porta em que a aplicação deverá responder em cada serviço.

### Como rodar a aplicação:
- Clone o projeto
- Vá até raiz do projeto
- Execute o seguinte comando para subir o banco de dados (MySQL) e o message broker (RabbitMQ):
```
docker compose up
```
- Excute o seguinte comando para aplicar a migration:
```
make migrate-up
```
- Para subir o projeto execute o seguinte comando:
```
go run cmd/ordemsystem/main.go cmd/ordemsystem/wire_gen.go
```
- Aproveite 😎

## Portas úteis:
- Servidor Rest: 8000
- Servidor gRPC: 50051
- Servidor GraphQL: 8099
- RabbitMQ: 5672
- MySQL: 3306

## Outros comandos úteis:
- Para reverter a migration
```
make migrate-down
```
- Para atualizar o schema do GraphQL
```
make graph-gen
```
- Para atulizar o schema do gRPC
```
make proto-gen
```
