# Order Test

![build](https://github.com/jordanfinners/order-test/workflows/deploy/badge.svg?branch=master)

All commands haven't been tried/tested on Windows only on Linux.

I added a basic CI/CD pipeline in github actions to run the go linting and tests on every change.
Pushes to master also trigger a deployment to Google Cloud Project, where the endpoints are Cloud Functions, the urls of which are detailed below.

I decided to use Google Cloud Functions as I haven't used them previously and thought it would be good to see what they are like as a cloud offering. I also discovered they have a 'Functions Framework' which can be used to run the functions locally.

## API

Requires Go 1.13 or greater. This can be downloaded [here](https://golang.org/doc/install).

If I were to keep developing the API my next features I would think about adding to store historic orders and provide a `GET orders` and `GET orders/:id` endpoints to list the historic orders.

### Endpoints

There are two endpoints exposed by this API.

#### Orders

This accepts POST requests with a Form Data body. Here is an example request:

```bash
curl 'http://localhost:8080/orders' \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  --data-raw 'items=251'
```

Locally this is accessible at `http://localhost:8080/orders`

Deployed version is accessible at `https://us-central1-ordertest.cloudfunctions.net/orders`

It responds with a JSON response, an example of which is below:

```json
{"items":251,"packs":[{"quantity":500}]}
```

#### Website

This accepts GET requests only. Here is an example request:

```bash
curl 'http://localhost:8080/website' 
```

Locally this is accessible at `http://localhost:8080/website`

Deployed version is accessible at `https://us-central1-ordertest.cloudfunctions.net/website`

It responds with a simple HTML page, where orders can be submitted. It templates in the appropriate url of the orders endpoint.

### Running Locally

Running the functions locally takes advantage of [Googles Functions Framework](https://github.com/GoogleCloudPlatform/functions-framework-go). This will serve up both website and orders endpoints on http://localhost:8080

To run: 

```bash
cd api/local/
go run local.go
```


### Testing
To run the tests, issue the following command(s):

```bash
go fmt ./...
go vet ./...
go test -cover -race ./...
```

To generate a code coverage report issue the following commands:

```bash
go test -cover -coverprofile=c.out ./...
go tool cover -html=c.out -o coverage.html
```

## Data 

This folder hosts the pack information external to the application logic.
This is loaded in on a per request basis by the application.

If I were to keep going on this project this data would be moved to a database, I thought loading it from a remote JSON file would be enough for the concept of externally loaded data for this test. 

A pack is in the format of:
```json
{
    "quantity": 250
}
```
As I thought this would allow in future cost to be added for example and be more extensible than a list of numbers only.
