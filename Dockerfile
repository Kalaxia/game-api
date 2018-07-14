FROM golang:1.10

WORKDIR /go/src/kalaxia-game-api
COPY . .

RUN go get -d -v ./... \
    && go install -v ./... \
    && go get -u -d github.com/mattes/migrate/cli github.com/lib/pq \
    && go build -tags 'postgres' -o /usr/local/bin/migrate github.com/mattes/migrate/cli \
    && mkdir -p /var/log/api

EXPOSE 80

CMD ["kalaxia-game-api"]
