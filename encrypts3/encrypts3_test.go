package encrypts3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChange(T *testing.T) {

	s3SessionConfig := GetSession(
		"csye6225-s3-general",
		"my_Wie1d_pwDDD!",
	)

	filename := "encrypted-text-assignment2.txt"

	importantMessage := "My important password and credit card"

	encryptedBytes := Encrypt(s3SessionConfig, importantMessage)

	assert.NotEqual(T, importantMessage, encryptedBytes, "The plain text should be encrypted")

	decryptedText := Decrypt(s3SessionConfig, encryptedBytes)
	assert.Equal(T, importantMessage, decryptedText, "The ciphertext should be decrypted")

	Save(s3SessionConfig, filename, encryptedBytes)

	fetchedBytes := Fetch(s3SessionConfig, filename)

	assert.Equal(T, importantMessage, Decrypt(s3SessionConfig, fetchedBytes), "The ciphertext should be decrypted")

}
