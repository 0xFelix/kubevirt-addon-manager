# Build the manager binary
FROM golang:1.19 as builder

WORKDIR /workspace

# Copy the Go Modules manifests and vendor directory
COPY go.mod go.mod
COPY go.sum go.sum
COPY vendor/ vendor/

# Copy the go source
COPY Makefile Makefile
COPY main.go main.go
COPY pkg/ pkg/

RUN make build

FROM registry.access.redhat.com/ubi9/ubi-minimal

RUN microdnf update -y && microdnf clean all

WORKDIR /

COPY --from=builder /workspace/bin/manager .

USER 65534:65534

ENTRYPOINT ["/manager"]
CMD ["controller"]
