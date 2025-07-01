# Swagger Documentation

Este proyecto incluye documentación automática de la API usando Swagger (OpenAPI 3.0).

## Acceso a la Documentación

Una vez que el servidor esté ejecutándose, puedes acceder a la documentación de Swagger en:

```
http://localhost:8080/swagger/index.html
```

## Endpoints de Documentación

- **Swagger UI**: `http://localhost:8080/swagger/index.html` - Interfaz interactiva para explorar y probar la API
- **JSON**: `http://localhost:8080/swagger/doc.json` - Documentación en formato JSON
- **YAML**: `http://localhost:8080/swagger/swagger.yaml` - Documentación en formato YAML

## Funcionalidades de Swagger UI

La interfaz de Swagger te permite:

1. **Explorar todos los endpoints** de la API organizados por categorías
2. **Ver detalles completos** de cada endpoint incluyendo:
   - Parámetros requeridos y opcionales
   - Tipos de datos
   - Ejemplos de request/response
   - Códigos de estado HTTP
3. **Probar los endpoints directamente** desde la interfaz
4. **Autenticación**: Usar el botón "Authorize" para agregar tu token JWT

## Uso de Autenticación

Para endpoints que requieren autenticación:

1. Click en el botón "Authorize" en la esquina superior derecha
2. Ingresa tu token JWT en el formato: `Bearer tu_token_aqui`
3. Click "Authorize" y luego "Close"

Ahora puedes probar endpoints autenticados directamente desde la interfaz.

## Regenerar Documentación

Si modificas los comentarios de Swagger en el código, regenera la documentación con:

```bash
~/go/bin/swag init -g cmd/main.go -o docs
```

## Categorías de API

La API está organizada en las siguientes categorías:

- **health**: Endpoints de estado del servidor
- **courses**: Gestión de cursos (CRUD)
- **enrollment**: Matriculación en cursos
- **assignments**: Gestión de tareas
- **submissions**: Envío y calificación de tareas
- **feedback**: Sistema de retroalimentación
- **resources**: Gestión de recursos y módulos
- **statistics**: Estadísticas y métricas

## Ejemplos de Uso

### Crear un Curso
```json
POST /course
{
  "title": "Introduction to Programming",
  "description": "Learn the basics of programming with Python",
  "created_by": "teacher123",
  "capacity": 30,
  "eligibility_criteria": ["Computer Science Major"],
  "teaching_assistants": ["ta1@example.com"]
}
```

### Obtener Todos los Cursos
```
GET /courses
Authorization: Bearer your_jwt_token
```

La documentación completa con todos los endpoints y ejemplos está disponible en la interfaz de Swagger UI.
