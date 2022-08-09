package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"log"
	"math/big"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
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
	debug           bool
	passwordSize    int
	passwordCharset string
	checkUpper      bool
	checkLower      bool
	checkNumerics   bool
	checkSymbols    bool
)

type Payload struct {
	Secret string `json:"secret" binding:"required"`
}

func init() {
	// Configuration initialisation
	viper.SetDefault("General.debug", false) // That the solution of Life, the Universe, and Encryption
	viper.SetDefault("Password Generation.size", 42)
	viper.SetDefault("Password Generation.charset", "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz !#$%&()*+,-./:;<=>?@[]^_`{|}~")
	viper.SetDefault("Password Generation.check_uppercase", true)
	viper.SetDefault("Password Generation.check_lowercase", true)
	viper.SetDefault("Password Generation.check_numerics", true)
	viper.SetDefault("Password Generation.check_symbols", true)

	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileType)

	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.luccryptous/")
	viper.AddConfigPath("$HOME/.config/luccryptous/")
	viper.AddConfigPath("/etc/luccryptous/")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Config file %s cannot be read\n", configFileName+"."+configFileType)
	} else {
		log.Printf("Config file used: %s\n", viper.ConfigFileUsed())
	}

	viper.SetEnvPrefix("luccrypt")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", " ", "_"))
	viper.AutomaticEnv()

	debug = viper.GetBool("General.debug")
	encodedKey := viper.GetString("General.key")
	passwordSize = viper.GetInt("Password Generation.size")
	passwordCharset = viper.GetString("Password Generation.charset")

	hasUpper, hasLower, hasNumerics, hasSymbols := checkCharsetVerifications(passwordCharset)

	checkUpper = viper.GetBool("Password Generation.check_uppercase") && hasUpper
	checkLower = viper.GetBool("Password Generation.check_lowercase") && hasLower
	checkNumerics = viper.GetBool("Password Generation.check_numerics") && hasNumerics
	checkSymbols = viper.GetBool("Password Generation.check_symbols") && hasSymbols

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

	// Gin debug mode
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	router := gin.Default()

	// Handle CORS, allow all origins
	router.Use(cors.Default())

	// Static routing
	router.Use(static.Serve("/", static.LocalFile("./views", true)))

	// API routing
	api := router.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
		api.GET("/uuid", getUUID)
		api.GET("/pass", getPass)
		api.POST("/crypt", msgCrypt)
	}

	_ = router.Run(":3000")
}

/* Check if password verifications like hasUpper is ok with the current
   password charset */
func checkCharsetVerifications(charset string) (bool, bool, bool, bool) {
	var (
		hasUpper    = false
		hasLower    = false
		hasNumerics = false
		hasSymbols  = false
	)

	for _, c := range charset {
		switch {
		case c >= 48 && c <= 57:
			hasNumerics = true
		case c >= 65 && c <= 90:
			hasUpper = true
		case c >= 97 && c <= 122:
			hasLower = true
		default:
			hasSymbols = true
		}
	}

	return hasUpper, hasLower, hasNumerics, hasSymbols
}

/* Check password validity by checking if password contains all
   mandatory characters sets */
func passwordIsValid(hasUpper, hasLower, hasNumerics, hasSymbols bool) bool {
	if checkUpper && !hasUpper {
		return false
	} else if checkLower && !hasLower {
		return false
	} else if checkNumerics && !hasNumerics {
		return false
	} else if checkSymbols && !hasSymbols {
		return false
	}

	return true
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
	var maxLenChar = big.NewInt(int64(len(passwordCharset)))

	for !passwordIsValid(hasUpper, hasLower, hasNumerics, hasSymbols) {
		hasUpper, hasLower, hasNumerics, hasSymbols = false, false, false, false

		for i := 0; i < n; i++ {
			choice, err := rand.Int(rand.Reader, maxLenChar)
			if err != nil {
				panic("Error at random number generation: " + err.Error())
			}

			buf[i] = passwordCharset[choice.Int64()]

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
