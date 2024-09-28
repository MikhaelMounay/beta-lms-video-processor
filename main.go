package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"io"
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
	encryptedPath := getEncryptedFileName(filePath, youtubeLink)

	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		return
	}

	// get keys from env file
	secretKey := os.Getenv("ENC_SECRET_KEY_HEX")
	iv := os.Getenv("ENC_IV_HEX")

	if err := encryptFile(filePath, encryptedPath, secretKey, iv); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func getEncryptedFileName(filePath string, videoUrl string) string {
	id := (strings.Split(videoUrl[strings.Index(videoUrl, "v=")+2:len(videoUrl)-1], ""))

	for i, j := 0, len(id)-1; i < j; i, j = i+1, j-1 {
		id[i], id[j] = id[j], id[i]
	}

	videoName := strings.Join(id, "")

	finalPath := strings.Replace(filePath, filePath[strings.LastIndex(filePath, "\\"):], fmt.Sprintf("\\%s.enc", videoName), 1)

	return finalPath
}

func encryptFile(filePath string, encryptedPath string, secretKey string, iv string) error {
	// Convert the secret key and IV from hex strings to byte slices
	key, err := hex.DecodeString(secretKey)
	if err != nil {
		return fmt.Errorf("failed to decode secret key: %v", err)
	}

	ivBytes, err := hex.DecodeString(iv)
	if err != nil {
		return fmt.Errorf("failed to decode IV: %v", err)
	}

	// Create a new AES cipher block and CBC encrypter
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("failed to create AES cipher: %v", err)
	}

	cipherStream := cipher.NewCBCEncrypter(block, ivBytes)

	// Open the input file
	inputFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %v", err)
	}
	defer inputFile.Close()

	// Open the output file
	outputFile, err := os.Create(encryptedPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	// Create buffers for reading from the input file and writing to the output
	buffer := make([]byte, aes.BlockSize*1024)
	for {
		bytesRead, readErr := inputFile.Read(buffer)
		if bytesRead > 0 {
			data := buffer[:bytesRead]

			// Apply padding if it's the last read
			if bytesRead < len(buffer) || readErr == io.EOF {
				data = pkcs7Padding(data, aes.BlockSize)
			}

			cipherStream.CryptBlocks(data, data)
			if _, writeErr := outputFile.Write(data); writeErr != nil {
				return fmt.Errorf("failed to write to output file: %v", writeErr)
			}
		}

		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return fmt.Errorf("failed to read input file: %v", readErr)
		}
	}

	fmt.Println("File encrypted successfully.")
	return nil
}

func pkcs7Padding(data []byte, blockSize int) []byte {
	paddingSize := blockSize - len(data)%blockSize
	padding := bytes.Repeat([]byte{byte(paddingSize)}, paddingSize)
	return append(data, padding...)
}
