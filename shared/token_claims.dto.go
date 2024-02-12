package shared

import (
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/ctx"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaimsDTO struct {
	ID          uint     `bind:"id"`
	StoreID     uint     `bind:"store_id"`
	FirstName   string   `bind:"first_name"`
	LastName    string   `bind:"last_name"`
	Email       string   `bind:"email"`
	Permissions []string `bind:"permissions"`
}

func (instance TokenClaimsDTO) Transform(c gooh.Context, medata common.ArgumentMetadata) any {
	tokenClaims := c.Context().Value("tokenClaims").(jwt.MapClaims)
	tokenClaimsDTO, _ := ctx.BindStruct(tokenClaims, &[]ctx.FieldLevel{}, TokenClaimsDTO{}, "")

	return tokenClaimsDTO
}
