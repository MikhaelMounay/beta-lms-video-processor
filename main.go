package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.design/x/clipboard"
)

var (
	SecretKey     string
	Iv            string
	VodKeyId      string
	VodKey        string
	INSTANCE_NAME string
	BaseB2Folder  string
	appKeyID      string
	appKey        string
	bucketID      string
	UploadURL     string
	AuthToken     string
)

func main() {
	// Load environment variables from the .env file
	if err := godotenv.Load(); err != nil {
		panic(fmt.Sprintf("Error loading .env file: %v\n", err))
	}

	// Get keys from env file
	SecretKey = os.Getenv("ENC_SECRET_KEY_HEX")
	Iv = os.Getenv("ENC_IV_HEX")
	VodKeyId = os.Getenv("VOD_KEY_ID")
	VodKey = os.Getenv("VOD_KEY")
	INSTANCE_NAME = os.Getenv("INSTANCE_NAME")
	BaseB2Folder = os.Getenv("BASE_B2_FOLDER")
	appKeyID = os.Getenv("B2_APP_KEY_ID")
	appKey = os.Getenv("B2_APP_KEY")
	bucketID = os.Getenv("B2_BUCKET_ID")

	// SecretKey = ""
	// Iv = ""
	// VodKeyId = ""
	// VodKey = ""
	// INSTANCE_NAME = ""
	// BaseB2Folder = ""
	// appKeyID = ""
	// appKey = ""
	// bucketID = ""

	fmt.Println("\n----------------------------------------------------------------")
	fmt.Println("\n--------------------------  Beta LMS  --------------------------")
	fmt.Println("\n----------------------------------------------------------------")
	fmt.Println("\n----------------  The Only Secure LMS You Need  ----------------")
	fmt.Println("\n----------------------------------------------------------------")
	fmt.Println("\n---------  Welcome to Beta LMS Video Processor! (v2)  ----------")
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

		if err := GetB2UploadURL(appKeyID, appKey, bucketID); err != nil {
			panic(fmt.Sprintf("Error: %v\n", err))
		}
		ProcessAndUploadVideo(filePath, videoName)

		fmt.Print("\nProcess another video (y: enter, n: type anything) ? ")
		scanner.Scan()
		continueFlag = scanner.Text()
	}
}

func ProcessAndUploadVideo(filePath string, videoName string) {
	encryptedPath := GetEncryptedFilePath(filePath, videoName)

	if err := PackageVideoFile(filePath, encryptedPath, VodKeyId, VodKey, Iv); err != nil {
		panic(fmt.Sprintf("Error: %v\n", err))
	}

	if err := EncryptFile(filePath, encryptedPath+".enc", SecretKey, Iv); err != nil {
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

	fmt.Println("\nUploading Video ...")

	// Sequentially upload files
	filesToUpload := []string{
		encryptedPath + ".enc",
		encryptedPath + "_v.mp4",
		encryptedPath + "_a.mp4",
		encryptedPath + "_m.mpd",
	}

	for _, file := range filesToUpload {
		if err := UploadFile(file); err != nil {
			fmt.Printf("\nUpload failed: %s\nError: %v\n", file, err)
		}
	}
}
