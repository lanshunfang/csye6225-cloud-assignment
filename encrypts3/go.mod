module xiaofang.me/gomod/encrypts3

go 1.14

// must also have it in require
// require ( xiaofang.me/gomod/awsshared v0.0.0 )
replace xiaofang.me/gomod/awsshared => ../awsshared

require (
	github.com/aws/aws-sdk-go v1.35.5
	github.com/stretchr/testify v1.6.1
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2
	xiaofang.me/gomod/awsshared v0.0.0
)
