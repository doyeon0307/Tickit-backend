FROM golang:1.23.2

WORKDIR /usr/sr/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go get -u github.com/doyeon0307/tickit-backend

EXPOSE 7000
CMD [ "./tickit-backend" ]
