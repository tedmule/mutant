FROM golang:1.23.0-alpine AS build

# Set workdir
WORKDIR /work

# Add dependencies
COPY go.mod .
COPY go.sum .
RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 go build -o mt

# Generate final image
FROM scratch
COPY --from=build /work/mt /mt
ENTRYPOINT [ "/mt" ]