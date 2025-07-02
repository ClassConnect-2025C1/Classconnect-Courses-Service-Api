# Analytics Tables Documentation

## Overview
Este proyecto incluye **TRES tablas** de analytics/estadísticas completamente diferentes que sirven propósitos distintos, además de endpoints de estadísticas que proporcionan información agregada y procesada:

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

## 🌐 GlobalStatistics (profesor | estadísticas globales)
**Propósito**: Almacena estadísticas **GLOBALES AGREGADAS** de todos los cursos de un profesor específico.

### Qué contiene:
- Estadísticas consolidadas de TODOS los cursos del profesor
- Promedios generales calculados
- Métricas de rendimiento global del instructor
- Datos para dashboard principal del profesor

### Estructura de la tabla:
```go
type GlobalStatistics struct {
    ID                   uint    `json:"id" gorm:"primaryKey"`
    TeacherEmail         string  `json:"teacher_email" gorm:"uniqueIndex;not null"`
    GlobalAverageGrade   float64 `json:"global_average_grade" gorm:"default:0"`
    GlobalSubmissionRate float64 `json:"global_submission_rate" gorm:"default:0"`
}
```

### Campos:
- `teacher_email`: Email único del profesor (índice único)
- `global_average_grade`: Promedio de calificaciones de TODOS los cursos
- `global_submission_rate`: Tasa de entrega promedio de TODOS los cursos

## 🔍 Diferencias Clave

| Aspecto | CourseAnalytics | UserCourseAnalytics | GlobalStatistics |
|---------|----------------|-------------------|------------------|
| **Scope** | Todo el curso | Un usuario específico | Todos los cursos del profesor |
| **Ejemplo de Promedio** | Promedio de TODOS los estudiantes del curso | Promedio personal del usuario | Promedio de TODOS los cursos |
| **Registros por Curso** | 1 registro por curso | 1 registro por cada usuario inscrito | 1 registro por profesor |
| **Uso típico** | Dashboard del instructor (curso específico) | Dashboard del estudiante | Dashboard principal del profesor |
| **Actualización** | Cuando cambian métricas del curso | Cuando el usuario realiza actividades | Cuando se actualizan estadísticas globales |
| **Tipo de datos** | JSON flexible | JSON flexible | Campos estructurados |

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

### GlobalStatistics:
- Dashboard principal del profesor
- Vista general de rendimiento de todos los cursos
- Reportes ejecutivos y análisis global
- Comparación de rendimiento histórico general

## 🚀 Implementación

Las tablas **CourseAnalytics** y **UserCourseAnalytics** usan JSON como string para máxima flexibilidad. La tabla **GlobalStatistics** usa campos estructurados para eficiencia en consultas frecuentes.

### Consultas Típicas:

```go
// Obtener estadísticas generales del curso (JSON flexible)
var courseStats model.CourseAnalytics
db.Where("course_id = ?", courseID).First(&courseStats)

// Obtener estadísticas de un usuario específico en un curso (JSON flexible)
var userStats model.UserCourseAnalytics  
db.Where("user_id = ? AND course_id = ?", userID, courseID).First(&userStats)

// Obtener estadísticas globales del profesor (campos estructurados)
var globalStats model.GlobalStatistics
db.Where("teacher_email = ?", teacherEmail).First(&globalStats)

// Obtener todas las estadísticas de usuarios en un curso (para rankings)
var allUserStats []model.UserCourseAnalytics
db.Where("course_id = ?", courseID).Find(&allUserStats)
```

## 📊 Endpoints de Estadísticas

### 🌐 GET `/statistics/global` - Estadísticas Globales del Profesor
**Propósito**: Obtiene estadísticas **GLOBALES** promediadas de todos los cursos del profesor.

**Qué devuelve**:
```json
{
  "statistics": {
    "id": 1,
    "teacher_email": "profesor@ejemplo.com",
    "global_average_grade": 85.7,
    "global_submission_rate": 78.3
  }
}
```

**Campos del response**:
- `teacher_email`: Email del profesor
- `global_average_grade`: Promedio de calificaciones de todos los cursos
- `global_submission_rate`: Tasa de entrega promedio de todos los cursos

### 🎯 GET `/statistics/{course_id}` - Estadísticas de Curso Específico
**Propósito**: Obtiene estadísticas **DETALLADAS** de un curso específico.

