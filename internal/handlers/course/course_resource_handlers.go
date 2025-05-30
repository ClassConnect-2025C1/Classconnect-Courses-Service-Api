package course

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"templateGo/internal/model"
	"templateGo/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateModule creates a new module for a course
func (h *courseHandlerImpl) CreateModule(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}
	_, ok = h.getCourseByID(c, courseID)
	if !ok {
		return
	}

	name := c.PostForm("name")
	if name == "" {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Validation Error", "Module name is required")
		return
	}

	module := &model.Module{
		CourseID: courseID,
		Name:     name,
	}

	if err := h.repo.CreateModule(module); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Failed to create module", "Error creating module: "+err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Module created successfully", "module_id": module.ID})
}

func (h *courseHandlerImpl) CreateResource(c *gin.Context) {
	moduleID, ok := h.getModuleID(c)
	if !ok {
		return
	}
	_, ok = h.getModuleByID(c, moduleID)
	if !ok {
		return
	}

	file, err := c.FormFile("file")
	link := c.PostForm("link")
	if err != nil && (link == "") {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Validation Error", "Either file or link must be provided")
		return
	}
	if err == nil && link != "" {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Validation Error", "Only one of file or link can be provided")
		return
	}
	var resourceId, url, resourceType string
	if file != nil {
		resourceId, url, err = h.postResource(file, moduleID)
		if err != nil {
			utils.NewErrorResponse(c, http.StatusInternalServerError, "Failed to upload file", "Error uploading file: "+err.Error())
			return
		}
	} else {
		resourceId = uuid.New().String()
		url = link
		resourceType = "link"
	}

	resource := &model.Resource{
		ID:       resourceId,
		ModuleID: moduleID,
		Type:     resourceType,
		URL:      url,
	}

	if err := h.repo.CreateResource(resource); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Failed to create resource", "Error creating resource: "+err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Resource created successfully"})
}

// Function to send POST request to Resource Service
func (h *courseHandlerImpl) postResource(fileHeader *multipart.FileHeader, courseID uint) (string, string, error) {
	fileContent, err := fileHeader.Open()
	if err != nil {
		return "", "", fmt.Errorf("error opening file: %w", err)
	}
	defer fileContent.Close()

	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, fileContent); err != nil {
		return "", "", fmt.Errorf("error reading file: %w", err)
	}

	client := &http.Client{}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if err := writer.WriteField("uploader_id", fmt.Sprintf("%d", courseID)); err != nil {
		return "", "", fmt.Errorf("error adding uploader_id to form: %w", err)
	}

	part, err := writer.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		return "", "", fmt.Errorf("error creating form file: %w", err)
	}
	if _, err = io.Copy(part, bytes.NewReader(buffer.Bytes())); err != nil {
		return "", "", fmt.Errorf("error copying file: %w", err)
	}
	err = writer.Close()
	if err != nil {
		return "", "", fmt.Errorf("error finalizing form: %w", err)
	}

	resourceServiceURL := os.Getenv("URL_RESOURCES")
	fmt.Printf("Resource Service URL: %s\n", resourceServiceURL)
	req, err := http.NewRequest("POST", resourceServiceURL+"/resource", body)
	if err != nil {
		return "", "", fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("error uploading file: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("error reading response: %w", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", "", fmt.Errorf("service returned error: %s", string(respBody))
	}
	var response struct {
		ResourceID string `json:"resource_id"`
		Link       string `json:"link"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return "", "", fmt.Errorf("error parsing response: %w", err)
	}
	return response.ResourceID, response.Link, nil
}
