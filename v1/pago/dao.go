package pago

import (
	"context"
	"log"

	"github.com/mauricioryan/pagosgo/tools/db"
	"github.com/mauricioryan/pagosgo/tools/errors"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type daoStruct struct {
	collection db.Collection
}

// Dao es la interface que exponse los servicios de acceso a la DB
type Dao interface {
	Insert(pago *Pago) (*Pago, error)
	Update(pago *Pago) (*Pago, error)
	FindAll() ([]*Pago, error)
	FindByID(pagoID string) (*Pago, error)
	//FindByLogin(login string) (*User, error)
}

// New dao es interno a este modulo, nadie fuera del modulo tiene acceso
func newDao() (Dao, error) {
	database, err := db.Get()
	if err != nil {
		return nil, err
	}

	collection := database.Collection("pagos")

	_, err = collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.NewDocument(
				bson.EC.String("login", ""),
			),
			Options: bson.NewDocument(
				bson.EC.Boolean("unique", true),
			),
		},
	)
	if err != nil {
		log.Output(1, err.Error())
	}

	coll := db.WrapCollection(collection)
	return daoStruct{
		collection: coll,
	}, nil
}

// MockedDao sirve para poder mockear el db.Collection y testear el modulo
func MockedDao(coll db.Collection) Dao {
	return daoStruct{
		collection: coll,
	}
}

func (d daoStruct) Insert(pago *Pago) (*Pago, error) {
	if err := pago.ValidateSchema(); err != nil {
		return nil, err
	}

	if _, err := d.collection.InsertOne(context.Background(), pago); err != nil {
		return nil, err
	}

	return pago, nil
}

/*
func (d daoStruct) Update(pago *pago) (*pago, error) {
	if err := pago.ValidateSchema(); err != nil {
		return nil, err
	}

	pago.Updated = time.Now()

	doc, err := bson.NewDocumentEncoder().EncodeDocument(pago)
	if err != nil {
		return nil, err
	}

	//TODO arreglar estructura para update
	_, err = d.collection.UpdateOne(context.Background(),
		bson.NewDocument(doc.LookupElement("_id")),
		bson.NewDocument(
			bson.EC.SubDocumentFromElements("$set",
				doc.LookupElement("password"),
				doc.LookupElement("name"),
				doc.LookupElement("enabled"),
				doc.LookupElement("updated"),
				doc.LookupElement("permissions"),
			),
		))

	if err != nil {
		return nil, err
	}

	return pago, nil
}
*/
// FindAll devuelve todos los usuarios
func (d daoStruct) FindAll() ([]*Pago, error) {
	filter := bson.NewDocument()
	cur, err := d.collection.Find(context.Background(), filter, nil)
	defer cur.Close(context.Background())

	if err != nil {
		return nil, err
	}

	pagos := []*Pago{}
	for cur.Next(context.Background()) {
		pago := &Pago{}
		if err := cur.Decode(pago); err != nil {
			return nil, err
		}
		pagos = append(pagos, pago)
	}

	return pagoss, nil
}

// FindByID lee un usuario desde la db
func (d daoStruct) FindByID(pagoID string) (*Pago, error) {
	_id, err := objectid.FromHex(pagoID)
	if err != nil {
		return nil, errors.ErrID
	}

	pago := &Pago{}
	filter := bson.NewDocument(bson.EC.ObjectID("_id", _id))
	if err = d.collection.FindOne(context.Background(), filter).Decode(pago); err != nil {
		return nil, err
	}

	return pago, nil
}

/*
// FindByLogin lee un usuario desde la db
func (d daoStruct) FindByLogin(login string) (*User, error) {
	user := &User{}
	filter := bson.NewDocument(bson.EC.String("login", login))
	err := d.collection.FindOne(context.Background(), filter).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrLogin
		}
		return nil, err
	}

	return user, nil
}
*/
