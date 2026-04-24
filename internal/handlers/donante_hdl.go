package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"roman-sangre/internal/models"
	"roman-sangre/internal/repository"
)

type DatosRegistro struct {
	Padecimientos []models.ItemCatalogo
	Medicamentos  []models.ItemCatalogo
}

func ShowDonorRegister(w http.ResponseWriter, r *http.Request) {
	// ---------------------------------------------------
	// METODO GET: Mostrar formulario con datos de la BD
	// ---------------------------------------------------
	if r.Method == http.MethodGet {
		padecimientos, medicamentos, err := repository.GetCatalogosMedicos()
		if err != nil {
			log.Println("Error obteniendo catálogos:", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return
		}

		datos := DatosRegistro{
			Padecimientos: padecimientos,
			Medicamentos:  medicamentos,
		}

		tmpl := template.Must(template.ParseFiles("templates/donante_registro.html"))
		tmpl.Execute(w, datos)
		return
	}

	// ---------------------------------------------------
	// METODO POST: Procesar y guardar registro
	// ---------------------------------------------------
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error procesando el formulario", http.StatusBadRequest)
			return
		}

		// Convertir tipos de datos
		edad, _ := strconv.Atoi(r.FormValue("edad"))
		peso, _ := strconv.ParseFloat(r.FormValue("peso"), 64)
		latitud, _ := strconv.ParseFloat(r.FormValue("latitud"), 64)
		longitud, _ := strconv.ParseFloat(r.FormValue("longitud"), 64)

		// Extraer tipo de sangre y factor RH ("O+" -> Tipo: "O", RH: "+")
		tipoSangreCompleto := r.FormValue("tipoSangre")
		tipo := strings.TrimRight(tipoSangreCompleto, "+-")
		factorStr := "positivo"
		if strings.HasSuffix(tipoSangreCompleto, "-") {
			factorStr = "negativo"
		}

		// Construir el modelo Donante
		nuevoDonante := models.Donante{
			// Nota: Mongo creará un ObjectID único automáticamente
			Nombre:     r.FormValue("nombre"),
			Password:   r.FormValue("password"), // ¡Importante: En el futuro hay que encriptarla!
			Edad:       edad,
			Genero:     r.FormValue("genero"),
			Peso:       peso,
			TipoSangre: tipo,
			FactorRH:   factorStr,
			UbicacionGeografica: models.Ubicacion{
				Coordenadas: models.Coordenadas{
					Latitud:  latitud,
					Longitud: longitud,
				},
				Direccion: r.FormValue("direccion"),
				Zona:      r.FormValue("zona"),
			},
			DatosContacto: models.DatosContacto{
				Correo:   r.FormValue("email"),
				Telefono: r.FormValue("telefono"),
				Whatsapp: r.FormValue("whatsapp"),
			},
			CondicionesMedicas: models.CondicionesMedicas{
				Padecimientos: r.Form["padecimientos"], // Parsea automáticamente los arrays de checkboxes
				Medicamentos:  r.Form["medicamentos"],
			},
			PreferenciasNotificacion: r.Form["notificaciones"],
			FechaRegistro:            time.Now().Format("2006-01-02"),
		}

		// Guardar en la base de datos
		err = repository.CreateDonante(nuevoDonante)
		if err != nil {
			log.Println("Error creando donante:", err)
			http.Error(w, "Error al registrar la cuenta", http.StatusInternalServerError)
			return
		}

		// Redirigir al dashboard tras el registro exitoso
		http.Redirect(w, r, "/donante/dashboard", http.StatusSeeOther)
		return
	}

	// Si no es GET ni POST
	http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
}

// Mostrar la página principal de opciones para el donante (Auth)
func ShowDonorAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "templates/donante_auth.html")
		return
	}
	http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
}

// Mostrar y procesar el inicio de sesión del donante
func ShowDonorLogin(w http.ResponseWriter, r *http.Request) {
	// Mostrar formulario de login
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "templates/donante_login.html")
		return
	}

	// Procesar formulario de login (POST)
	if r.Method == http.MethodPost {
		// Por ahora solo haremos una redirección simple al dashboard.
		// TODO: Más adelante implementaremos la validación en la base de datos comparando correos y contraseñas.

		http.Redirect(w, r, "/donante/dashboard", http.StatusSeeOther)
		return
	}

	http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
}
