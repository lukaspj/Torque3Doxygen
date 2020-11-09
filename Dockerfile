FROM golang:1.15.4 AS builder

WORKDIR /go/src/app

COPY . .

RUN CGO_ENABLED=0 go build -i -v -o ScriptExecServer

FROM scratch

COPY --from=builder /go/src/app/ScriptExecServer /
COPY --from=builder /go/src/app/files /files

ENTRYPOINT ["/ScriptExecServer"]