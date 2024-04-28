package items

import (
	"context"
	"fmt"
	"log"

	"github.com/Risuii/invoice/src/entity"
	"github.com/google/uuid"
)

func (i *ItemsRepository) Create(ctx context.Context, data []*entity.Item) error {
	namedStmt, err := i.getNamedStatement(ctx, InsertItems)
	if err != nil {
		log.Println("getNamedStatement err: ", err)
		return err
	}

	for _, v := range data {
		itemData := entity.Item{
			ItemData: entity.ItemData{
				InvoiceID: v.InvoiceID,
				ItemID:    v.ItemID,
				Name:      v.Name,
				Type:      v.Type,
				Quantity:  v.Quantity,
				UnitPrice: v.UnitPrice,
				Amount:    v.Amount,
			},
		}
		_, err = namedStmt.ExecContext(ctx, itemData)
	}

	if err != nil {
		log.Println("create items err: ", err)
		return err
	}

	redisErr := i.redis.DelWithPattern(ctx, DeleteItemRedisKey)
	if redisErr != nil {
		log.Println(redisErr)
	}

	return nil
}

func (i *ItemsRepository) GetByInvoiceID(ctx context.Context, invID string) ([]*entity.Item, error) {
	var Item []*entity.Item
	err := i.redis.WithCache(ctx, fmt.Sprintf(GetItemsByInvoiceIDRedisKey, invID), &Item, func() (interface{}, error) {
		var itemData []*entity.Item
		err := i.masterStmts[GetByInvoiceID].SelectContext(ctx, &itemData, invID)
		return itemData, err
	})

	if err != nil {
		log.Println(err)
		return Item, err
	}

	return Item, nil
}

func (i *ItemsRepository) Update(ctx context.Context, data []*entity.Item) error {

	namedStmt, err := i.getNamedStatement(ctx, UpdateItems)
	if err != nil {
		log.Println("get named statement err: ", err)
		return err
	}

	for _, v := range data {
		itemData := entity.Item{
			ItemData: entity.ItemData{
				InvoiceID: v.InvoiceID,
				ItemID:    v.ItemID,
				Name:      v.Name,
				Type:      v.Type,
				Quantity:  v.Quantity,
				UnitPrice: v.UnitPrice,
				Amount:    v.Amount,
			},
		}

		_, err = namedStmt.ExecContext(ctx, itemData)
	}

	if err != nil {
		log.Println("exec err: ", err)
		return err
	}

	redisErr := i.redis.DelWithPattern(ctx, DeleteItemRedisKey)
	if redisErr != nil {
		log.Println(redisErr)
	}

	return nil
}

func (i *ItemsRepository) Delete(ctx context.Context, ids []uuid.UUID) error {
	for _, v := range ids {
		_, err := i.masterStmts[DeleteItemByItemID].ExecContext(ctx, v)
		if err != nil {
			log.Println("DeleteProduct err: ", err)
			return err
		}
	}

	redisErr := i.redis.DelWithPattern(ctx, DeleteItemRedisKey)
	if redisErr != nil {
		log.Println("Redis delete error:", redisErr)
	}

	return nil
}
