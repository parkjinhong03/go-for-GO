module auth

go 1.13

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-playground/validator/v10 v10.3.0
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/jinzhu/gorm v1.9.14
	github.com/micro/go-micro/v2 v2.9.0
	github.com/micro/go-plugins/broker/rabbitmq/v2 v2.8.0
	github.com/stretchr/testify v1.4.0
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37
	google.golang.org/protobuf v1.24.0
)
