FROM golang:latest AS builder

WORKDIR /app

RUN export GO111MODULE=on

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -v -o ./task_manager ./cmd/task_manager

FROM alpine:latest AS runner

COPY --from=builder /app/task_manager .
COPY --from=builder /app/config.yaml ./config.yaml
COPY --from=builder /app/templates ./templates

EXPOSE 8080 3306 27017 6379

CMD [ "./task_manager" ]