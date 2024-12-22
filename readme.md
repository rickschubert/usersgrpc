GRPC Server
===========

A small project to explore GRPC. Using DynamoDB on Localstack as a database, the server supports 5 different operations and a health check:

```
GetUser(GetUserInput) returns (User)
AddUser(AddUserInput) returns (User)
ModifyUser(ModifyUserInput) returns (User)
RemoveUser(RemoveUserInput) returns (User)
ListUsers(ListUsersInput) returns (ListUsersResponse)
```

The `e2e` directory contains end-to-end test acting as a GRPC client in order to test the features from a user perspective.


# Running and Developing

## E2E Tests with the application running inside Docker

- Start the application and DynamoDB database (via localstack) through docker compose:

```sh
docker compose up --build
```

- Wait until the server is running (it will log "Starting server on port 50051")

- Run the end-to-end tests on your own machine:

```sh
go clean -testcache && LOCAL=true go test -v ./e2e
```

## E2E Tests with the application running locally

- Spin up DynamoDB through localstack:

```sh
docker compose up localstack
```

- Run the server on your own machine:

```sh
LOCAL=true go run ./server
```

- Run the end-to-end tests on your own machine:

```sh
go clean -testcache && LOCAL=true go test -v ./e2e
```

## Developing locally

- Spin up DynamoDB through localstack:

```sh
docker compose up localstack
```

- Run the server on your own machine:

```sh
LOCAL=true go run ./server
```

- Seed the necessary 'users' table in DynamoDB with a bunch of test users:

```sh
LOCAL=true sh seed_db_with_test_data.sh
```

## Recreating users package through protobufs

- The `.proto` definition lives in `users/users.proto`. Whenever you update it, run the following command to re-create the Golang interfaces:

```sh
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative users/users.proto
```
