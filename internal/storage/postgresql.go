package storage

import (
	"github.com/jackc/pgx"
	"log"
	"payment-service/configs"
	"payment-service/internal/entity"
	"strconv"
	"time"
)

type BankStore struct {
	conn *pgx.ConnPool
}

func New(cfg configs.DatabaseConfig) *BankStore {

	dbPort, err := strconv.Atoi(cfg.Port)
	if err != nil {
		log.Fatal("Error converting DB_PORT to integer")
	}

	connConf := pgx.ConnConfig{
		Host:     cfg.Host,
		Port:     uint16(dbPort),
		User:     cfg.Username,
		Password: cfg.DBPassword,
		Database: cfg.DBName,
	}

	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     connConf,
		MaxConnections: 10,
	})
	//TODO:тут снизу что то очень странное
	if err != nil {
		log.Println("Error creating new connection pool", err)
		//TODO:Сразу возврящяем ошибку, в сигне ошибки нет, но если метож может вернуть ошибку, то всегда возвращаем
		// и получается ретюрн будет как return nil, err
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
	//TODO:Изменить время на время из бд
	err := s.conn.QueryRow(request, user.UserID, user.ChangingInBalance, dateUpdate)
	if err != nil {
		log.Println(err)
		//TODO:Опять же оишбку нужно обработать
		// даже в слушае когда у тебя ничего не меняется, лучше написать return
	}

}
