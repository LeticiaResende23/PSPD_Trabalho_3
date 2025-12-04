package main

import (
	"context"
	"fmt"
	"log"
	"net"
    "time"       // Adicionado
    "math/rand"  // Adicionado
    "os"         // Adicionado
    "strconv"    // Adicionado
	"google.golang.org/grpc"
	pb "projeto-catalogo/proto" // Mantenha o import original se der erro
)

type servidor struct {
	pb.UnimplementedInventarioServer
}

func getEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func (s *servidor) GetEstoque(ctx context.Context, req *pb.ProdutoRequest) (*pb.EstoqueResponse, error) {
    idProduto := req.GetId()
    
    // --- SIMULAÇÃO DE CARGA ---
    atrasoMs := getEnvInt("ATRASO_MS", 100)
    cargaCpu := getEnvInt("CARGA_CPU", 100000)

    log.Printf("[Inventário] Req ID: %v | Config: Delay=%dms, CPU=%d loops", idProduto, atrasoMs, cargaCpu)

    start := time.Now()
    for i := 0; i < cargaCpu; i++ {
        _ = float64(i) * rand.Float64()
    }
    
    // Jitter (variação) para ficar mais realista
    jitter := rand.Intn(atrasoMs/10 + 1)
    time.Sleep(time.Duration(atrasoMs + jitter) * time.Millisecond)

    log.Printf("[Inventário] Processamento finalizado em %v", time.Since(start))
    // -----------------------------

	if idProduto == "1" {
		return &pb.EstoqueResponse{Preco: 150.00, Quantidade: 100}, nil
	} else if idProduto == "2" {
		return &pb.EstoqueResponse{Preco: 350.50, Quantidade: 50}, nil
	}

	return nil, fmt.Errorf("produto com ID %s não encontrado no estoque", idProduto)
}

func main() {
	porta := ":50052"
	lis, err := net.Listen("tcp", porta)
	if err != nil {
		log.Fatalf("Falha ao escutar na porta %s: %v", porta, err)
	}

	s := grpc.NewServer()
	pb.RegisterInventarioServer(s, &servidor{})

	log.Printf("Servidor Inventário escutando em %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}