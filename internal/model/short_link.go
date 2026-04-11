package model

import (
	"oat431/fluffy-mouton/pkg/common"

	"github.com/google/uuid"
)

type ShortLink struct {
	common.BaseEntity

	OwnBy uuid.UUID `db:"own_by" json:"own_by"`
	View  int       `db:"view" json:"view"`

	TargetURL string   `db:"target_url" json:"target_url"`
	ShortURL  string   `db:"short_url" json:"short_url"`
	Type      LinkType `db:"type" json:"type"`
}
