package pago

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
	validator "gopkg.in/go-playground/validator.v9"
)

// User data structure
type Pago struct {
	ID              objectid.ObjectID `bson:"_id"`
	Autorizado      bool              `bson:"autorizado"`
	FechaPago       time.Time         `bson:"fechaPago"`
	FechaAutorizado time.Time         `bson:"fechaAutorizado"`
}

func NewPago() *Pago {
	return &Pago{
		ID:         objectid.New(),
		Autorizado: false,
		FechaPago:  time.Now(),
	}
}

func (e *Pago) ValidateSchema() error {
	validate := validator.New()
	return validate.Struct(e)
}
