package model

import "oat431/fluffy-mouton/pkg/common"

type Auth struct {
	common.BaseEntity

	Username   string `db:"username" json:"username"`
	Email      string `db:"email" json:"email"`
	Password   string `db:"password" json:"-"`
	IsVerified bool   `db:"is_verified" json:"is_verified"`
}
