package customers

import (
	"context"
	"fmt"
	"log"

	frsAtomic "github.com/Risuii/frs-lib/atomic"
	atomicSqlx "github.com/Risuii/frs-lib/atomic/sqlx"
	frsRedis "github.com/Risuii/frs-lib/redis"
	sqlxUtils "github.com/Risuii/frs-lib/sqlx"
	"github.com/jmoiron/sqlx"
)

const (
	AllFields = `id, customer_id, name, address, created_at, updated_at`

	GetByID = iota + 100

	InsertCustomer = iota + 200
	UpdateCustomer

	// Redis Key

	GetDetailCustomersRedisKey = "invoice:customers:getdetail:%s"
	DeleteCustomerRedisKey     = "invoice:customers:*"
)

var (
	masterQueries = []string{
		GetByID: fmt.Sprintf("SELECT %s FROM customers WHERE customer_id = $1 AND deleted_at IS NULL", AllFields),
	}

	masterNamedQueries = []string{
		InsertCustomer: `INSERT INTO customers (customer_id, name, address) VALUES (:customer_id, :name, :address)`,
		UpdateCustomer: `UPDATE customers SET (customer_id, name, address) = (:customer_id, :name, :address) WHERE customer_id = :customer_id`,
	}
)

type CustomersRepository struct {
	db                *sqlx.DB
	masterStmts       []*sqlx.Stmt
	masterNamedStmpts []*sqlx.NamedStmt
	redis             frsRedis.Redis
}

func InitCustomersRepository(ctx context.Context, db *sqlx.DB, redis frsRedis.Redis) (*CustomersRepository, error) {
	stmpts, err := sqlxUtils.PrepareQueries(db, masterQueries)
	if err != nil {
		log.Println("PrepareQueries err:", err)
		return nil, err
	}

	namedStmpts, err := sqlxUtils.PrepareNamedQueries(db, masterNamedQueries)
	if err != nil {
		log.Println("PrepareNamedQueries err:", err)
		return nil, err
	}

	return &CustomersRepository{
		db:                db,
		masterStmts:       stmpts,
		masterNamedStmpts: namedStmpts,
		redis:             redis,
	}, nil
}

func (r *CustomersRepository) getStatement(ctx context.Context, queryId int) (*sqlx.Stmt, error) {
	var err error
	var statement *sqlx.Stmt
	if atomicSessionCtx, ok := ctx.(*frsAtomic.AtomicSessionContext); ok {
		if atomicSession, ok := atomicSessionCtx.AtomicSession.(*atomicSqlx.SqlxAtomicSession); ok {
			statement, err = atomicSession.Tx().PreparexContext(ctx, masterQueries[queryId])
		} else {
			err = frsAtomic.InvalidAtomicSessionProvider
		}
	} else {
		statement = r.masterStmts[queryId]
	}
	return statement, err
}

func (r *CustomersRepository) getNamedStatement(ctx context.Context, queryId int) (*sqlx.NamedStmt, error) {
	var err error
	var namedStmt *sqlx.NamedStmt
	if atomicSessionCtx, ok := ctx.(*frsAtomic.AtomicSessionContext); ok {
		if atomicSession, ok := atomicSessionCtx.AtomicSession.(*atomicSqlx.SqlxAtomicSession); ok {
			namedStmt, err = atomicSession.Tx().PrepareNamedContext(ctx, masterNamedQueries[queryId])
		} else {
			err = frsAtomic.InvalidAtomicSessionProvider
		}
	} else {
		namedStmt = r.masterNamedStmpts[queryId]
	}
	return namedStmt, err
}
