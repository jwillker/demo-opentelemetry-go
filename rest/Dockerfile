FROM golang:1.13
WORKDIR /src/workspace
COPY . /src/workspace
RUN go mod download
RUN go build -o rest main.go

FROM gcr.io/distroless/base:nonroot
COPY --from=0 /src/workspace/rest .
CMD ["./rest"]
