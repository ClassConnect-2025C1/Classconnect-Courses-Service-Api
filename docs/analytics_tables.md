# Analytics Tables Documentation

## Overview
Este proyecto incluye dos tablas de analytics completamente diferentes que sirven propósitos distintos:

## 📊 CourseAnalytics (curso | estadística)
**Propósito**: Almacena estadísticas **GENERALES DEL CURSO** como un todo.

### Qué contiene:
- Métricas agregadas de todos los estudiantes
- Estadísticas del curso completo
- Información del desempeño general
- Datos de popularidad y engagement

### Ejemplos de datos que podrías almacenar:
```json
{
  "total_enrolled_students": 150,
  "total_active_students": 142,
  "course_average_grade": 85.5,
  "course_completion_rate": 78.3,
  "total_assignments_created": 12,
  "most_difficult_assignment": "Final Project",
  "most_popular_resource": "Chapter 5 Video",
  "instructor_engagement_rate": 92.1,
  "last_course_activity": "2025-06-25T10:30:00Z"
}
```

## 👤 UserCourseAnalytics (user | curso | estadística)
**Propósito**: Almacena estadísticas **INDIVIDUALES DE CADA USUARIO** inscrito en un curso específico.

### Qué contiene:
- Métricas específicas de UN usuario en particular
- Progreso individual del estudiante
- Comportamiento personal de estudio
- Desempeño individual comparado con otros

### Ejemplos de datos que podrías almacenar:
```json
{
  "user_personal_grade": 92.5,
  "user_assignments_completed": 8,
  "user_total_assignments": 10,
  "user_time_spent_minutes": 450,
  "user_last_activity": "2025-06-25T09:15:00Z",
  "user_performance_trend": "improving",
  "user_rank_in_course": 15,
  "user_favorite_resources": ["Video Lecture 3", "PDF Chapter 2"],
  "user_study_pattern": "evening"
}
```

## 🔍 Diferencias Clave

| Aspecto | CourseAnalytics | UserCourseAnalytics |
|---------|----------------|-------------------|
| **Scope** | Todo el curso | Un usuario específico |
| **Ejemplo de Promedio** | Promedio de TODOS los estudiantes | Promedio personal del usuario |
| **Registros por Curso** | 1 registro por curso | 1 registro por cada usuario inscrito |
| **Uso típico** | Dashboard del instructor | Dashboard del estudiante |
| **Actualización** | Cuando cambian métricas del curso | Cuando el usuario realiza actividades |

## 💡 Casos de Uso

### CourseAnalytics:
- Dashboard del instructor para ver el desempeño general
- Reportes administrativos
- Análisis de efectividad del curso
- Comparación entre diferentes cursos

### UserCourseAnalytics:
- Dashboard personal del estudiante
- Seguimiento individual de progreso
- Recomendaciones personalizadas
- Gamificación (rankings, logros)

## 🚀 Implementación

Ambas tablas usan JSON como string para máxima flexibilidad. Puedes almacenar cualquier estructura de datos que necesites sin modificar la base de datos.

### Consultas Típicas:

```go
// Obtener estadísticas generales del curso
var courseStats model.CourseAnalytics
db.Where("course_id = ?", courseID).First(&courseStats)

// Obtener estadísticas de un usuario específico en un curso
var userStats model.UserCourseAnalytics  
db.Where("user_id = ? AND course_id = ?", userID, courseID).First(&userStats)

// Obtener todas las estadísticas de usuarios en un curso (para rankings)
var allUserStats []model.UserCourseAnalytics
db.Where("course_id = ?", courseID).Find(&allUserStats)
```
