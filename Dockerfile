# Etapa 1: Construção do binário
FROM golang:1.21 AS builder

# Configura o diretório de trabalho
WORKDIR /app

# Copia o código para o container
COPY . .

# Baixa dependências e compila a aplicação
RUN go mod tidy
RUN go build -o main ./cmd/api/main.go

# Etapa 2: Imagem para produção
FROM debian:bullseye-slim

# Configura o diretório de trabalho
WORKDIR /app

# Copia o binário compilado da etapa anterior
COPY --from=builder /app/main .

# Expõe a porta (ajuste conforme necessário)
EXPOSE 8080

# Comando para rodar a aplicação
CMD ["./main"]
