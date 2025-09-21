package utils

import (
	"testing"

	"github.com/fazendapro/FazendaPro-api/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {
	t.Run("HashPassword_Success", func(t *testing.T) {
		password := "senha123"
		hashedPassword, err := utils.HashPassword(password)

		assert.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
		assert.NotEqual(t, password, hashedPassword)
		assert.True(t, len(hashedPassword) > 0)
	})

	t.Run("HashPassword_EmptyPassword", func(t *testing.T) {
		password := ""
		hashedPassword, err := utils.HashPassword(password)

		assert.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
	})

	t.Run("CheckPassword_Success", func(t *testing.T) {
		password := "senha123"
		hashedPassword, _ := utils.HashPassword(password)

		isValid := utils.CheckPasswordHash(password, hashedPassword)
		assert.True(t, isValid)

		isValid = utils.CheckPasswordHash("senhaerrada", hashedPassword)
		assert.False(t, isValid)
	})
}
