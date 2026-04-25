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

	"crypto/rand"
	"encoding/hex"
)

type DatosRegistro struct {
	Padecimientos []models.ItemCatalogo
	Medicamentos  []models.ItemCatalogo
	ErrorMessage  string
}

// Función auxiliar para crear IDs seguros
func generateSessionID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
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
			tmpl := template.Must(template.ParseFiles("templates/donante_login.html"))
			tmpl.Execute(w, LoginData{ErrorMessage: "Correo o contraseña incorrectos"})
			return
		}

		// 🔐 Comparar password
		err = bcrypt.CompareHashAndPassword([]byte(donante.Password), []byte(password))
		if err != nil {
			tmpl := template.Must(template.ParseFiles("templates/donante_login.html"))
			tmpl.Execute(w, LoginData{ErrorMessage: "Correo o contraseña incorrectos"})
			return
		}

		// 🎟️ CREAR LA SESIÓN EN MONGODB
		sessionID := generateSessionID()
		expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 días

		nuevaSesion := models.Sesion{
			ID:        sessionID,
			UserEmail: donante.DatosContacto.Correo,
			CreatedAt: time.Now(),
			ExpiresAt: expiresAt,
			IsActive:  true,
			IP:        r.RemoteAddr,
			UserAgent: r.UserAgent(),
		}

		err = repository.CreateSession(nuevaSesion)
		if err != nil {
			http.Error(w, "Error creando sesión", http.StatusInternalServerError)
			return
		}

		// 🍪 GUARDAR EL ID DE SESIÓN EN UNA COOKIE
		http.SetCookie(w, &http.Cookie{
			Name:     "roman_session",
			Value:    sessionID,
			Expires:  expiresAt,
			Path:     "/",
			HttpOnly: true,  // Importante: evita ataques XSS
			Secure:   false, // Cambiar a 'true' cuando uses HTTPS en producción
			SameSite: http.SameSiteLaxMode,
		})

		// Redirigir al dashboard
		http.Redirect(w, r, "/donante/dashboard", http.StatusSeeOther)
	}
}

// =======================
// DASHBOARD (Protegido)
// =======================
func ShowDonorDashboard(w http.ResponseWriter, r *http.Request) {
	// 1. Pedir el "gafete" (la cookie de sesión)
	cookie, err := r.Cookie("roman_session")
	if err != nil {
		// No trae gafete -> Pa' fuera (al login)
		http.Redirect(w, r, "/donante/login", http.StatusSeeOther)
		return
	}

	// 2. Revisar si la sesión existe y está activa en MongoDB
	sesion, err := repository.GetSession(cookie.Value)
	if err != nil || sesion.ExpiresAt.Before(time.Now()) {
		// Gafete falso o expirado -> Limpiar cookie y pa' fuera
		http.SetCookie(w, &http.Cookie{
			Name:   "roman_session",
			Value:  "",
			MaxAge: -1,
			Path:   "/",
		})
		http.Redirect(w, r, "/donante/login", http.StatusSeeOther)
		return
	}

	// 3. Buscar todos los datos del donante usando el email guardado en la sesión
	donante, err := repository.GetDonanteByEmail(sesion.UserEmail)
	if err != nil {
		http.Redirect(w, r, "/donante/logout", http.StatusSeeOther)
		return
	}

	// 4. Enviar los datos al HTML
	tmpl := template.Must(template.ParseFiles("templates/donante_dashboard.html"))
	tmpl.Execute(w, donante)
}

// =======================
// CERRAR SESIÓN (Logout)
// =======================
func LogoutDonante(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener la cookie actual
	cookie, err := r.Cookie("roman_session")
	if err == nil {
		// 2. Borrar la sesión de MongoDB para que ya no sirva
		_ = repository.DeleteSession(cookie.Value)
	}

	// 3. Destruir la cookie en el navegador del usuario
	http.SetCookie(w, &http.Cookie{
		Name:     "roman_session",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Fecha en el pasado
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	})

	// 4. Redirigir a la página de inicio o login
	http.Redirect(w, r, "/donante/login", http.StatusSeeOther)
}
