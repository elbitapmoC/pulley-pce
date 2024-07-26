package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Challenge struct {
	Challenger       string `json:"challenger"`
	EncryptedPath    string `json:"encrypted_path"`
	EncryptionMethod string `json:"encryption_method"`
	ExpiresIn        string `json:"expires_in"`
	Hint             string `json:"hint"`
	Instructions     string `json:"instructions"`
	Level            int    `json:"level"`
}

func decodeASCIIValues(asciiStr string) (string, error) {
	asciiStr = strings.Trim(asciiStr, "[]")
	asciiValuesStr := strings.Split(asciiStr, ",")

	decodedRunes := make([]rune, len(asciiValuesStr))
	for i, val := range asciiValuesStr {
		intVal, err := strconv.Atoi(strings.TrimSpace(val))
		if err != nil {
			return "", err
		}
		decodedRunes[i] = rune(intVal)
	}
	return string(decodedRunes), nil
}

func getChallenge(c echo.Context) error {
	email := "bazellMP@gmail.com"
	url := fmt.Sprintf("https://ciphersprint.pulley.com/%s", email)

	resp, err := http.Get(url)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()

	var challenge Challenge
	if err := json.NewDecoder(resp.Body).Decode(&challenge); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Decrypt the encrypted_path if the encryption method is ASCII conversion
	if challenge.EncryptionMethod == "converted to a JSON array of ASCII values" {
		decryptedPath, err := decodeASCIIValues(challenge.EncryptedPath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		challenge.EncryptedPath = decryptedPath
	}

	// Return the challenge details with the decrypted path
	return c.JSON(http.StatusOK, challenge)
}

func followChallenge(c echo.Context) error {
	encryptedPath := c.Param("encryptedPath")
	subsequentURL := fmt.Sprintf("https://ciphersprint.pulley.com/%s", encryptedPath)

	subsequentResp, err := http.Get(subsequentURL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer subsequentResp.Body.Close()

	subsequentBody, err := ioutil.ReadAll(subsequentResp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Return the response from the subsequent request
	return c.JSON(http.StatusOK, string(subsequentBody))
}

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{http.MethodGet},
	}))

	e.GET("/get-challenge", getChallenge)
	e.GET("/follow-challenge/:encryptedPath", followChallenge)
	e.Logger.Fatal(e.Start(":8080"))
}
