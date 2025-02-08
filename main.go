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
	vodKeyId := os.Getenv("VOD_KEY_ID")
	vodKey := os.Getenv("VOD_KEY")
	INSTANCE_NAME := os.Getenv("INSTANCE_NAME")
	// var (
	// 	secretKey     = ""
	// 	iv            = ""
	// 	vodKeyId      = ""
	// 	vodKey        = ""
	// 	INSTANCE_NAME = ""
	// )

	fmt.Println("\n----------------------------------------------------------------")
	fmt.Println("\n--------------------------  Beta LMS  --------------------------")
	fmt.Println("\n----------------------------------------------------------------")
	fmt.Println("\n----------------  The Only Secure LMS You Need  ----------------")
	fmt.Println("\n----------------------------------------------------------------")
	fmt.Println("\n---------  Welcome to Beta LMS Video Processor! (v1)  ----------")
	fmt.Println("\n----------------------------------------------------------------")

	fmt.Printf("\nWelcome to Beta LMS Instance: %v\n\n", INSTANCE_NAME)

	continueFlag := ""
	scanner := bufio.NewScanner(os.Stdin)

	for continueFlag == "" {
		// Get inputs from user
		fmt.Print("Enter the path of the video file: ")
		scanner.Scan()
		filePath := strings.ReplaceAll(scanner.Text(), "\"", "")

		fmt.Print("Enter video encoded name (press enter directly to generate random name): ")
		scanner.Scan()
		videoName := scanner.Text()
		encryptedPath := GetEncryptedFilePath(filePath, videoName)

		if err := PackageVideoFile(filePath, encryptedPath, vodKeyId, vodKey, iv); err != nil {
			panic(fmt.Sprintf("Error: %v\n", err))
		}

		if err := EncryptFile(filePath, encryptedPath+".enc", secretKey, iv); err != nil {
			panic(fmt.Sprintf("Error: %v\n", err))
		}

		fmt.Println("\nFile encrypted successfully.")
		fmt.Printf("\nEncrypted file path: %s\n", encryptedPath)
		fmt.Printf("\nEncrypted file name (already copied to clipboard): %s\n", encryptedPath[strings.LastIndex(encryptedPath, "\\")+1:])
		if encFileHash, err := ComputeSHA256Hash(encryptedPath + ".enc"); err != nil {
			panic(fmt.Sprintf("Error: %v\n", err))
		} else {
			fmt.Printf("\nEncrypted file hash: %s\n", encFileHash)
		}

		if err := clipboard.Init(); err != nil {
			panic(fmt.Sprintf("Error initializing clipboard: %v\n", err))
		}
		clipboard.Write(clipboard.FmtText, []byte(encryptedPath[strings.LastIndex(encryptedPath, "\\")+1:]))

		fmt.Print("\nProcess another video (y: enter, n: type anything) ? ")
		scanner.Scan()
		continueFlag = scanner.Text()
	}
}
