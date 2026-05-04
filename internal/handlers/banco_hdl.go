package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"roman-sangre/internal/models"
	"roman-sangre/internal/repository"
)

type DatosRegistroBanco struct {
	ErrorMessage string
}

func ShowBancoRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/banco_registro.html"))
		tmpl.Execute(w, DatosRegistroBanco{})
		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error procesando el formulario", http.StatusBadRequest)
			return
		}

		email := r.FormValue("email")

		_, err = repository.GetBancoByEmail(email)
		if err == nil {
			datos := DatosRegistroBanco{
				ErrorMessage: "El correo electrónico institucional ya está registrado. Inicia sesión.",
			}
			tmpl := template.Must(template.ParseFiles("templates/banco_registro.html"))
			tmpl.Execute(w, datos)
			return
		} else if err != mongo.ErrNoDocuments {
			http.Error(w, "Error verificando la base de datos", http.StatusInternalServerError)
			return
		}

		capacidad, _ := strconv.Atoi(r.FormValue("capacidad"))
		latitud, _ := strconv.ParseFloat(r.FormValue("latitud"), 64)
		longitud, _ := strconv.ParseFloat(r.FormValue("longitud"), 64)

		passwordPlano := r.FormValue("password")
		hash, err := bcrypt.GenerateFromPassword([]byte(passwordPlano), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error procesando contraseña", http.StatusInternalServerError)
			return
		}

		tiposSangre := []string{"O+", "O-", "A+", "A-", "B+", "B-", "AB+", "AB-"}
		umbrales := make(map[string]int)
		for _, tipo := range tiposSangre {
			val, _ := strconv.Atoi(r.FormValue("umbral_" + tipo))
			umbrales[tipo] = val
		}

		inventarioInicial := []models.ResumenInventario{}
		for _, tipo := range tiposSangre {
			inventarioInicial = append(inventarioInicial, models.ResumenInventario{
				TipoSangre: tipo,
				Cantidad:   0,
			})
		}

		nuevoBanco := models.BancoSangre{
			NombreBanco:     r.FormValue("nombreBanco"),
			TipoInstitucion: r.FormValue("tipoInstitucion"),
			Email:           email,
			Password:        string(hash),

			UbicacionGeografica: models.Ubicacion{
				Coordenadas: models.Coordenadas{Latitud: latitud, Longitud: longitud},
				Direccion:   r.FormValue("direccion"),
			},

			Horarios: []models.Horario{
				{
					Dia:      "Lunes a Viernes",
					Apertura: r.FormValue("apertura"),
					Cierre:   r.FormValue("cierre"),
				},
			},

			CapacidadAlmacenamiento: capacidad,

			ContactoEmergencia: models.Contacto{
				Nombre:   "Emergencias",
				Telefono: r.FormValue("tel_emergencia"),
			},

			PersonalResponsable: []models.Personal{
				{
					Nombre: r.FormValue("resp_nombre"),
					Rol:    r.FormValue("resp_cargo"),
				},
			},

			UmbralMinimo: umbrales,
			Inventario:   inventarioInicial,
		}

		err = repository.CreateBanco(nuevoBanco)
		if err != nil {
			log.Println("Error creando banco de sangre:", err)
			http.Error(w, "Error al registrar la cuenta institucional", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/banco/login", http.StatusSeeOther)
		return
	}

	http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
}

func ShowBancoAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "templates/banco_auth.html")
		return
	}
	http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
}
