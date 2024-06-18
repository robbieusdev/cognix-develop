# select image
FROM golang:1.22.1-alpine

# install dependensies and migration tool
RUN apk add --no-cache build-base git && go install github.com/pressly/goose/cmd/goose@latest

# copy migration files.
COPY versions /versions

# run migration
WORKDIR /versions
CMD goose postgres "${COCKROACH_CLIENT_DATABASE_URL_GOLANG}" up
