FROM golang:1-bookworm

ARG DATABASE_URL=postgres://postgres:1234@localhost:5432/fineasy?sslmode=disable
ENV DATABASE_URL=${DATABASE_URL}

WORKDIR /usr/fineasy

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN make build

EXPOSE 8080

ENTRYPOINT [ "./bin/app" ]