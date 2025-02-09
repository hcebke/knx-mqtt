FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./
COPY models/*.go ./models/
COPY parser/*.go ./parser/
COPY protocols/*.go ./protocols/
COPY utils/*.go ./utils/

# CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o /knx-mqtt

FROM alpine:latest

WORKDIR /app

COPY --from=build /knx-mqtt .

CMD ["/app/knx-mqtt"]