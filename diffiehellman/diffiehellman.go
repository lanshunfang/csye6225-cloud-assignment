package diffiehellman

import (
	"math"
)

var base float64 = 4
var modulusInPrime float64 = 17

/*
EncryptSharedKey ...
*/
func EncryptSharedKey(secretNumber float64) float64 {
	return math.Mod(math.Pow(base, secretNumber), modulusInPrime)
}

/*
DecryptSharedKey ...
*/
func DecryptSharedKey(encryptedNumber float64, secretNumber float64) float64 {
	return math.Mod(math.Pow(encryptedNumber, secretNumber), modulusInPrime)
}
