FROM golang:1.13
WORKDIR /src/workspace
COPY . /src/workspace
RUN go mod download
RUN go build -o basic main.go

FROM gcr.io/distroless/base:nonroot
COPY --from=0 /src/workspace/basic .
CMD ["./basic"]
