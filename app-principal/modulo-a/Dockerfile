# Etapa 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copiar todo o contexto (igual ao modulo-b) para garantir acesso a tudo
COPY . .

# Mudar para o diretório do modulo-a
WORKDIR /app/modulo-a

# Baixar as dependências e compilar
RUN go mod download
RUN CGO_ENABLED=0 go build -o /main .

# Etapa 2: Imagem final
FROM alpine:latest

WORKDIR /app

# Copiar o binário construído da etapa anterior
COPY --from=builder /main .

# Expor a porta
EXPOSE 50051

# Comando para rodar
CMD ["./main"]