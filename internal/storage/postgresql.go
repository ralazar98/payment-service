package storage

import (
	"github.com/jackc/pgx"
	"log"
	"os"
	"payment-service/internal/entity"
	"strconv"
	"time"
)

type BankStore struct {
	conn *pgx.ConnPool
}

func New() *BankStore {
	dbHost := os.Getenv("DB_HOST")
	dbPortStr := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatal("Error converting DB_PORT to integer")
	}

	connConf := pgx.ConnConfig{
		Host:     dbHost,
		Port:     uint16(dbPort),
		User:     dbUser,
		Password: dbPassword,
		Database: dbName,
	}

	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     connConf,
		MaxConnections: 10,
	})
	if err != nil {
		log.Println("Error creating new connection pool")
	} else {
		log.Println("Successfully created new connection pool")
		return &BankStore{
			conn: pool,
		}
	}
	return &BankStore{}
}

func (s *BankStore) Update(user *entity.UpdateBalance) {
	request := ` UPDATE public.bank_storage SET balance = balance+$2,date_updated=$3 WHERE user_id = $1 `

	dateUpdate := time.Now().UTC().Format("2006-01-02")
	err := s.conn.QueryRow(request, user.UserID, user.ChangingInBalance, dateUpdate)
	if err != nil {
		log.Println(err)
	}

}
