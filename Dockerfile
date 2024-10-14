# Fase de construção
FROM golang:latest AS builder

WORKDIR /app

# Copiar arquivos de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copiar o restante do código
COPY . .

# Compilar um binário estático
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /botany-back ./cmd/api/main.go

# Verifique se o binário foi realmente criado
RUN ls -l /botany-back

# Fase de execução
FROM alpine:latest

RUN apk --no-cache add ca-certificates

ENV GIN_MODE=release

# Copiar o binário do estágio de construção
COPY --from=builder /botany-back /botany-back

# Verifique se o binário está na imagem final
RUN ls -l /botany-back

EXPOSE ${API_PORT}

# Executar o binário
CMD ["/botany-back"]
