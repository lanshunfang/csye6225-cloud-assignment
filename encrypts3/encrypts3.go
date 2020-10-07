package encrypts3

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"log"
	"net/http"

	"xiaofang.me/gomod/awsshared"

	"golang.org/x/crypto/scrypt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Config struct {
	S3_BUCKET string
}

func Save(s3Config S3Config, AWSSession awsshared.AWSSession, fileDir string, buffer []byte) {

	_, err := s3.New(AWSSession.SESSION).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(s3Config.S3_BUCKET),
		Key:                  aws.String(fileDir),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(int64(len(buffer))),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})

	if err != nil {
		log.Fatal(err.Error())
		panic(err.Error())
	}
}

func Fetch(s3Config S3Config, AWSSession awsshared.AWSSession, fileDir string) []byte {

	svc := s3.New(AWSSession.SESSION, &aws.Config{
		DisableRestProtocolURICleaning: aws.Bool(true),
	})
	// out: https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#GetObjectOutput
	out, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s3Config.S3_BUCKET),
		Key:    aws.String(fileDir),
	})
	if err != nil {
		// panic(err.Error())
		log.Fatal(err.Error())
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(out.Body)
	newBytes := buf.Bytes()

	return newBytes
}

func Encrypt(aesPassword string, plaintext string) []byte {

	data := []byte(plaintext)
	key, salt, err := DeriveKey(aesPassword, nil)
	if err != nil {
		panic(err.Error())
	}
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	ciphertext = append(ciphertext, salt...)
	return ciphertext
}

func Decrypt(aesPassword string, data []byte) string {

	salt, data := data[len(data)-32:], data[:len(data)-32]
	key, _, err := DeriveKey(aesPassword, salt)
	if err != nil {
		panic(err.Error())
	}
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		panic(err.Error())
	}
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return string(plaintext)

}

func DeriveKey(password string, salt []byte) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, 32)
		if _, err := rand.Read(salt); err != nil {
			return nil, nil, err
		}
	}
	key, err := scrypt.Key([]byte(password), salt, 1048576, 8, 1, 32)
	if err != nil {
		return nil, nil, err
	}
	return key, salt, nil
}
