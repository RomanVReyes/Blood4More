package models

import "time"

type Donante struct {
	ID             string    `bson:"_id,omitempty" json:"id"`
	Nombre         string    `bson:"nombre" json:"nombre"`
	Email          string    `bson:"email" json:"email"`
	Password       string    `bson:"password" json:"-"` // El password no se envía en JSON
	Edad           int       `bson:"edad" json:"edad"`
	Peso           float64   `bson:"peso" json:"peso"`
	TipoSangre     string    `bson:"tipoSangre" json:"tipoSangre"`
	Ubicacion      string    `bson:"ubicacion" json:"ubicacion"`
	UltimaDonacion time.Time `bson:"ultimaDonacion" json:"ultimaDonacion"`
	FechaRegistro  time.Time `bson:"fechaRegistro" json:"fechaRegistro"`
}
