FROM golang:1.19-alpine3.16 AS  build

RUN apk --no-cache add ca-certificates

COPY . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o wisdom .

FROM scratch

COPY --from=build /app/wisdom /

ENTRYPOINT [ "/wisdom" ]