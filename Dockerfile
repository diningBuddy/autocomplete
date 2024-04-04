FROM public.ecr.aws/docker/library/golang:1.18 as builder

WORKDIR /app

ARG ACCESS_TOKEN

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod .
COPY go.sum .

RUN git config --global url."https://bp-deployer:$ACCESS_TOKEN@github.com/".insteadOf "https://github.com/"

ENV GOPRIVATE=github.com/skku

RUN go mod download

COPY . .

RUN go build

FROM public.ecr.aws/docker/library/alpine:latest as app

WORKDIR /app
COPY --from=builder app /app

ENTRYPOINT ["./autocomplete"]