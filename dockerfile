# Imagen base de Go
FROM golang:latest

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar los archivos del directorio actual al contenedor en /app
COPY . .

# Exponer el puerto que utiliza la aplicación
EXPOSE 4008

# Comando para ejecutar la aplicación
CMD ["go", "run", "main.go"]
