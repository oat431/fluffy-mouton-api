package model

import (
	"oat431/fluffy-mouton/pkg/common"
	"time"

	"github.com/google/uuid"
)

type EmailVerifyToken struct {
	common.BaseEntity

	AuthID    uuid.UUID `db:"auth_id" json:"auth_id"`
	Token     string    `db:"token" json:"token"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
}
