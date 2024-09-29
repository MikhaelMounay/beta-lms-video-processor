package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	// Get inputs from user
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the path of the video file: ")
	scanner.Scan()
	filePath := strings.ReplaceAll(scanner.Text(), "\"", "")

	fmt.Print("Enter the link of the YouTube video: ")
	scanner.Scan()
	youtubeLink := scanner.Text()
	encryptedPath := GetEncryptedFileName(filePath, youtubeLink)

	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		return
	}

	// get keys from env file
	secretKey := os.Getenv("ENC_SECRET_KEY_HEX")
	iv := os.Getenv("ENC_IV_HEX")

	if err := EncryptFile(filePath, encryptedPath, secretKey, iv); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
