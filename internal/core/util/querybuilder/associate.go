package querybuilder

import (
	"gorm.io/gorm"
)

type associate struct {
	cond        *Cond
	listBuilder []Builder
}

func Associate(args ...Builder) *associate {
	asso := &associate{
		listBuilder: args,
	}
	return asso
}

func (asso *associate) Build(tx *gorm.DB) *gorm.DB {
	// build condition
	if asso.cond != nil {
		tx = asso.cond.Build(tx)
	}

	// build preload
	for _, v := range asso.listBuilder {
		tx = v.Build(tx)
	}

	return tx
}

func Preload(table string, conds ...*Cond) *preload {
	if notEmptyString(table) {
		return &preload{
			table: table,
			cond:  And(conds...),
		}
	}
	return nil
}

type preload struct {
	cond  *Cond
	table string
}

func (p *preload) Build(tx *gorm.DB) *gorm.DB {
	if notEmptyString(p.table) {
		tx = tx.Preload(p.table, func(tmpTx *gorm.DB) *gorm.DB {
			if p.cond != nil {
				tmpTx = p.cond.Build(tmpTx)
			}
			return tmpTx
		})
	}

	return tx
}

func Join(table string, conds ...*Cond) *join {
	if notEmptyString(table) {
		return &join{
			cond:  And(conds...),
			table: table,
		}
	}
	return nil
}

type join struct {
	cond  *Cond
	table string
}

func (p *join) Build(tx *gorm.DB) *gorm.DB {
	if notEmptyString(p.table) {
		tx = tx.Joins(p.table, func(tmpTx *gorm.DB) *gorm.DB {
			if p.cond != nil {
				tmpTx = p.cond.Build(tmpTx)
			}
			return tmpTx
		})
	}

	return tx
}