**Qué devuelve**:
```json
{
  "statistics": {
    "id": 1,
    "course_id": 123,
    "course_name": "Introducción a la Programación",
    "global_average_grade": 87.5,
    "global_submission_rate": 82.1,
    "last_10_assignments_average_grade_tendency": "crescent",
    "last_10_assignments_submission_rate_tendency": "stable",
    "suggestions": "El curso muestra una tendencia positiva en las calificaciones",
    "statistics_for_assignments": [
      {
        "date": "2025-06-01T10:00:00Z",
        "average_grade": 85.0,
        "submission_rate": 90.0
      },
      {
        "date": "2025-06-15T10:00:00Z", 
        "average_grade": 88.5,
        "submission_rate": 85.0
      }
    ]
  }
}
```

**Campos del response**:
- `course_id`: ID del curso
- `course_name`: Nombre del curso
- `global_average_grade`: Promedio general de calificaciones del curso
- `global_submission_rate`: Tasa de entrega general del curso
- `last_10_assignments_average_grade_tendency`: Tendencia de las últimas 10 tareas ("crescent", "decrescent", "stable")
- `last_10_assignments_submission_rate_tendency`: Tendencia de entrega de las últimas 10 tareas
- `suggestions`: Sugerencias generadas por IA basadas en las estadísticas
- `statistics_for_assignments`: Array con estadísticas históricas por tarea

### 👤 GET `/statistics/course/{course_id}/user/{user_id}` - Estadísticas de Usuario
**Propósito**: Obtiene estadísticas **INDIVIDUALES** de un usuario específico en un curso.

**Qué devuelve**:
```json
{
  "statistics": {
    "id": 1,
    "course_id": 123,
    "user_id": "user123",
    "average_grade": 92.5,
    "submission_rate": 95.0,
    "last_10_assignments_average_grade_tendency": "crescent",
    "last_10_assignments_submission_rate_tendency": "stable",
    "statistics_for_assignments": [
      {
        "date": "2025-06-01T10:00:00Z",
        "average_grade": 90.0,
        "submission_rate": 100.0
      },
      {
        "date": "2025-06-15T10:00:00Z",
        "average_grade": 95.0,
        "submission_rate": 90.0
      }
    ]
  }
}
```

**Campos del response**:
- `course_id`: ID del curso
- `user_id`: ID del usuario
- `average_grade`: Promedio personal de calificaciones del usuario
- `submission_rate`: Tasa de entrega personal del usuario
- `last_10_assignments_average_grade_tendency`: Tendencia personal de las últimas 10 tareas
- `last_10_assignments_submission_rate_tendency`: Tendencia personal de entrega
- `statistics_for_assignments`: Array con estadísticas históricas personales por tarea

## 💾 Relación con las Tablas Analytics

### CourseAnalytics vs Endpoints de Estadísticas:
- **CourseAnalytics**: Almacena datos en **JSON flexible** para analytics complejos
- **Statistics endpoints**: Devuelven datos **estructurados y procesados** para dashboards

### UserCourseAnalytics vs User Statistics:
- **UserCourseAnalytics**: Almacena **datos personalizados** en JSON
- **User Statistics endpoint**: Devuelve **métricas calculadas** en tiempo real

### GlobalStatistics vs Global Statistics Endpoint:
- **GlobalStatistics**: Almacena **datos agregados precalculados** en campos estructurados
- **Global Statistics endpoint**: Devuelve directamente los datos de la tabla GlobalStatistics

## 📋 Resumen de las 3 Tablas

| Tabla | Tipo de Datos | Propósito | Endpoint Relacionado |
|-------|---------------|-----------|---------------------|
| **CourseAnalytics** | JSON flexible | Analytics complejos del curso | N/A (almacenamiento) |
| **UserCourseAnalytics** | JSON flexible | Analytics personalizados del usuario | N/A (almacenamiento) |
| **GlobalStatistics** | Campos estructurados | Estadísticas globales del profesor | `GET /statistics/global` |

### Flujo de Datos:
1. **Datos Raw** → Se procesan y almacenan en `CourseAnalytics` y `UserCourseAnalytics`
2. **Agregación** → Se calculan y almacenan en `GlobalStatistics`
3. **API Response** → Los endpoints devuelven datos procesados y estructurados
