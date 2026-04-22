package handlers

import (
	"net/http"
)

// ShowDonorAuth sirve la página que pregunta: ¿Login o Registro?
func ShowDonorAuth(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/donante_auth.html")
}

// ShowDonorRegister sirve el formulario de registro
func ShowDonorRegister(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/donante_registro.html")
}

// ShowDonorLogin sirve el formulario de login
func ShowDonorLogin(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/donante_login.html")
}
