# Dockerfile
# Usa una imagen base de Go para construir tu aplicación
FROM golang:1.25-alpine AS builder

# Establece el directorio de trabajo
WORKDIR /app

# Copia los archivos go.mod y go.sum para descargar dependencias
COPY go.mod .
COPY go.sum .

# Descarga todas las dependencias
RUN go mod download

# Copia el resto del código fuente
COPY . .

# Compila el ejecutable
RUN go build -o /admira-etl ./cmd/server

# Usa una imagen más ligera para el ejecutable final
FROM alpine:latest

# Crea un usuario no-root para seguridad
RUN adduser -D -g '' appuser
USER appuser

# Establece el directorio de trabajo
WORKDIR /app

# Copia el ejecutable compilado desde la etapa de builder
COPY --from=builder /admira-etl .

# Expone el puerto por defecto de la aplicación
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./admira-etl"]