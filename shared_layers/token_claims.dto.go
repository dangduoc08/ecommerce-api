package sharedLayers

import (
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/ctx"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaimContextKey string

type TokenClaimsDTO struct {
	ID          uint     `bind:"id"`
	StoreID     uint     `bind:"store_id"`
	FirstName   string   `bind:"first_name"`
	LastName    string   `bind:"last_name"`
	Email       string   `bind:"email"`
	Permissions []string `bind:"permissions"`
}

func (instance TokenClaimsDTO) Transform(c gogo.Context, medata common.ArgumentMetadata) any {
	tokenClaims := c.Context().Value(TokenClaimContextKey("tokenClaims")).(jwt.MapClaims)
	tokenClaimsDTO, _ := ctx.BindStruct(tokenClaims, &[]ctx.FieldLevel{}, TokenClaimsDTO{}, "", "")

	return tokenClaimsDTO
}
