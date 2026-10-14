package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/HiogoPariz/notez/apps/api/internal/config"
	"github.com/HiogoPariz/notez/apps/api/internal/dto"
)

type FileIntegration struct {
	dto *dto.NoteDTO
}

func CreateFileIntegration(dto *dto.NoteDTO) *FileIntegration {
	return &FileIntegration{dto}
}

func (integration *FileIntegration) baseURL() string {
	return config.Get("FILES_URL", "http://localhost:3001")
}

func (integration *FileIntegration) GetFileContent() (string, error) {
	client := http.Client{}
	requestURL := fmt.Sprintf("%s/%s", integration.baseURL(), integration.dto.FileName)

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("files service returned %d for %s", res.StatusCode, integration.dto.FileName)
	}

	response, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	// O serviço de files guarda o conteúdo no formato {"content": "..."}.
	// Extrai apenas o conteúdo; se não estiver nesse formato, devolve o corpo cru.
	var payload struct {
		Content string `json:"content"`
	}
	if err := json.Unmarshal(response, &payload); err == nil {
		return payload.Content, nil
	}

	return string(response), nil
}

func (integration *FileIntegration) CreateFileContent(content string, file_name string) error {
	client := http.Client{}
	requestURL := fmt.Sprintf("%s/%s", integration.baseURL(), file_name)
	requestBody := []byte(fmt.Sprintf(`{"content": %q}`, content))

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(requestBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("files service returned %d", res.StatusCode)
	}

	return nil
}
