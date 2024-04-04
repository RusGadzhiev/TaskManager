FROM golang:latest AS builder

WORKDIR /apps

RUN export GO111MODULE=on

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY . .

RUN go build -v -o ./task_manager ./cmd/task_manager

FROM alpine:latest AS runner

COPY --from=builder /apps/task_manager .
COPY --from=builder /apps/config.yaml ./config.yaml

EXPOSE 8080 3306 27017 6379

CMD [ "./task_manager" ]