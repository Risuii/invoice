package Invoices

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	frsAtomic "github.com/Risuii/frs-lib/atomic"
	atomicSqlx "github.com/Risuii/frs-lib/atomic/sqlx"
	frsRedis "github.com/Risuii/frs-lib/redis"
	sqlxUtils "github.com/Risuii/frs-lib/sqlx"
)

const (
	AllFields           = `id, invoice_id, issue_date, subject, total_items, customer_id, due_date, status, sub_total, tax, grand_total, created_at, updated_at`
	AllFieldsForGetList = `t.id, t.invoice_id, t.issue_date, t.subject, t.total_items, c.name AS customer_name, t.due_date, t.status, t.sub_total, t.tax, t.grand_total, t.created_at, t.updated_at`

	BaseQuery = iota + 100
	GetByID
	GetByInvoiceID
	GetList
	GetCountList
	GetLatestInvoiceID

	InsertInvoice = iota + 200
	UpdateInvoice

	// Redis Key

	GetListInvoicesRedisKey   = "invoice:invoices:getlist:%s"
	GetDetailInvoicesRedisKey = "invoice:invoices:getdetail:%s"
	GetInvoicesCountRedisKey  = "invoice:invoices:getcount:%s"
	DeleteInvoiceRedisKey     = "invoice:invoices:*"
)

var (
	masterQueries = []string{
		BaseQuery:          fmt.Sprintf("SELECT %s FROM Invoices", AllFields),
		GetByID:            fmt.Sprintf("SELECT %s FROM Invoices WHERE invoice_id = $1 AND deleted_at IS NULL", AllFields),
		GetByInvoiceID:     fmt.Sprintf("SELECT %s FROM Invoices WHERE invoice_id = $1 And deleted_at IS NULL", AllFields),
		GetList:            fmt.Sprintf(`SELECT %s FROM Invoices as t INNER JOIN customers as c ON t.customer_id = c.customer_id  WHERE t.deleted_at IS NULL`, AllFieldsForGetList),
		GetCountList:       `SELECT COUNT(*) FROM Invoices WHERE deleted_at IS NULL`,
		GetLatestInvoiceID: `SELECT MAX(invoice_id) FROM invoices`,
	}

	masterNamedQueries = []string{
		InsertInvoice: `INSERT INTO invoices (invoice_id, issue_date, subject, total_items, customer_id, due_date, status, sub_total, tax, grand_total) VALUES (:invoice_id, :issue_date, :subject, :total_items, :customer_id, :due_date, :status, :sub_total, :tax, :grand_total) RETURNING invoice_id, customer_id`,
		UpdateInvoice: `UPDATE invoices SET (issue_date, subject, total_items, due_date, sub_total, tax, grand_total) = (:issue_date, :subject, :total_items, :due_date, :sub_total, :tax, :grand_total) WHERE invoice_id = :invoice_id`,
	}
)

type InvoicesRepository struct {
	db                *sqlx.DB
	masterStmts       []*sqlx.Stmt
	masterNamedStmpts []*sqlx.NamedStmt
	redis             frsRedis.Redis
}

func InitInvoicesRepository(ctx context.Context, db *sqlx.DB, redis frsRedis.Redis) (*InvoicesRepository, error) {
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

	return &InvoicesRepository{
		db:                db,
		masterStmts:       stmpts,
		masterNamedStmpts: namedStmpts,
		redis:             redis,
	}, nil
}

func (r *InvoicesRepository) getStatement(ctx context.Context, queryId int) (*sqlx.Stmt, error) {
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

func (r *InvoicesRepository) getNamedStatement(ctx context.Context, queryId int) (*sqlx.NamedStmt, error) {
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
