package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"log"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

const (
	configFileName = "luccryptous"
	configFileType = "toml"
)

var (
	block           cipher.Block
	passwordSize    int
	passwordCharset string
)

type Payload struct {
	Secret string `json:"secret" binding:"required"`
}

func init() {
	// That the solution of Life, the Universe, and Encryption
	viper.SetDefault("Password Generation.size", 42)
	viper.SetDefault("Password Generation.charset", "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz !#$%&()*+,-./:;<=>?@[]^_`{|}~")
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileType)
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.luccryptous/")
	viper.AddConfigPath("$HOME/.config/luccryptous/")
	viper.AddConfigPath("/etc/luccryptous/")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Config file %s not found\n", configFileName+"."+configFileType)
	} else {
		log.Printf("Config file used: %s\n", viper.ConfigFileUsed())
	}

	viper.SetEnvPrefix("luccryptous")
	viper.AutomaticEnv()

	encodedKey := viper.GetString("General.key")
	passwordSize = viper.GetInt("Password_Generation.size")
	passwordCharset = viper.GetString("Password Generation.charset")

	log.Print(passwordSize)

	if len(encodedKey) != 64 {
		panic("Key must be composed of 64 hexadecimal characters")
	}

	key, err := hex.DecodeString(string(encodedKey))
	if err != nil {
		panic("Key must be composed of 64 hexadecimal characters")
	}

	// Cipher block initialisation
	block, err = aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}
}

func main() {
	router := gin.Default()

	router.Use(static.Serve("/", static.LocalFile("./views", true)))

	api := router.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
		api.GET("/guid", getUUID)
		api.GET("/pass", getPass)
		api.POST("/crypt", msgCrypt)
	}

	_ = router.Run(":3000")
}

/* Generate a random password with a secure random number generator,
   passwords have at least one Uppercase letter, one Lowercase letter,
   one Numeric and one Symbol. */
func generateRandomString(n int) ([]byte, error) {
	var (
		hasUpper    = false
		hasLower    = false
		hasNumerics = false
		hasSymbols  = false
	)

	var buf = make([]byte, n)

	for !(hasSymbols && hasNumerics && hasLower && hasUpper) {
		hasUpper, hasLower, hasNumerics, hasSymbols = false, false, false, false

		_, err := rand.Read(buf)
		if err != nil {
			return nil, err
		}

		for i, b := range buf {
			buf[i] = passwordCharset[int(b)%len(passwordCharset)]

			switch {
			case buf[i] >= 48 && buf[i] <= 57:
				hasNumerics = true
			case buf[i] >= 65 && buf[i] <= 90:
				hasUpper = true
			case buf[i] >= 97 && buf[i] <= 122:
				hasLower = true
			default:
				hasSymbols = true
			}
		}
	}

	return buf, nil
}

/* Encrypt plaintext using AES 256 CFB */
func encrypt(plaintext []byte) ([]byte, error) {
	// Buffer for IV + encrypted secret
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	// Initialise a random IV
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

func processEncryption(c *gin.Context, data interface{}) {
	var plaintext []byte

	switch v := data.(type) {
	case string:
		plaintext = []byte(v)
	case []byte:
		plaintext = v
	default:
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
			"message": "Type error at encryption",
		})
		return
	}

	ciphertext, err := encrypt(plaintext)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
			"message": "Error at encryption",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"secret": base64.StdEncoding.EncodeToString(ciphertext),
		})
	}
}

func getUUID(c *gin.Context) {
	if secret, err := uuid.NewRandom(); err != nil {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
			"message": "Error at UUID generation",
		})
	} else {
		processEncryption(c, secret.String())
	}
}

func getPass(c *gin.Context) {
	if secret, err := generateRandomString(passwordSize); err != nil {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
			"message": "Error at password generation",
		})
	} else {
		processEncryption(c, secret)
	}
}

func msgCrypt(c *gin.Context) {
	var payload Payload

	if err := c.BindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "*secret* field is required",
		})
	} else {
		processEncryption(c, payload.Secret)
	}
}
