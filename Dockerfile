FROM golang:1.27-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/go_final_project .

FROM alpine:latest
WORKDIR /app

COPY --from=build /app/go_final_project .
COPY web ./web

EXPOSE 7540
ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

ENTRYPOINT ["/app/go_final_project"]