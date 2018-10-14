package pago

import (
	//"github.com/nmarsollier/authgo/security"
	"github.com/mauricioryan/pagosgo/tools/db"
	//"github.com/nmarsollier/authgo/security"
)

type serviceImpl struct {
	dao        Dao
	secService security.Service
}

// Service es la interfaz ue define el servicio
type Service interface {
	NewPago(pago *PagoRequest) (string, error)
	//Enable(userID string) error
	Pagos() ([]*User, error)
}

// NewService retorna una nueva instancia del servicio
func NewService() (Service, error) {
	secService, err := security.NewService()
	if err != nil {
		return nil, err
	}

	dao, err := newDao()
	if err != nil {
		return nil, err
	}

	return serviceImpl{
		dao:        dao,
		secService: secService,
	}, nil
}

// MockedService permite mockear el servicio
func MockedService(fakeDao Dao, fakeTRepo security.Service) Service {
	return serviceImpl{
		dao:        fakeDao,
		secService: fakeTRepo,
	}
}

// editar para crear un nuevo pago
func (s serviceImpl) NewPago(pago *PagoRequest) (string, error) {
	newPago := NewPago()
	newPago.Autorizado = false
	newPago.FechaPago = pago.FechaPago

	newPago, err := s.dao.Insert(newPago)
	if err != nil {
		if db.IsUniqueKeyError(err) {
			return "", ErrLoginExist
		}
		return "", err
	}

}

// Get wrapper para obtener un pago
func (s serviceImpl) Get(pagoID string) (*Pago, error) {
	return s.dao.FindByID(pagoID)
}

/*
//Enable habilita un usuario
func (s serviceImpl) Enable(userID string) error {
	usr, err := s.dao.FindByID(userID)
	if err != nil {
		return err
	}

	usr.Enabled = true
	_, err = s.dao.Update(usr)

	return err
}
*/

// Users wrapper para obtener todos los pagos
func (s serviceImpl) Pagos() ([]*Pago, error) {
	return s.dao.FindAll()
}
