# Usar uma imagem base do Python
FROM python:3.10-slim

# Definir o diretório de trabalho na raiz do contexto (para pegar o produtos.proto)
WORKDIR /app

# Copiar TODO o projeto (igual ao modulo-b)
COPY . .

# Entrar na pasta do módulo específico
WORKDIR /app/modulo-p

# Instalar dependências
RUN pip install -r requirements.txt

# Expor a porta
EXPOSE 5000

# Comando para iniciar
CMD ["python", "app.py"]