FROM golang:1.23.5 AS build

ARG CGO_ENABLED=0

WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./
RUN CGO_ENABLED=${CGO_ENABLED} go build -ldflags="-w" -o /app .


FROM alpine AS final

RUN apk add git

LABEL maintainer="soerenschneider"
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=build /app /app
USER appuser

WORKDIR /repo

ENTRYPOINT ["/app"]
