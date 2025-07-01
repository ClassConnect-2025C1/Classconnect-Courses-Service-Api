# Configuración de Swagger para Render

## Variables de Entorno en Render

Para habilitar Swagger en producción, agrega estas variables de entorno en tu servicio de Render:

### Variables requeridas:
```bash
# Para habilitar Swagger en producción
ENABLE_SWAGGER=true

# Modo de Gin (opcional, por defecto será 'release' en producción)
GIN_MODE=release
```

## URLs de acceso:

### En desarrollo (local):
```bash
# API Health Check
http://localhost:8080/

# Swagger UI
http://localhost:8080/swagger/index.html

# Swagger JSON
http://localhost:8080/swagger/doc.json
```

### En producción (Render):
```bash
# API Health Check
https://classconnect-courses-service-api.onrender.com/

# Swagger UI
https://classconnect-courses-service-api.onrender.com/swagger/index.html

# Swagger JSON
https://classconnect-courses-service-api.onrender.com/swagger/doc.json
```

**✨ Configuración inteligente:** Swagger detecta automáticamente la URL base del entorno donde se está ejecutando, por lo que las peticiones se harán siempre a la URL correcta (local o producción) sin necesidad de configuración adicional.

