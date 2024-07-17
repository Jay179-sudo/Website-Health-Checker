##############################################

# STEP 1 build exectuable binary

##############################################


FROM golang:alpine AS BuildStage

WORKDIR /app

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

WORKDIR /app/cmd/server
RUN go build -o /go-docker


#############################################

# STEP 2 execute binary

#############################################

FROM alpine:latest

WORKDIR /

COPY --from=BuildStage /go-docker /go-docker

ENTRYPOINT ["/go-docker"]




