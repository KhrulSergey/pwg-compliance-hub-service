ARG GOLANG_VERSION=1.19

#===========================================================
# Base
#===========================================================
FROM golang:${GOLANG_VERSION}-alpine as base

#===========================================================
# Builder
#===========================================================
FROM base as builder
WORKDIR /app
COPY . ./
ARG BUILD_VERSION
ARG APP_VERSION
ARG PW_GITLAB_USER_LOGIN
ARG PW_GITLAB_USER_TOKEN
RUN apk --no-cache add --update bash
RUN apk --no-cache add --update alpine-sdk
# These line is added to be able to go get an internal dependency
RUN git config --global --global url."https://${PW_GITLAB_USER_LOGIN}:${PW_GITLAB_USER_TOKEN}@gitlab.com".insteadOf "https://gitlab.com"
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.build=${BUILD_VERSION} -X main.version=${APP_VERSION}" -v -o compliance-hub-service cmd/main.go


#===========================================================
# Release
#===========================================================
FROM alpine as release
WORKDIR /
RUN apk --no-cache add --update curl
COPY --from=builder /app/compliance-hub-service /compliance-hub-service
CMD ["./compliance-hub-service"]
