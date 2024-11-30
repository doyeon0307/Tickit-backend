FROM golang:1.23.2

WORKDIR /usr/sr/app

ENV JWT_SECRET_KEY=${JWT_SECRET_KEY}
ENV AWS_ACCESS_KEY=${AWS_ACCESS_KEY}
ENV AWS_SECRET_KEY=${AWS_SECRET_KEY}


COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go get -u github.com/doyeon0307/tickit-backend

EXPOSE 7000
CMD [ "./tickit-backend" ]
