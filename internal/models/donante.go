package models

// Estructura para leer los catálogos de Mongo
type ItemCatalogo struct {
	ID     string `bson:"_id"`
	Nombre string `bson:"nombre"`
}

// Estructuras anidadas para el Donante
type Coordenadas struct {
	Latitud  float64 `bson:"latitud" json:"latitud"`
	Longitud float64 `bson:"longitud" json:"longitud"`
}

type Ubicacion struct {
	Coordenadas Coordenadas `bson:"coordenadas" json:"coordenadas"`
	Direccion   string      `bson:"direccion" json:"direccion"`
	Zona        string      `bson:"zona" json:"zona"`
}

type DatosContacto struct {
	Telefono string `bson:"telefono" json:"telefono"`
	Correo   string `bson:"correo" json:"correo"`
	Whatsapp string `bson:"whatsapp" json:"whatsapp"`
}

type CondicionesMedicas struct {
	Padecimientos []string `bson:"padecimientos" json:"padecimientos"`
	Medicamentos  []string `bson:"medicamentos" json:"medicamentos"`
}

// Estructura principal del Donante
type Donante struct {
	ID                       string             `bson:"_id,omitempty" json:"_id,omitempty"` // omitempty para que Mongo genere uno si no lo pasas
	Nombre                   string             `bson:"nombre" json:"nombre"`
	Password                 string             `bson:"password" json:"-"` // No lo devolvemos en JSON
	Edad                     int                `bson:"edad" json:"edad"`
	Genero                   string             `bson:"genero" json:"genero"`
	Peso                     float64            `bson:"peso" json:"peso"`
	TipoSangre               string             `bson:"tipoSangre" json:"tipoSangre"`
	FactorRH                 string             `bson:"factorRH" json:"factorRH"`
	UbicacionGeografica      Ubicacion          `bson:"ubicacionGeografica" json:"ubicacionGeografica"`
	DatosContacto            DatosContacto      `bson:"datosContacto" json:"datosContacto"`
	CondicionesMedicas       CondicionesMedicas `bson:"condicionesMedicas" json:"condicionesMedicas"`
	PreferenciasNotificacion []string           `bson:"preferenciasNotificacion" json:"preferenciasNotificacion"`
	FechaRegistro            string             `bson:"fechaRegistro" json:"fechaRegistro"`
}
