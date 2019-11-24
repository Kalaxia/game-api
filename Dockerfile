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

RUN go build .

FROM alpine

RUN apk add make

WORKDIR /go/src/kalaxia-game-api

COPY --from=build-env /usr/local/bin/migrate /usr/local/bin/migrate
COPY --from=build-env /go/src/kalaxia-game-api .
COPY --from=build-env /go/src/kalaxia-game-api/kalaxia-game-api /usr/local/bin/kalaxia-game-api
RUN mkdir -p /var/log/api

EXPOSE 80

CMD ["kalaxia-game-api"]