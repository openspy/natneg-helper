FROM golang:latest as build
WORKDIR /app
COPY go.mod go.sum .
COPY src/ src
RUN go mod download
WORKDIR /app/src
RUN GOOS=linux go build -o ../natneg-helper

FROM golang:latest
WORKDIR /app
COPY --from=build /app/natneg-helper natneg-helper
ENTRYPOINT ["/app/natneg-helper"]