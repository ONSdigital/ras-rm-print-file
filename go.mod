module github.com/ONSdigital/ras-rm-print-file

go 1.19

require (
	cloud.google.com/go/datastore v1.1.0
	cloud.google.com/go/pubsub v1.3.1
	cloud.google.com/go/storage v1.10.0
	github.com/blendle/zapdriver v1.3.1
	github.com/gorilla/mux v1.8.0
	github.com/pkg/sftp v1.12.0
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	golang.org/x/sys v0.3.0 // indirect
	google.golang.org/api v0.28.0
	google.golang.org/grpc v1.29.1
)
