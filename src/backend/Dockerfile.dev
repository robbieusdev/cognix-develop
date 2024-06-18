# Dockerfile

FROM golang:1.22-alpine
ARG service
ENV SERVICE_NAME $service
WORKDIR /backend/${SERVICE_NAME}
COPY ../go.mod ../go.sum ../
RUN go install github.com/cosmtrek/air@latest
RUN go mod download
CMD ["air", "-d"]
