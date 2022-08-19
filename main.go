package main

import (
	"database/sql"
	"fmt"
	"github.com/patrickrhamon/codebank/infrastructure/grpc/server"
	"github.com/patrickrhamon/codebank/infrastructure/kafka"
	"github.com/patrickrhamon/codebank/infrastructure/repository"
	"github.com/patrickrhamon/codebank/usecase"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
}

func main() {
	db := setupDb()
	defer db.Close()
	producer := setupKafkaProducer()
	processTransactionUseCase := setupTransactionUseCase(db, producer)
	serveGrpc(processTransactionUseCase)
}

func setupKafkaProducer() kafka.KafkaProducer {
	producer := kafka.NewKafkaProducer()
	producer.SetupProducer(os.Getenv("KafkaBootstrapServers"))
	return producer
}

func setupTransactionUseCase(db *sql.DB, producer kafka.KafkaProducer) usecase.UseCaseTransaction {
	transactionRepository := repository.NewTransactionRepositoryDb(db)
	useCase := usecase.NewUseCaseTransaction(transactionRepository)
	useCase.KafkaProducer = producer
	return useCase
}

func setupDb() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("host"),
		os.Getenv("port"),
		os.Getenv("user"),
		os.Getenv("password"),
		os.Getenv("dbname"),
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("error connection to database")
	}
	return db
}

func serveGrpc(processTransactionUseCase usecase.UseCaseTransaction) {
	grpcServer := server.NewGRPCServer()
	grpcServer.ProcessTransactionUseCase = processTransactionUseCase
	fmt.Println("Rodando gRPC Server")
	grpcServer.Serve()
}

// cc := domain.NewCreditCard()
// cc.Number = "1234"
// cc.Name = "Patrick"
// cc.ExpirationMonth = 7
// cc.ExpirationYear = 2022
// cc.CVV = 123
// cc.Limit = 10000
// cc.Balance = 500

// repo := repository.NewTransactionRepositoryDb(db)
// err := repo.CreateCreditCard(*cc)

// if err != nil {
// 	fmt.Println(err)
// }