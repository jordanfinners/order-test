# Order Test

![build](https://github.com/jordanfinners/order-test/workflows/deploy/badge.svg?branch=master)

All commands haven't been tried/tested on Windows only on Linux.

I added a basic CI/CD pipeline in github actions to run the go linting and tests on every change.
Pushes to master also trigger a deployment to Azure, where the api is a Function, the urls of which are detailed below.
Azure Functions don't support Go Functions natively so you have to use an [Azure Custom Handler](https://docs.microsoft.com/en-us/azure/azure-functions/functions-custom-handlers)

I decided to move to Azure from Google Cloud to demonstrate working across 'clouds'.

## API

Requires Go 1.13 or greater. This can be downloaded [here](https://golang.org/doc/install).
Requires Docker. This can be downloaded [here](https://docs.docker.com/engine/install/).

If I were to keep developing the API my next features I would think about adding to store historic orders and provide `orders/:id` endpoints to interact with specific orders e.g. refund, cancel order.

I would also look at securing the endpoints and website more, with security headers, authentication on the endpoint and more protection around what values can be submitted.

### Endpoints

There are is an endpoint exposed by this API.

#### Orders

#### POST

This accepts POST requests with a JSON body. Here is an example request:

```bash
curl 'http://localhost:7071/api/orders' \
  -H 'Content-Type: application/json' \
  --data-raw '{"items":251}'
```

Locally this is accessible at `http://localhost:7071/orders`

Deployed version is accessible at `https://order-test.azurewebsites.net/api/orders`

It responds with a JSON response, an example of which is below:

```json
{"id": "UUID", "items":251,"packs":[{"quantity":500}]}
```

#### GET

This accepts GET requests. Here is an example request:

```bash
curl 'http://localhost:7071/api/orders' \
  -H 'Content-Type: application/json'
```

Locally this is accessible at `http://localhost:7071/orders`

Deployed version is accessible at `https://order-test.azurewebsites.net/api/orders`

It responds with a JSON response, an example of which is below:

```json
[{"id": "UUID", "items":251,"packs":[{"quantity":500}]}]
```

### Architecture 

The entry point to the functions is in the `server/main.go` which handles all requests.
Those to `/orders` will be passed off to the `router/router.go` and has to be configured with the file `orders/function.json`.
The router directs it to the correct handler and converts the inbound function request to simplified request and a wrapper to handle converting the response to the required outbound format for the Azure Custom Handler.

The handlers deal with the simplified request and responds, and runs the business logic.
The handlers use the storage layer, which in the tests is a local mongodb instance span up by `/api/storage/test_db.go`.
In Azure this will connect to a Cosmos DB instance with the mongo API.

### Running Locally

Running the functions locally takes advantage of [Azure Functions Core Tools](https://docs.microsoft.com/en-us/azure/azure-functions/functions-run-local?tabs=linux%2Ccsharp%2Cbash). This will allow you to run the functions locally but it does need to be built first.

It needs to have a local.settings.json file creating at the API folder layer. It should look like the following, filling in the values to connect it to the development environment:
```json
{
    "IsEncrypted": false,
    "Values": {
        "FUNCTIONS_WORKER_RUNTIME": "Custom",
        "AzureWebJobsStorage": "YOUR_STORAGE_CONNECTION_STRING",
        "DATABASE_CONNECTION_STRING": "mongodb://",
        "DATABASE_NAME": "staging"
    }
}
```

You can use either a local mongodb docker image using:
```bash
docker pull mongo
docker run -d -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=password --name db mongo
```
Then set the variable `DATABASE_CONNECTION_STRING` to `mongodb://admin:password@localhost:27017/` the `DATABASE_NAME` can be set to any value.

Or you can connect it to a development instance on Azure using the connection string for the Cosmos DB instance, which will need to have a database created and 2 collections - `packs` and `orders`.

To run: 

```bash
cd /api/server
CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o ../main

cd ../
func host start
```

### Building 

To build a Linux executable required by Azure Functions
```bash
CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o ../main
```

### Testing

To run the tests, issue the following command(s) in the api directory:

```bash
go fmt ./...
go vet ./...
go test -cover -race ./...
```

To generate a code coverage report issue the following commands in the api directory:

```bash
go test -cover -coverprofile=c.out ./...
go tool cover -html=c.out -o coverage.html
```
