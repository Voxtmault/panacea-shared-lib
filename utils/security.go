package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"

	"golang.org/x/crypto/pbkdf2"

	config "github.com/voxtmault/panacea-shared-lib/config"
)

const (
	characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+,.?/:;{}[]~"
	// Constants for PBKDF2 parameters
	saltSize       = 16 // You can adjust the salt size as needed
	iterationCount = 10000
	keySize        = 32 // Key size in bytes
)

// type JWTClaims struct {
// 	AccountType string `json:"account_type"`
// 	jwt.RegisteredClaims
// }

// Security Utils
func EncryptAES_CBC(plaintext []byte) (string, error) {

	cfg := config.GetConfig()

	key := []byte(cfg.SecurityConfig.AESKey)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "-1", err
	}

	paddedPlaintext := pkcs7Padding(plaintext, block.BlockSize())

	ciphertext := make([]byte, block.BlockSize()+len(paddedPlaintext))
	iv := ciphertext[:block.BlockSize()]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "-1", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[block.BlockSize():], paddedPlaintext)

	// Encode the ciphertext using base64 and URL-safe encoding
	encodedCiphertext := base64.URLEncoding.EncodeToString(ciphertext)

	return encodedCiphertext, nil
}

// DecryptAES_CBC returns an empty string, "", if the CT provided is invalid. Return 400 if the ciphertext is too short meaning that the encrypted form might be forged or corrupted.
// Returns 500 if something went wrong in the process of decrypting the CT
func DecryptAES_CBC(ciphertext string) (string, error) {

	cfg := config.GetConfig()

	key := []byte(cfg.SecurityConfig.AESKey)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	encrypted, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(encrypted) < aes.BlockSize {
		//Ciphertext too short
		return "", errors.New("400")
	}

	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	if len(encrypted)%aes.BlockSize != 0 {
		//Ciphertext is not a multiple of the block size
		return "", errors.New("500")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encrypted, encrypted)

	plaintext, err := pkcs7Unpadding(encrypted)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func pkcs7Padding(input []byte, blockSize int) []byte {
	padding := blockSize - len(input)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(input, padText...)
}

func pkcs7Unpadding(input []byte) ([]byte, error) {
	length := len(input)
	unpadding := int(input[length-1])
	if length < unpadding {
		return nil, fmt.Errorf("invalid padding")
	}
	return input[:(length - unpadding)], nil
}

// Password Utils
func GenerateRandomPassword(length int) string {
	charLength := big.NewInt(int64(len(characters)))
	password := make([]byte, length)

	for i := 0; i < length; i++ {
		randomIndex, _ := rand.Int(rand.Reader, charLength)
		password[i] = characters[randomIndex.Int64()]
	}

	return string(password)
}

// Password Storing Utils

// hexDecode decodes a hexadecimal string to bytes.
func hexDecode(hexStr string) ([]byte, error) {
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// HashPassword generates a PBKDF2 hash of the password and returns the hash and salt.
func HashPassword(password string) (string, string, error) {
	// Generate a random salt
	salt := make([]byte, saltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return "", "", err
	}

	// Compute the PBKDF2 hash of the password using the salt
	hashedPassword := pbkdf2.Key([]byte(password), salt, iterationCount, keySize, sha256.New)

	// Encode the salt and hash as hexadecimal strings
	saltHex := fmt.Sprintf("%x", salt)
	hashHex := fmt.Sprintf("%x", hashedPassword)

	return hashHex, saltHex, nil
}

// VerifyPassword verifies if a given password matches a stored hash and salt.
func VerifyPassword(password, salt, storedHash string) bool {
	// Decode the salt and stored hash from hexadecimal strings
	saltBytes, _ := hexDecode(salt)
	storedHashBytes, _ := hexDecode(storedHash)

	// Compute the PBKDF2 hash of the input password using the stored salt
	computedHash := pbkdf2.Key([]byte(password), saltBytes, iterationCount, keySize, sha256.New)

	// Use subtle.ConstantTimeCompare to compare the computed hash with the stored hash
	return subtle.ConstantTimeCompare(computedHash, storedHashBytes) == 1
}

func GenerateRandomString() (string, error) {

	length := config.GetConfig().PasswordMinLength

	numBytes := length / 4 * 3
	if length%4 != 0 {
		numBytes = (length/4 + 1) * 3
	}

	temp := make([]byte, numBytes)
	if _, err := rand.Read(temp); err != nil {
		return "", err
	}

	randomString := base64.URLEncoding.EncodeToString(temp)

	randomString = randomString[:length]

	return randomString, nil
}
