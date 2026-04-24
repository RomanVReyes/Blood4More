package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"roman-sangre/internal/models"
	"roman-sangre/internal/repository"
)

type DatosRegistro struct {
	Padecimientos []models.ItemCatalogo
	Medicamentos  []models.ItemCatalogo
	ErrorMessage  string
}

func ShowDonorRegister(w http.ResponseWriter, r *http.Request) {
	// ---------------------------------------------------
	// METODO GET: Mostrar formulario
	// ---------------------------------------------------
	if r.Method == http.MethodGet {
		padecimientos, medicamentos, _ := repository.GetCatalogosMedicos()
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

		email := r.FormValue("email")

		// 1. VALIDACIÓN DE DUPLICADOS
		_, err = repository.GetDonanteByEmail(email)
		if err == nil {
			// Si err es nil, significa que SÍ encontró un usuario con ese correo.
			padecimientos, medicamentos, _ := repository.GetCatalogosMedicos()
			datos := DatosRegistro{
				Padecimientos: padecimientos,
				Medicamentos:  medicamentos,
				ErrorMessage:  "El correo electrónico ingresado ya está registrado. Intenta iniciar sesión.",
			}

			// Volvemos a renderizar la página mostrando el error
			tmpl := template.Must(template.ParseFiles("templates/donante_registro.html"))
			tmpl.Execute(w, datos)
			return
		} else if err != mongo.ErrNoDocuments {
			// Si el error no es "ErrNoDocuments" (no encontrado), es un error real de la BD
			http.Error(w, "Error verificando la base de datos", http.StatusInternalServerError)
			return
		}

		// Si llegamos aquí, el error fue mongo.ErrNoDocuments, lo que significa que el correo está libre.
		// ... (Aquí va el resto de la conversión de datos: edad, peso, latitud, etc.)

		edad, _ := strconv.Atoi(r.FormValue("edad"))
		peso, _ := strconv.ParseFloat(r.FormValue("peso"), 64)
		latitud, _ := strconv.ParseFloat(r.FormValue("latitud"), 64)
		longitud, _ := strconv.ParseFloat(r.FormValue("longitud"), 64)

		tipoSangreCompleto := r.FormValue("tipoSangre")
		tipo := strings.TrimRight(tipoSangreCompleto, "+-")
		factorStr := "positivo"
		if strings.HasSuffix(tipoSangreCompleto, "-") {
			factorStr = "negativo"
		}

		passwordPlano := r.FormValue("password")
		hash, err := bcrypt.GenerateFromPassword([]byte(passwordPlano), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error procesando contraseña", 500)
			return
		}

		nuevoDonante := models.Donante{
			Nombre: r.FormValue("nombre"),

			Password: string(hash),

			Edad:       edad,
			Genero:     r.FormValue("genero"),
			Peso:       peso,
			TipoSangre: tipo,
			FactorRH:   factorStr,
			UbicacionGeografica: models.Ubicacion{
				Coordenadas: models.Coordenadas{Latitud: latitud, Longitud: longitud},
				Direccion:   r.FormValue("direccion"),
			},
			DatosContacto: models.DatosContacto{
				Correo:   email,
				Telefono: r.FormValue("telefono"),
				Whatsapp: r.FormValue("whatsapp"),
			},
			CondicionesMedicas: models.CondicionesMedicas{
				Padecimientos: r.Form["padecimientos"],
				Medicamentos:  r.Form["medicamentos"],
			},
			PreferenciasNotificacion: r.Form["notificaciones"],
			FechaRegistro:            time.Now().Format("2006-01-02"),
		}

		err = repository.CreateDonante(nuevoDonante)
		if err != nil {
			log.Println("Error creando donante:", err)
			http.Error(w, "Error al registrar la cuenta", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/donante/login", http.StatusSeeOther)
		return
	}

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

// =======================
// LOGIN
// =======================
func ShowDonorLogin(w http.ResponseWriter, r *http.Request) {
	// Estructura local para pasar errores a la vista de login
	type LoginData struct {
		ErrorMessage string
	}

	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/donante_login.html"))
		tmpl.Execute(w, LoginData{})
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		// 🔎 Buscar usuario en DB
		donante, err := repository.GetDonanteByEmail(email)
		if err != nil {
			// Usuario no encontrado
			tmpl := template.Must(template.ParseFiles("templates/donante_login.html"))
			tmpl.Execute(w, LoginData{ErrorMessage: "Correo o contraseña incorrectos"})
			return
		}

		// 🔐 Comparar password
		err = bcrypt.CompareHashAndPassword([]byte(donante.Password), []byte(password))
		if err != nil {
			// Contraseña incorrecta
			tmpl := template.Must(template.ParseFiles("templates/donante_login.html"))
			tmpl.Execute(w, LoginData{ErrorMessage: "Correo o contraseña incorrectos"})
			return
		}

		// TODO: aquí puedes crear sesión luego
		http.Redirect(w, r, "/donante/dashboard", http.StatusSeeOther)
	}
}
