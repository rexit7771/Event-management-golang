from golang:1.22.0-alpine

RUN apk add --no-cache gcc musl-dev curl git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download -x

COPY . .

RUN go build -gcflags="all=-N -l" -o main .

EXPOSE 8080

CMD ["./main"]