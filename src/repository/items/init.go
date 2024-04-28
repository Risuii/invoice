package items

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
	AllFields = `id, invoice_id, item_id, name, type, quantity, unit_price, amount`

	GetByInvoiceID = iota + 100
	DeleteItemByItemID

	InsertItems = iota + 200
	UpdateItems

	// Redis Key

	GetItemsByInvoiceIDRedisKey = "invoice:items:invoiceid:%s"
	DeleteItemRedisKey          = "invoice:items:*"
)

var (
	masterQueries = []string{
		GetByInvoiceID:     fmt.Sprintf("SELECT %s FROM items WHERE invoice_id = $1 AND deleted_at IS NULL", AllFields),
		DeleteItemByItemID: `UPDATE items SET deleted_at = now() WHERE item_id = $1 and deleted_at IS NULL`,
	}

	masterNamedQueries = []string{
		InsertItems: `INSERT INTO items (invoice_id, item_id, name, type, quantity, unit_price, amount) VALUES (:invoice_id, :item_id, :name, :type, :quantity, :unit_price, :amount)`,
		UpdateItems: `UPDATE items SET (invoice_id, name, type, quantity, unit_price, amount) = (:invoice_id, :name, :type, :quantity, :unit_price, :amount) WHERE item_id = :item_id`,
	}
)

type ItemsRepository struct {
	db                *sqlx.DB
	masterStmts       []*sqlx.Stmt
	masterNamedStmpts []*sqlx.NamedStmt
	redis             frsRedis.Redis
}

func InitItemsRepository(ctx context.Context, db *sqlx.DB, redis frsRedis.Redis) (*ItemsRepository, error) {
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

	return &ItemsRepository{
		db:                db,
		masterStmts:       stmpts,
		masterNamedStmpts: namedStmpts,
		redis:             redis,
	}, nil
}

func (r *ItemsRepository) getStatement(ctx context.Context, queryId int) (*sqlx.Stmt, error) {
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

func (r *ItemsRepository) getNamedStatement(ctx context.Context, queryId int) (*sqlx.NamedStmt, error) {
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
