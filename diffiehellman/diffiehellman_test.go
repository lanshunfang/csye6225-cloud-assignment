package diffiehellman

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChange(T *testing.T) {
	var secrentNum1 float64 = 3
	var secrentNum2 float64 = 6
	aliceSentEncryptedSharedKey := EncryptSharedKey(secrentNum1)
	bobSentEncryptedSharedKey := EncryptSharedKey(secrentNum2)

	aliceDecryptedSharedKey := DecryptSharedKey(bobSentEncryptedSharedKey, secrentNum1)
	bobDecryptedSharedKey := DecryptSharedKey(aliceSentEncryptedSharedKey, secrentNum2)

	assert.Equal(T, aliceDecryptedSharedKey, bobDecryptedSharedKey, "Two parties should be able to encrypt and decrypt shared key")

}
