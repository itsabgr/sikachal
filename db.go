package sikachal

import (
	"context"
	"database/sql"
	"errors"
	"github.com/itsabgr/sikachal/internal/common"
	"math/rand"
	"runtime"
	"time"
)

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type DB struct {
	db *sql.DB
}

func NewDB(db *sql.DB) *DB {
	return &DB{db: db}
}
func (db *DB) GetUserByID(ctx context.Context, id uint64) (*User, error) {
	user, err := common.QueryRow[User](ctx, db.db, "select id,firstName,lastName from user where id = ?", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomWord(randSrc *rand.Rand) string {
	n := randSrc.Intn(5) + 5
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[randSrc.Intn(len(letterRunes))]
	}
	return string(b)
}
func (db *DB) Create(ctx context.Context) error {
	_, err := db.db.ExecContext(ctx, `create table main.user
(
    id        integer not null
        constraint user_pk
            primary key autoincrement,
    firstName text    not null,
    lastName  text    not null
)`)
	return err
}
func (db *DB) BulkUserInsert(ctx context.Context, count uint64) error {
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()
	randSrc := rand.New(rand.NewSource(time.Now().UnixNano()))
	for range make([]struct{}, count) {

		if err := ctx.Err(); err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, "insert into user (firstName,lastName) values (?,?)", randomWord(randSrc), randomWord(randSrc))
		if err != nil {
			return err
		}
		runtime.Gosched()
	}
	return tx.Commit()
}
