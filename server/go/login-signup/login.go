package loginsignup

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Login struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

func UserLogIn(c *gin.Context) {
	var user Login

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error: ": err.Error()})
	}

	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error: ": err.Error()})
	}

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
	}

	cipherTxt := gcm.Seal(nonce, nonce, []byte(user.Password), nil)
	user.Password = string(cipherTxt)
	c.JSON(200, gin.H{
		"message":  "data received",
		"username": user.UserName,
		"pass":     user.Password,
	})
}
