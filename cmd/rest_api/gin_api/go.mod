module github.com/ndjordjevic/go-sb/cmd/rest_api/gin_api

go 1.12

require (
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-gonic/gin v1.4.0
	github.com/gocql/gocql v0.0.0-20190423091413-b99afaf3b163
	github.com/golang/protobuf v1.3.1
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/ndjordjevic/go-sb/api/order v0.0.0-20190512100935-8ac66f944e3c
	github.com/ndjordjevic/go-sb/api/price v0.0.0-20190512100935-8ac66f944e3c
	github.com/ndjordjevic/go-sb/cmd/common v0.0.0-20190512100935-8ac66f944e3c
	github.com/olivere/elastic/v7 v7.0.0
	google.golang.org/grpc v1.20.1
	gopkg.in/olahol/melody.v1 v1.0.0-20170518105555-d52139073376
)
