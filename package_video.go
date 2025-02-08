package main

import (
	"fmt"
	"os"
	"os/exec"
)

func PackageVideoFile(filePath string, encryptedPath string, vodKeyId string, vodKey string, iv string) error {
	packagerPath, err := GetPackagerExecutable()
	defer os.Remove(packagerPath)

	if err != nil {
		fmt.Println("Error getting Shaka Packager:", err)
		os.Exit(1)
	}

	// Construct the packager command
	cmd := exec.Command(packagerPath,
		fmt.Sprintf("in=%s,stream=video,output=%s,input_format=mp4,output_format=mp4,drm_label=HD", filePath, encryptedPath+"_v.mp4"),
		fmt.Sprintf("in=%s,stream=audio,output=%s,input_format=mp4,output_format=mp4,drm_label=AUDIO", filePath, encryptedPath+"_a.mp4"),
		"--enable_raw_key_encryption",
		"--keys", fmt.Sprintf("key_id=%s:key=%s", vodKeyId, vodKey),
		"--iv", iv,
		"--protection_scheme", "cenc",
		"--clear_lead", "0",
		fmt.Sprintf("--mpd_output=%s", encryptedPath+"_m.mpd"),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Print("\n")
	if err := cmd.Run(); err != nil {
		fmt.Println("Error running packager:", err)
		os.Exit(1)
	}

	return nil
}
