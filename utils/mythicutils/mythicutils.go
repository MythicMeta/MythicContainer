package mythicutils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/MythicMeta/MythicContainer/config"
	"github.com/MythicMeta/MythicContainer/logging"
)

func SendFileToMythic(ctx context.Context, content *[]byte, fileID string, directFileToken string) error {
	if config.MythicConfig.MythicServerHost == "" {
		log.Fatalf("[-] Missing MYTHIC_SERVER_HOST environment variable point to mythic server IP")
	}
	if directFileToken == "" {
		return errors.New("missing direct file upload token")
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileWriter, err := writer.CreateFormFile("file", "payload")
	if err != nil {
		logging.LogError(err, "Failed to create new form file to upload payload")
		return err
	}
	_, err = io.Copy(fileWriter, bytes.NewReader(*content))
	if err != nil {
		logging.LogError(err, "Failed to write payload bytes to form")
		return err
	}
	fieldWriter, err := writer.CreateFormField("agent-file-id")
	if err != nil {
		logging.LogError(err, "Failed to add new form field to upload payload")
		return err
	}
	_, err = fieldWriter.Write([]byte(fileID))
	if err != nil {
		logging.LogError(err, "Failed to add in agent-file-id to form")
		return err
	}
	if err = writer.Close(); err != nil {
		logging.LogError(err, "Failed to finish multipart form for payload upload")
		return err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("http://%s:%d/direct/upload/%s",
			config.MythicConfig.MythicServerHost,
			config.MythicConfig.MythicServerPort,
			fileID), body)
	if err != nil {
		logging.LogError(err, "Failed to create new POST request to send payload to Mythic")
		return err
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", directFileToken))
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		logging.LogError(err, "Failed to send payload over to Mythic")
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		logging.LogError(nil, "Failed to send payload to Mythic", "status code", resp.StatusCode)
		return errors.New(fmt.Sprintf("Failed to send payload to Mythic with status code: %d\n", resp.StatusCode))
	}
	return nil
}

func GetFileFromMythic(ctx context.Context, fileID string, directFileToken string) (*[]byte, error) {
	if config.MythicConfig.MythicServerHost == "" {
		log.Fatalf("[-] Missing MYTHIC_SERVER_HOST environment variable point to mythic server IP")
	}
	if directFileToken == "" {
		return nil, errors.New("missing direct file download token")
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("http://%s:%d/direct/download/%s", config.MythicConfig.MythicServerHost,
			config.MythicConfig.MythicServerPort, fileID), nil)
	if err != nil {
		logging.LogError(err, "Failed to create new GET request to get file from Mythic")
		return nil, err
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", directFileToken))
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		logging.LogError(err, "Failed to send payload over to Mythic")
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		logging.LogError(nil, "Failed to read file from Mythic", "status code", resp.StatusCode)
		return nil, errors.New(fmt.Sprintf("Failed to read file from Mythic with status code: %d\n", resp.StatusCode))
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.LogError(err, "Failed to read body response from Mythic")
		return nil, errors.New(fmt.Sprintf("Failed to get file from Mythic: %s", err.Error()))
	}
	return &content, nil
}
