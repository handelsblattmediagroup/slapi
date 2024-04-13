FROM golang:1.22.1-alpine AS build
COPY . .
RUN CGO_ENABLED=0 go build \
	-ldflags '-d -s -w' \
	-o /slapi ./cmd/slapi

FROM gcr.io/distroless/static
USER nonroot:nonroot
COPY --from=build --chown=nonroot:nonroot /slapi /slapi
ENTRYPOINT ["/slapi"]
