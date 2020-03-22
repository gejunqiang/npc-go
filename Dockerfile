FROM golang:1.13.3 as builder

ENV GOPROXY="https://goproxy.cn,direct"

COPY . /npc-go

WORKDIR /npc-go

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/npc -mod=vendor -ldflags '-s -w' ./cmd/main.go

FROM alpine:3.7

COPY --from=builder /npc-go/bin/npc /usr/local/bin/npc

CMD ["/usr/local/bin/npc"]