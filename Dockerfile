FROM golang:1.13-alpine AS build-env

ENV GO111MODULE=on

WORKDIR /go/src/kalaxia-game-api

RUN apk add git make gcc g++

RUN go get -u -d github.com/mattes/migrate/cli github.com/lib/pq \
    && go build -tags 'postgres' -o /usr/local/bin/migrate github.com/mattes/migrate/cli

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
COPY docker-entrypoint.sh /entrypoint.sh

RUN chmod a+x /entrypoint.sh \
    && chown root:root /entrypoint.sh
RUN go build -o /usr/local/bin/kalaxia-game-api .

EXPOSE 80

ENTRYPOINT ["/entrypoint.sh"]

CMD ["kalaxia-game-api"]