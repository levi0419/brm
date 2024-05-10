package shared

import (
	
	db "github.com/brm/db/sqlc"
	"github.com/brm/token"
	"github.com/brm/utils"
)

type IServer interface {
	GetConfig() utils.Config
	GetTokenMaker() token.ITokenMaker
	GetDbStore() db.Store

}
