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
		resourceType = "file"
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

// getResources retrieves all resources from a course as modules with their resources
func (h *courseHandlerImpl) GetResources(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}
	_, ok = h.getCourseByID(c, courseID)
	if !ok {
		return
	}
	modules, err := h.repo.GetModulesByCourseID(courseID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve modules", "Error retrieving modules: "+err.Error())
		return
	}

	resources := make([]gin.H, 0, len(modules))
	for _, module := range modules {
		moduleResources, err := h.repo.GetResourcesByModuleID(module.ID)
		if err != nil {
			utils.NewErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve resources", "Error retrieving resources: "+err.Error())
			return
		}
		resources = append(resources, gin.H{
			"module_id":   module.ID,
			"order":       module.Order,
			"module_name": module.Name,
			"resources":   moduleResources,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"modules": resources,
	})
}

func (h *courseHandlerImpl) DeleteResource(c *gin.Context) {
	resourceID := c.Param("resource_id")
	if resourceID == "" {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Validation Error", "Resource ID is required")
		return
	}

	if err := h.repo.DeleteResource(resourceID); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Failed to delete resource", "Error deleting resource: "+err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Resource deleted successfully"})
}

func (h *courseHandlerImpl) DeleteModule(c *gin.Context) {
	moduleID, ok := h.getModuleID(c)
	if !ok {
		return
	}
	_, ok = h.getModuleByID(c, moduleID)
	if !ok {
		return
	}

	if err := h.repo.DeleteModule(moduleID); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Failed to delete module", "Error deleting module: "+err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Module deleted successfully"})
}

// PatchResources updates the order of resources in a module
func (h *courseHandlerImpl) PatchResources(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		// getCourseID already sends an error response if needed
		return
	}

	// Ensure the course itself exists
	_, ok = h.getCourseByID(c, courseID)
	if !ok {
		// getCourseByID already sends an error response if needed
		return
	}

	var req model.CourseOrderUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid request payload", "Error binding JSON: "+err.Error())
		return
	}

	for moduleIndex, moduleUpdate := range req.Modules {
		module, moduleOk := h.getModuleByID(c, moduleUpdate.ModuleID)
		if !moduleOk {
			// Not handling rollback
			return
		}
		if module.CourseID != courseID {
			utils.NewErrorResponse(c, http.StatusForbidden, "Module does not belong to this course",
				fmt.Sprintf("Module ID %d is not part of course ID %d", moduleUpdate.ModuleID, courseID))
			// Not handling rollback
			return
		}

		if err := h.repo.UpdateModuleOrder(moduleUpdate.ModuleID, moduleIndex); err != nil {
			utils.NewErrorResponse(c, http.StatusInternalServerError, "Failed to update module order",
				fmt.Sprintf("Error updating order for module ID %d: %s", moduleUpdate.ModuleID, err.Error()))
			// Not handling rollback
			return
		}

		for resourceIndex, resourceUpdate := range moduleUpdate.Resources {
			resource, err := h.repo.GetResourceByID(resourceUpdate.ID)
			if err != nil {
				utils.NewErrorResponse(c, http.StatusNotFound, "Resource not found",
					fmt.Sprintf("Error finding resource ID %s: %s", resourceUpdate.ID, err.Error()))
				// Not handling rollback
				return
			}
			if resource.ModuleID != moduleUpdate.ModuleID {
				utils.NewErrorResponse(c, http.StatusForbidden, "Resource does not belong to this module",
					fmt.Sprintf("Resource ID %s is not part of module ID %d", resourceUpdate.ID, moduleUpdate.ModuleID))
				// Not handling rollback
				return
			}

			if err := h.repo.UpdateResourceOrder(resourceUpdate.ID, resourceIndex); err != nil {
				utils.NewErrorResponse(c, http.StatusInternalServerError, "Failed to update resource order",
					fmt.Sprintf("Error updating order for resource ID %s: %s", resourceUpdate.ID, err.Error()))
				// Not handling rollback
				return
			}
		}
	}

	// If using transactions, this is where you would commit.
	// See defer block above for commit/rollback logic.

	c.JSON(http.StatusOK, gin.H{"message": "Modules and resources order updated successfully"})
}
