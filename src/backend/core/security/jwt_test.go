package security

import (
	"cognix.ch/api/v2/core/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJwtService_Create(t *testing.T) {

	js := NewJWTService("secret", 60)

	identity := Identity{
		AccessToken:  "ya29.a0Ad52N3-ivglS4hoLsoLmgGC6Xp-culrkTZEY5sCs1AknfFsxOOiiVIKtEU7zoA7hDUAJ8y2wbaaoRjI7e8h0l_1iOo0cJSLJBId7G7P6nX370UWx_dvMcaS_ExEBG5N64f6ZWNAMXPyoeESkUi7rOzxWf3yqSQlmLgaCgYKAckSARMSFQHGX2Mi5FKu4kDXffW-8Da4is1WnA0169",
		RefreshToken: "",
		User: &model.User{
			ID:         uuid.New(),
			TenantID:   uuid.New(),
			UserName:   "andrey.paladiychuk@gmail.com",
			FirstName:  "Andrii",
			LastName:   "Paladiichuk",
			ExternalID: "102087316285176952515",
			Roles:      nil,
			Tenant:     nil,
		},
	}
	token, err := js.Create(&identity)
	assert.NoError(t, err)
	t.Log(token)
}
