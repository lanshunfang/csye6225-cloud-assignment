package encrypts3

import (
	"testing"

	"xiaofang.me/gomod/awsshared"

	"github.com/stretchr/testify/assert"
)

func TestChange(T *testing.T) {

	importantMessage := "My important password and credit card"

	AWSSession := awsshared.GetSession()

	aesPassword := "my_Wie1d_pwDDD!"

	s3Config := S3Config{"csye6225-s3-general"}
	filename := "encrypted-text-assignment2.txt"

	encryptedBytes := Encrypt(aesPassword, importantMessage)

	assert.NotEqual(T, importantMessage, encryptedBytes, "The plain text should be encrypted")

	decryptedText := Decrypt(aesPassword, encryptedBytes)
	assert.Equal(T, importantMessage, decryptedText, "The ciphertext should be decrypted")

	Save(s3Config, AWSSession, filename, encryptedBytes)

	fetchedBytes := Fetch(s3Config, AWSSession, filename)

	assert.Equal(T, importantMessage, Decrypt(aesPassword, fetchedBytes), "The ciphertext should be decrypted")

}
