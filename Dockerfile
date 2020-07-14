FROM golang:1.13-alpine as build

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cloudeye-exporter

FROM gcr.io/distroless/static:nonroot
USER nonroot:nonroot

COPY --from=build /app/cloudeye-exporter /usr/local/bin/cloudeye-exporter 

ENTRYPOINT [ "cloudeye-exporter" ]
