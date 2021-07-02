FROM go:latest

WORKDIR /app

COPY . /app/

RUN go mod tidy

RUN go build main.go && echo "Building done" && ls

CMD ["./main"]