#* BUILD STAGE
FROM golang:1.22.3-alpine3.20 AS build_stage
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o ./main ./main.go

#* RUNTIME STAGE
FROM alpine:3.19
WORKDIR /app
COPY --from=build_stage /app/main .
COPY ./keys/id_rsa ./keys/id_rsa
CMD [ "/app/main" ]