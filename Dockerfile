##############################################

# STEP 1 build exectuable binary

##############################################


FROM golang:alpine

WORKDIR /app

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

WORKDIR /app/cmd/server
RUN go build -o /go-docker

CMD ["/go-docker"]




