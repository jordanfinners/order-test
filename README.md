# Order Test

![build](https://github.com/jordanfinners/order-test/workflows/build/badge.svg?branch=master)

This repo is split into two folders.
All commands haven't been tried/tested on Windows only on Linux.

I added a basic CI pipeline in github actions to run the go linting and test on every change, as it was mentioned in the job description.

## API

Requires Go 1.11 or greater. This can be downloaded [here](https://golang.org/doc/install).

This assumes you are using python 3 and the appropriate pip version for it.

To run:

```bash
go fmt ./...
go vet ./...
go test -cover -race ./...
```


## Data 
