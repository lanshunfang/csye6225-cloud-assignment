package awsshared

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChange(T *testing.T) {

	AWS_REGION := os.Getenv("AWS_REGION")
	AWSSession := GetSession()

	assert.Equal(T, AWSSession.REGION, AWS_REGION, "The AWS Region should be set")
	assert.NotNil(T, AWSSession.SESSION, "The AWS SESSION should be got")

}
