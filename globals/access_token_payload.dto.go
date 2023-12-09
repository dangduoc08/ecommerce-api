package globals

import (
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/context"
	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenDTO struct {
	ID        uint   `bind:"id"`
	StoreID   uint   `bind:"store_id"`
	FirstName string `bind:"first_name"`
	LastName  string `bind:"last_name"`
	Email     string `bind:"email"`
}

func (self AccessTokenDTO) Transform(ctx gooh.Context, medata common.ArgumentMetadata) any {
	tokenClaims := ctx.Context().Value("tokenClaims").(jwt.MapClaims)
	return context.BindStruct(tokenClaims, AccessTokenDTO{})
}
