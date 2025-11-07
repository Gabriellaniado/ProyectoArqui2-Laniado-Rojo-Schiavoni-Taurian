# Usar una imagen base de Go
FROM golang:1.24-alpine

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar los archivos de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código fuente
COPY . .

# Compilar la aplicación
RUN go build -o /users-api cmd/api/main.go

# Exponer el puerto que usará la aplicación
EXPOSE 8082

# Comando para correr la aplicación
CMD [ "/users-api" ]