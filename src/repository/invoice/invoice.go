package Invoices

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Risuii/invoice/src/entity"
	"github.com/Risuii/invoice/src/v1/contract"
)

func BuildFilter(query string, params contract.GetListParam) (string, contract.GetListParam) {

	if params.InvoiceID != "" {
		query += ` AND invoice_id = :invoice_id`
	}

	if params.IssueDate != "" {
		query += ` AND issue_date = :issue_date`
	}

	if params.Subject != "" {
		query += ` AND subject = :subject`
	}

	if params.TotalItem > 0 {
		query += ` AND total_items = :total_item`
	}

	if params.Customer != "" {
		query += ` AND c.name = :customer`
	}

	if params.DueDate != "" {
		query += ` AND due_date = :due_date`
	}

	if params.Status != "" {
		query += ` AND status = :status`
	}

	if params.Offset != 0 || params.Page == 1 {
		query += " LIMIT :limit OFFSET :offset"
	}

	return query, params
}

func (t *InvoicesRepository) Create(ctx context.Context, data *entity.Invoices) (contract.InvoiceResponseDB, error) {
	var res contract.InvoiceResponseDB

	namedStmt, err := t.getNamedStatement(ctx, InsertInvoice)
	if err != nil {
		log.Println("getNamedStatement err: ", err)
		return res, err
	}

	if err = namedStmt.GetContext(ctx, &res, data); err != nil {
		log.Println("get invoice err: ", err)
		return res, err
	}

	redisErr := t.redis.DelWithPattern(ctx, DeleteInvoiceRedisKey)
	if redisErr != nil {
		log.Println(redisErr)
	}

	return res, nil
}

func (t *InvoicesRepository) GetList(ctx context.Context, params contract.GetListParam) ([]*entity.Invoices, error) {
	var Invoices []*entity.Invoices

	stringQuery := masterQueries[GetList]
	query, params := BuildFilter(stringQuery, params)

	param, err := json.Marshal(params)
	if err != nil {
		log.Println("marshal err: ", err)
		return nil, err
	}

	err = t.redis.WithCache(ctx, fmt.Sprintf(GetListInvoicesRedisKey, param), &Invoices, func() (interface{}, error) {
		rows, err := t.db.NamedQueryContext(ctx, query, params)
		if err != nil {
			log.Println("named query err: ", err)
			return nil, err
		}

		for rows.Next() {
			var dataInvoices entity.Invoices
			err = rows.StructScan(&dataInvoices)
			if err != nil {
				return nil, err
			}

			Invoices = append(Invoices, &dataInvoices)
		}

		return Invoices, nil
	})

	if err != nil {
		log.Println("GetInvoicesList err: ", err)
		return nil, err
	}

	return Invoices, nil
}

func (t *InvoicesRepository) GetInvoicesCount(ctx context.Context, param contract.GetListParam) (int64, error) {
	var count int64

	params, err := json.Marshal(param)
	if err != nil {
		log.Println("marshal err: ", err)
		return 0, err
	}

	err = t.redis.WithCache(ctx, fmt.Sprintf(GetInvoicesCountRedisKey, params), &count, func() (interface{}, error) {
		var countData int64
		err := t.masterStmts[GetCountList].Get(&countData)
		return countData, err
	})

	if err != nil {
		log.Println("GetInvoicesCount err: ", err)
		return 0, err
	}

	return count, nil
}

func (t *InvoicesRepository) Get(ctx context.Context, id string) (entity.Invoices, error) {
	var Invoices entity.Invoices
	err := t.redis.WithCache(ctx, fmt.Sprintf(GetDetailInvoicesRedisKey, id), &Invoices, func() (interface{}, error) {
		var InvoicesData entity.Invoices
		err := t.masterStmts[GetByID].GetContext(ctx, &InvoicesData, id)
		return InvoicesData, err
	})

	if err != nil {
		log.Println(err)
		return Invoices, err
	}

	return Invoices, nil
}

func (t *InvoicesRepository) GetLatestInvoiceID(ctx context.Context) (string, error) {
	var res string

	err := t.masterStmts[GetLatestInvoiceID].GetContext(ctx, &res)
	if err != nil {
		log.Println("get latest invoice id err: ", err)
		return res, err
	}

	return res, nil
}

func (t InvoicesRepository) Update(ctx context.Context, data *entity.Invoices) error {
	var rowsAffected int64

	namedStmt, err := t.getNamedStatement(ctx, UpdateInvoice)
	if err != nil {
		log.Println("get named statement err: ", err)
		return err
	}

	res, err := namedStmt.ExecContext(ctx, data)
	if err != nil {
		log.Println("exec err: ", err)
		return err
	}

	rowsAffected, err = res.RowsAffected()
	if err != nil {
		log.Println("Get rows affected err: ", err)
		return err
	}

	if rowsAffected == 0 {
		log.Println("ID not exist err: ", sql.ErrNoRows)
		return sql.ErrNoRows
	}

	redisErr := t.redis.DelWithPattern(ctx, DeleteInvoiceRedisKey)
	if redisErr != nil {
		log.Println(redisErr)
	}

	return nil
}
