FROM golang:1.9

WORKDIR /go/src/kalaxia-game-api
COPY . .

RUN go-wrapper download \
    && go-wrapper install \
    && go get -u -d github.com/mattes/migrate/cli github.com/lib/pq \
    && go build -tags 'postgres' -o /usr/local/bin/migrate github.com/mattes/migrate/cli

EXPOSE 80

CMD ["go-wrapper", "run"]
