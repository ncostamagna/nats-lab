# Etapa de build
FROM golang:1.24 AS builder

# Configura el directorio de trabajo
WORKDIR /app

# Copia los archivos del proyecto
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compila el binario de manera estática
RUN CGO_ENABLED=0 GOOS=linux go build -o main

# Etapa final: imagen mínima
FROM scratch

# (Opcional) Copiar certificados si tu app hace HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copiar el binario compilado
COPY --from=builder /app/main /main

# Define el punto de entrada
ENTRYPOINT ["/main"]