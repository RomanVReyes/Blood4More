package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Horario struct {
	Dia        string `bson:"dia"`
	Apertura string `bson:"horaInicio"`
	Cierre    string `bson:"horaFin"`
}

type Contacto struct {
	Nombre   string `bson:"nombre"`
	Telefono string `bson:"telefono"`
}

type Personal struct {
	Nombre string `bson:"nombre"`
	Rol    string `bson:"rol"`
}

type ResumenInventario struct {
	TipoSangre string `bson:"tipoSangre"`
	Cantidad   int    `bson:"cantidad"`
}

type BancoSangre struct {
	ID                      primitive.ObjectID `bson:"_id,omitempty"`
	NombreBanco             string             `bson:"nombreBanco"`
	TipoInstitucion         string             `bson:"tipoInstitucion"`
	Email                   string             `bson:"email"`    // Para el login separado
	Password                string             `bson:"password"` // Hash
	UbicacionGeografica     Ubicacion          `bson:"ubicacionGeografica"`
	Horarios                []Horario          `bson:"horarios"`
	CapacidadAlmacenamiento int                `bson:"capacidadAlmacenamiento"`
	ContactoEmergencia      Contacto           `bson:"contactoEmergencia"`
	PersonalResponsable     []Personal         `bson:"personalResponsable"`
	UmbralMinimo            map[string]int     `bson:"umbralMinimo"`
	// El inventario empezará vacío al registrarse
	Inventario []ResumenInventario `bson:"inventario"`
}
