package query

import "github.com/andrei-cosmin/hakkt/entity"

type Response interface {
	GetEntities() []entity.Id
}
