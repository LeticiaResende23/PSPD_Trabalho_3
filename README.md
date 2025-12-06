
# Projeto de Microserviços com Kubernetes e Monitoramento

Este projeto implementa uma arquitetura de microserviços para um e-commerce simplificado, utilizando **Python (Flask)** para o Gateway e **Go (Golang)** para os serviços de backend (Catálogo e Inventário). A infraestrutura é orquestrada via **Kubernetes (Minikube)** e monitorada com **Prometheus**.

---

## Pré-requisitos

Certifique-se de ter as seguintes ferramentas instaladas:
* **Docker**
* **Minikube**
* **Kubectl**
* **Helm** (para instalar o Prometheus)
* **k6** (para testes de carga)

---

## Como Rodar o Projeto

### 1. Iniciar o Cluster (Multi-Node)
O projeto requer um cluster com 3 nós para testar a distribuição de carga.
```bash
minikube start --nodes 3
````

### 2\. Construir e Carregar as Imagens

Como o Minikube roda isolado, precisamos construir as imagens e enviá-las para dentro dele.

```bash
# Entre na pasta do código
cd app-principal

# Build das imagens
docker build -t modulo-p:latest -f modulo-p/Dockerfile .
docker build -t modulo-a:latest -f modulo-a/Dockerfile .
docker build -t modulo-b:latest -f modulo-b/Dockerfile .

# Enviar para o cluster (Pode demorar um pouco)
minikube image load modulo-p:latest
minikube image load modulo-a:latest
minikube image load modulo-b:latest
```

### 3\. Deploy da Aplicação

```bash
# Volte para a raiz e entre na pasta kubernete
cd ../kubernete
kubectl apply -f .
```

### 4\. Instalar o Monitoramento (Prometheus)

Usamos o Helm para instalar o Prometheus, desativando a persistência para evitar erros no Minikube.

```bash
helm repo add prometheus-community [https://prometheus-community.github.io/helm-charts](https://prometheus-community.github.io/helm-charts)
helm repo update

helm upgrade --install prometheus prometheus-community/prometheus \
  --set server.persistentVolume.enabled=false \
  --set alertmanager.persistentVolume.enabled=false
```

-----

## 🖥️ Acessando a Aplicação

Para acessar os serviços locais, utilize o `port-forward` em terminais separados.

**Terminal 1: Aplicação (Módulo P)**

```bash
kubectl port-forward deployment/modulo-p-deployment 5000:5000
```

  * Acesse: [http://localhost:5000/produtos/1](https://www.google.com/search?q=http://localhost:5000/produtos/1)

**Terminal 2: Dashboard do Prometheus**

```bash
kubectl port-forward deployment/prometheus-server 9090:9090
```

  * Acesse: [http://localhost:9090](https://www.google.com/search?q=http://localhost:9090)
  * Query para gráfico: `container_memory_usage_bytes{pod=~"modulo-.*"}`

-----

## Como Rodar o Teste de Carga (Elasticidade)

Para validar se o Kubernetes escala e distribui a carga corretamente:

1.  **Escalar a Aplicação (3 réplicas):**

    ```bash
    kubectl scale deployment modulo-a-deployment --replicas=3
    kubectl scale deployment modulo-b-deployment --replicas=3
    ```

2.  **Rodar o Teste com k6:**
    Certifique-se de que o **Terminal 1** (porta 5000) esteja aberto.

    ```bash
    # Na raiz do projeto
    k6 run k6Testes/teste_carga_caso_3.js
    ```

3.  **Monitorar:** Acompanhe o gráfico no Prometheus e veja o consumo subindo em todos os pods simultaneamente.
