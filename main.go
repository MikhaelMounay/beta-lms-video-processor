package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.design/x/clipboard"
)

func main() {
	// Load environment variables from the .env file
	if err := godotenv.Load(); err != nil {
		panic(fmt.Sprintf("Error loading .env file: %v\n", err))
	}

	// Get keys from env file
	secretKey := os.Getenv("ENC_SECRET_KEY_HEX")
	iv := os.Getenv("ENC_IV_HEX")
	INSTANCE_NAME := os.Getenv("INSTANCE_NAME")

	fmt.Println("\n----------------------------------------------------------------")
	fmt.Println("\n--------------------------  Beta LMS  --------------------------")
	fmt.Println("\n----------------------------------------------------------------")
	fmt.Println("\n----------------  The Only Secure LMS You Need  ----------------")
	fmt.Println("\n----------------------------------------------------------------")
	fmt.Println("\n---------  Welcome to Beta LMS Video Processor! (v1)  ----------")
	fmt.Println("\n----------------------------------------------------------------")

	fmt.Printf("\nWelcome to Beta LMS Instance: %v\n\n", INSTANCE_NAME)

	// Get inputs from user
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the path of the video file: ")
	scanner.Scan()
	filePath := strings.ReplaceAll(scanner.Text(), "\"", "")

	fmt.Print("\nEnter the link of the YouTube video: ")
	scanner.Scan()
	youtubeLink := scanner.Text()
	encryptedPath := GetEncryptedFilePath(filePath, youtubeLink)

	if err := EncryptFile(filePath, encryptedPath, secretKey, iv); err != nil {
		panic(fmt.Sprintf("Error: %v\n", err))
	}

	fmt.Println("\nFile encrypted successfully.")
	fmt.Printf("\nEncrypted file path: %s\n", encryptedPath)
	fmt.Printf("\nEncrypted file name (already copied to clipboard): %s\n", encryptedPath[strings.LastIndex(encryptedPath, "\\")+1:strings.LastIndex(encryptedPath, ".")])

	if err := clipboard.Init(); err != nil {
		panic(fmt.Sprintf("Error initializing clipboard: %v\n", err))
	}
	clipboard.Write(clipboard.FmtText, []byte(encryptedPath[strings.LastIndex(encryptedPath, "\\")+1:strings.LastIndex(encryptedPath, ".")]))

	fmt.Println("\nPress Enter to exit.")
	scanner.Scan()
}
