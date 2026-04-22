docker run --rm -it \
  -v "$PWD":/usr/src/app \
  -w /usr/src/app/agent \
  -e SERVER_URL=http://host.docker.internal:8080/api/v1/metrics \
  -e COLLECTION_INTERVAL=2s \
  -e REQUEST_TIMEOUT=3s \
  golang:1.25.5 \
  go run ./cmd/agent

to run agent from it's dir

and server: from server/cmd/server 
# go run main.go
