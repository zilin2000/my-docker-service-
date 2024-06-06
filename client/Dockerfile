FROM golang:1.16

WORKDIR /client/

COPY ../client/main.go .
COPY ../go.mod .
COPY ../go.sum .

COPY ../proto ./proto


RUN go mod download

CMD ["go", "run", "main.go"]
