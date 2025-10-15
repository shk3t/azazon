FROM golang:1.25-alpine

ENV APP_NAME="stock"
WORKDIR /${APP_NAME}-app

RUN mkdir -p common $APP_NAME
COPY $APP_NAME/go.mod $APP_NAME/go.sum $APP_NAME
COPY common/go.mod common/go.sum common

RUN go work init common $APP_NAME
RUN cd $APP_NAME && go mod download
RUN cd common && go mod download

COPY $APP_NAME $APP_NAME
COPY common common
COPY .env .

RUN cd $APP_NAME && CGO_ENABLED=0 GOOS=linux go build -o build/main cmd/main.go
WORKDIR /${APP_NAME}-app/$APP_NAME
CMD ["build/main"]
