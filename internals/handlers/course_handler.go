package handlers

import (
	"net/http"
	"strconv"
	"templateGo/internals/models"
	"templateGo/internals/services"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	service *services.CourseService
}

func NewCourseHandler(service *services.CourseService) *CourseHandler {
	return &CourseHandler{service}
}

// Crear curso
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateCourse(&course); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el curso"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Curso creado con éxito", "data": course})
}

// Obtener todos los cursos
func (h *CourseHandler) GetAllCourses(c *gin.Context) {
	courses, err := h.service.GetAllCourses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener cursos"})
		return
	}
	// Retornar lista vacía en vez de null
	if courses == nil {
		courses = []models.Course{}
	}
	c.JSON(http.StatusOK, gin.H{"data": courses})
}

// Obtener curso por ID
func (h *CourseHandler) GetCourseByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	course, err := h.service.GetCourseByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Curso no encontrado"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": course})
}

// Editar curso
func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course.ID = uint(id)
	if err := h.service.UpdateCourse(&course); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el curso"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Curso actualizado con éxito", "data": course})
}

// Eliminar curso (lógicamente)
func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.service.DeleteCourse(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el curso"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Curso eliminado con éxito"})
}
