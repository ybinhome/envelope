package envelopes

import (
	"context"
	"github.com/segmentio/ksuid"
	"github.com/tietang/dbx"
	"github.com/ybinhome/envelope/infra/base"
	"github.com/ybinhome/envelope/services"
)

type itemDomain struct {
	RedEnvelopeItem
}

// 1. 生成 ItemNo
func (i *itemDomain) createItemNo() {
	i.ItemNo = ksuid.New().Next().String()
}

// 2. 创建 Item
func (i *itemDomain) Create(item services.RedEnvelopeItemDTO) {
	i.RedEnvelopeItem.FromDTO(&item)
	i.RecvUsername.Valid = true
	i.createItemNo()
}

// 3. 保存 Item
func (i *itemDomain) Save(ctx context.Context) (id int64, err error) {
	err = base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeItemDao{runner: runner}
		id, err = dao.Insert(&i.RedEnvelopeItem)
		return err
	})
	return id, err
}

// 4. 通过 ItemNo 查询抢红包明细数据
func (i *itemDomain) GetOne(ctx context.Context, itemNo string) (dto *services.RedEnvelopeItemDTO) {
	err := base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeItemDao{runner: runner}
		po := dao.GetOne(itemNo)
		if po != nil {
			dto = po.ToDTO()
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return dto
}

// 5. 通过 envelopeNo 查询已抢到的红包列表
func (i *itemDomain) FindItems(envelopeNo string) (itemDtos []*services.RedEnvelopeItemDTO) {
	var items []*RedEnvelopeItem
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeItemDao{runner: runner}
		items = dao.FindItems(envelopeNo)
		return nil
	})
	if err != nil {
		return itemDtos
	}

	itemDtos = make([]*services.RedEnvelopeItemDTO, 0)
	for _, po := range items {
		itemDtos = append(itemDtos, po.ToDTO())
	}

	return itemDtos
}
