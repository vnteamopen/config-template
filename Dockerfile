FROM golang:1.18 as gobuilder
WORKDIR /app
COPY ./ .
RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./config-template .;

FROM alpine:latest  
WORKDIR /app/
COPY --from=gobuilder /app/config-template .
CMD ["./config-template", "help"] 
