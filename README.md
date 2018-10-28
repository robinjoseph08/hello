# hello

A dummy Hello World Go API to be used for testing deployments.

## Run

To run the API yourself, you can pull the Docker image and run it with the following commands.

```sh
$ docker run --rm -p 9990:9990 --name hello robinjoseph08/hello:latest
```

## Development

To run this API locally, you can run the following commands.

```sh
$ make setup # needed the first time you clone the repo
$ make install
$ go run cmd/hello/main.go
```
