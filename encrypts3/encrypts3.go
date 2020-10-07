package encrypts3

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/scrypt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var password = ""

var symmetryKeyHex = "6D795F766572795F703077657246756C215F736563526554"

type S3SessionConfig struct {
	S3_REGION   string
	S3_BUCKET   string
	session     *session.Session
	aesPassword []byte
}

func GetSession(s3Bucket string, rawPassword string) S3SessionConfig {
	AWS_REGION := os.Getenv("AWS_REGION")

	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
		})
	s, err := session.NewSession(&aws.Config{
		Region:      aws.String(AWS_REGION),
		Credentials: creds,
	})
	if err != nil {
		log.Fatal(err)
	}

	return S3SessionConfig{AWS_REGION, s3Bucket, s, []byte(rawPassword)}
}

func Save(s3Config S3SessionConfig, fileDir string, buffer []byte) {

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	_, err := s3.New(s3Config.session).PutObject(&s3.PutObjectInput{
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

func Fetch(s3Config S3SessionConfig, fileDir string) []byte {

	svc := s3.New(s3Config.session, &aws.Config{
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

func Encrypt(s3SessionConfig S3SessionConfig, plaintext string) []byte {

	data := []byte(plaintext)
	key, salt, err := DeriveKey(s3SessionConfig.aesPassword, nil)
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

func Decrypt(s3SessionConfig S3SessionConfig, data []byte) string {

	salt, data := data[len(data)-32:], data[:len(data)-32]
	key, _, err := DeriveKey(s3SessionConfig.aesPassword, salt)
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

func DeriveKey(password, salt []byte) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, 32)
		if _, err := rand.Read(salt); err != nil {
			return nil, nil, err
		}
	}
	key, err := scrypt.Key(password, salt, 1048576, 8, 1, 32)
	if err != nil {
		return nil, nil, err
	}
	return key, salt, nil
}
