FROM golang:alpine AS build

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o my-app

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/my-app ./

ENTRYPOINT ["./my-app"]

# EXPOSE 5000
#sudo docker build --pull --rm -f "Dockerfile" -t serchiboi/backend_vmgo_p2so1:latest "./"
#sudo docker push serchiboi/backend_vmgo_p2so1