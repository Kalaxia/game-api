FROM golang:1.9

WORKDIR /go/src/kalaxia-game-api
COPY . .

RUN go-wrapper download
RUN go-wrapper install

EXPOSE 80

CMD ["go-wrapper", "run"]

