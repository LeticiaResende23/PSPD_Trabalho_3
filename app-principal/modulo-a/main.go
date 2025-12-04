package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"encoding/json"
    "net/http"
	"time"
	"os"
	"strconv"
	"google.golang.org/grpc"
	pb "projeto-catalogo/proto" 
)

type servidor struct {
	pb.UnimplementedCatalogoServer 
}

// Dados mocados (simulando banco de dados)
var produtos = map[string]*pb.InfoBasicaResponse{
	"1": {Id: "1", Nome: "Mouse Gamer Pro", Descricao: "Mouse óptico com 16.000 DPI"},
	"2": {Id: "2", Nome: "Teclado Mecânico RGB", Descricao: "Teclado com switches blue e iluminação customizável"},
}

// Função auxiliar para ler variáveis de ambiente com valor padrão
func getEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Erro ao converter variável %s: %v. Usando padrão: %d", key, err, defaultValue)
		return defaultValue
	}
	return value
}

func (s *servidor) GetInfoBasica(ctx context.Context, req *pb.ProdutoRequest) (*pb.InfoBasicaResponse, error) {
	idProduto := req.GetId()
	
	// --- CONFIGURAÇÃO DINÂMICA VIA VARIÁVEIS DE AMBIENTE ---
	atrasoMs := getEnvInt("ATRASO_MS", 100)      // Padrão: 100ms se não definir nada
	cargaCpu := getEnvInt("CARGA_CPU", 100000)   // Padrão: 100.000 loops

	log.Printf("Requisição recebida ID: %v | Config: Delay=%dms, CPU=%d loops", idProduto, atrasoMs, cargaCpu)

	// 1. Simulação de Carga de CPU
	start := time.Now()
	for i := 0; i < cargaCpu; i++ {
		_ = float64(i) * rand.Float64() // Cálculo matemático inútil para gastar ciclo de CPU
	}
	
	// 2. Simulação de Latência (IO/Rede/Banco)
	// Adicionamos um pequeno "jitter" (variação) de 10% para ficar realista e não ser um valor fixo robótico
	jitter := rand.Intn(atrasoMs/10 + 1) 
	tempoTotal := time.Duration(atrasoMs + jitter) * time.Millisecond
	time.Sleep(tempoTotal)
	
	log.Printf("Processamento finalizado em %v", time.Since(start))
	// --------------------------------------------------------

	produto, existe := produtos[idProduto]
	if existe {
		return produto, nil
	}

	return nil, fmt.Errorf("produto com ID %s não encontrado", idProduto)
}

func iniciarServidorHttp() {
    http.HandleFunc("/produto/", func(w http.ResponseWriter, r *http.Request) {
		// Aplica a mesma lógica para o REST, se for usado
		atrasoMs := getEnvInt("ATRASO_MS", 100)
		time.Sleep(time.Duration(atrasoMs) * time.Millisecond)

        id := r.URL.Path[len("/produto/"):]
        produto, existe := produtos[id]
        if existe {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(produto)
        } else {
            http.NotFound(w, r)
        }
    })
    log.Println("Servidor REST Catálogo escutando em :8081")
    http.ListenAndServe(":8081", nil)
}

func main() {
	go iniciarServidorHttp()
	porta := ":50051"
	lis, err := net.Listen("tcp", porta)
	if err != nil {
		log.Fatalf("Falha ao escutar na porta %s: %v", porta, err)
	}

	s := grpc.NewServer()
	pb.RegisterCatalogoServer(s, &servidor{})

	log.Printf("Servidor Catálogo escutando em %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}