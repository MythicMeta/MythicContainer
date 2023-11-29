package mythicutils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/MythicMeta/MythicContainer/config"
	"github.com/MythicMeta/MythicContainer/logging"
)

func init() {
	if config.MythicConfig.MythicServerHost == "" {
		log.Fatalf("[-] Missing MYTHIC_SERVER_HOST environment variable point to mythic server IP")
	}
}

func SendFileToMythic(content *[]byte, fileID string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if fileWriter, err := writer.CreateFormFile("file", "payload"); err != nil {
		logging.LogError(err, "Failed to create new form file to upload payload")
		return err
	} else if _, err = io.Copy(fileWriter, bytes.NewReader(*content)); err != nil {
		logging.LogError(err, "Failed to write payload bytes to form")
		return err
	} else if fieldWriter, err := writer.CreateFormField("agent-file-id"); err != nil {
		logging.LogError(err, "Failed to add new form field to upload payload")
		return err
	} else if _, err := fieldWriter.Write([]byte(fileID)); err != nil {
		logging.LogError(err, "Failed to add in agent-file-id to form")
		return err
	}
	writer.Close()
	if request, err := http.NewRequest("POST",
		fmt.Sprintf("http://%s:%d/direct/upload/%s", config.MythicConfig.MythicServerHost,
			config.MythicConfig.MythicServerPort, fileID), body); err != nil {
		logging.LogError(err, "Failed to create new POST request to send payload to Mythic")
		return err
	} else {
		request.Header.Add("Content-Type", writer.FormDataContentType())
		if resp, err := http.DefaultClient.Do(request); err != nil {
			logging.LogError(err, "Failed to send payload over to Mythic")
			return err
		} else {
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				logging.LogError(nil, "Failed to send payload to Mythic", "status code", resp.StatusCode)
				return errors.New(fmt.Sprintf("Failed to send payload to Mythic with status code: %d\n", resp.StatusCode))
			}
		}
	}

	return nil
}

func GetFileFromMythic(fileID string) (*[]byte, error) {
	if request, err := http.NewRequest("GET",
		fmt.Sprintf("http://%s:%d/direct/download/%s", config.MythicConfig.MythicServerHost,
			config.MythicConfig.MythicServerPort, fileID), nil); err != nil {
		logging.LogError(err, "Failed to create new GET request to get file from Mythic")
		return nil, err
	} else {
		if resp, err := http.DefaultClient.Do(request); err != nil {
			logging.LogError(err, "Failed to send payload over to Mythic")
			return nil, err
		} else {
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				logging.LogError(nil, "Failed to read file from Mythic", "status code", resp.StatusCode)
				return nil, errors.New(fmt.Sprintf("Failed to read file from Mythic with status code: %d\n", resp.StatusCode))
			} else if content, err := io.ReadAll(resp.Body); err != nil {
				logging.LogError(err, "Failed to read body response from Mythic")
				return nil, errors.New(fmt.Sprintf("Failed to get file from Mythic: %s", err.Error()))
			} else {
				return &content, nil
			}
		}
	}
}
