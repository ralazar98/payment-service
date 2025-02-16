package storage

import (
	"fmt"
	"github.com/jackc/pgx"
	"log"
	"payment-service/configs"
	"payment-service/internal/entity"
	"time"
)

type BankStore struct {
	conn *pgx.ConnPool
}

func New(cfg configs.DatabaseConfig) (*BankStore, error) {

	connConf := pgx.ConnConfig{
		Host:     cfg.Host,
		Port:     uint16(cfg.Port),
		User:     cfg.Username,
		Password: cfg.DBPassword,
		Database: cfg.DBName,
	}

	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     connConf,
		MaxConnections: 10,
	})
	if err != nil {
		return nil, fmt.Errorf("Error creating new connection pool %w ", err)
	}
	log.Println("Successfully created new connection pool")
	return &BankStore{
		conn: pool,
	}, nil

}

func (s *BankStore) Update(user *entity.UpdateBalance) {
	request := ` UPDATE public.bank_storage SET balance = balance+$2,date_updated=$3 WHERE user_id = $1 `

	err := s.conn.QueryRow(request, user.UserID, user.ChangingInBalance, time.Now())
	if err != nil {

		log.Println("update DB Error: ", err)
		//TODO:Опять же оишбку нужно обработать
		// даже в слушае когда у тебя ничего не меняется, лучше написать return
	}

}
