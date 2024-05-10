package db

import (
	"context"
	"database/sql"
	"fmt"

	_"github.com/lib/pq"

	"github.com/brm/utils"
)

type Store interface {
	Querier
	UserSignUpTx(context.Context, CreateUserParams, AddUserDeviceParams) (UserSignUpTxResult, error)
}

// SQLStore provides all functions to execute SQL (postgres) queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func NewPostgresDBConnection(config utils.Config) (*sql.DB, error) {
	con, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		return nil, err
	}

	if err := con.Ping(); err != nil {

		con.Close()

		return nil, err

	}

	return con, nil
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type UserSignUpTxResult struct {
	User        User           	`json:"user"`
	Discount    Discount    	`json:"discount"`
	UserDevice  UserDeviceList 	`json:"user_device"`
}

func (store *SQLStore) UserSignUpTx(ctx context.Context, createUserArg CreateUserParams, addUserDeviceArg AddUserDeviceParams) (UserSignUpTxResult, error) {
	var result UserSignUpTxResult

	err := store.execTx(ctx, func(q *Queries) error {

		var err error
		result.User, err = q.CreateUser(ctx, createUserArg)
		if err != nil {
			return err
		}

		userID := result.User.ID

		// Assign to addUserDeviceArg.UserID
		addUserDeviceArg.UserID = userID

		result.UserDevice, err = q.AddUserDevice(ctx, addUserDeviceArg)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
