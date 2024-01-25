FROM golang as build
WORKDIR /app
COPY . .
RUN go install
RUN go build -o spf


FROM debian:bookworm-slim
WORKDIR /app
COPY --from=build /app/spf spf
COPY --from=build /app/default-server-conf.yml conf.yml
CMD ["./spf", "server", "-c", "conf.yml"]
