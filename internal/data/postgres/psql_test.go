package postgres

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPsqlPool_ShouldReturnErrorWithoutConnection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool, cleanUp, err := NewPsqlPool(&Config{PostgresUrl: "some url"})

	assert.Nil(t, pool)
	assert.Nil(t, cleanUp)
	assert.NotNil(t, err)
}
