package worm

import (
	"context"
	"database/sql"
)

type SqlContex struct {
	ctx        	context.Context
	use_prepare 	sql.NullBool
	use_master 		sql.NullBool
	show_log    	sql.NullBool
}

func (ctx *SqlContex)Reset() {
	ctx.ctx = nil
	ctx.use_prepare.Valid = false
	ctx.use_master.Valid = false
	ctx.show_log.Valid = false
}