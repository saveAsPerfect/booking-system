
FROM golang:alpine AS builder


RUN apk add --no-cache git


RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest


WORKDIR /app


COPY . .


RUN go build -o main ./cmd/main.go


ADD https://github.com/vishnubob/wait-for-it/raw/master/wait-for-it.sh /usr/local/bin/wait-for-it
RUN chmod +x /usr/local/bin/wait-for-it


FROM alpine:latest


COPY --from=builder /app/main /app/main
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY --from=builder /usr/local/bin/wait-for-it /usr/local/bin/wait-for-it


COPY migrations /app/migrations


COPY configs /app/configs 

WORKDIR /app


CMD ["wait-for-it", "db:5432", "--", "./main"]
