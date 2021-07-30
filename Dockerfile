FROM golang:1.16 AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
COPY views/ ./views/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .
RUN ls -la && ls -la views

FROM gcr.io/distroless/base-debian10
WORKDIR /
ENV TZ Pacific/Tahiti
COPY --from=build /app/app .
COPY --from=build /app/views ./views
EXPOSE 3000
USER nonroot:nonroot
CMD [ "/app" ]