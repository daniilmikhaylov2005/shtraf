FROM golang:1.19

WORKDIR /testShtraf

COPY go.mod go.sum ./

RUN go mod download

COPY * ./
RUN  go get github.com/daniilmikhaylov2005/shtraf/api
RUN  go get github.com/daniilmikhaylov2005/shtraf/server/models
RUN CGO_ENABLED=0 GOOS=linux go build -o /testShtraf

EXPOSE 8800

CMD ["/testShtraf"]