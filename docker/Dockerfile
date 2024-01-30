FROM golang:1.19

RUN mkdir /code

WORKDIR /code

COPY ./main.go ./main.go
COPY ./src ./src
COPY docker/Makefile .


COPY go.mod go.sum ./
RUN go mod download

RUN go build -o treasury ./main.go
RUN chmod +x ./treasury

EXPOSE 8080


CMD [ "./treasury" ]