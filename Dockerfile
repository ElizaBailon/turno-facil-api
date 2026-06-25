# --- Etapa 1: Compilación (Build Stage) ---
FROM golang:1.22-alpine AS builder

# Instalar dependencias necesarias para CGO (necesario para SQLite)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copiar los archivos de dependencias primero para aprovechar la caché de Docker
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código fuente
COPY . .

# Compilar el binario optimizado habilitando CGO para SQLite
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o turno-facil-api cmd/api/main.go

# --- Etapa 2: Ejecución (Run Stage) ---
FROM alpine:3.19

# Instalar librerías mínimas de ejecución y zona horaria
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=America/Guayaquil

WORKDIR /app

# Copiar el binario compilado desde la etapa anterior
COPY --from=builder /app/turno-facil-api .

# Exponer el puerto en el que corre tu API (8080)
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./turno-facil-api"]