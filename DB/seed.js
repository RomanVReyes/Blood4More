

// Limpiar colecciones (opcional pero recomendable)
db.catalogo_padecimientos.drop()
db.catalogo_medicamentos.drop()
db.donantes.drop()
db.bancos_sangre.drop()
db.hospitales.drop()
db.solicitudes_sangre.drop()
db.notificaciones.drop()
db.donaciones.drop()
db.compatibilidad_sanguinea.drop()

// =======================
// CATÁLOGOS
// =======================

db.catalogo_padecimientos.insertMany([
  { _id: "pad_001", nombre: "Diabetes" },
  { _id: "pad_002", nombre: "Hipertensión" },
  { _id: "pad_003", nombre: "Asma" },
  { _id: "pad_004", nombre: "Enfermedad cardíaca" },
  { _id: "pad_005", nombre: "Anemia" }
])

db.catalogo_medicamentos.insertMany([
  { _id: "med_001", nombre: "Metformina" },
  { _id: "med_002", nombre: "Insulina" },
  { _id: "med_003", nombre: "Paracetamol" },
  { _id: "med_004", nombre: "Ibuprofeno" },
  { _id: "med_005", nombre: "Aspirina" }
])

// =======================
// DONANTES
// =======================

db.donantes.insertOne({
  _id: "donante_001",
  nombre: "Juan Perez",
  edad: 29,
  genero: "masculino",
  peso: 75,
  tipoSangre: "O+",
  factorRH: "positivo",

  ubicacionGeografica: {
    coordenadas: {
      latitud: 19.4326,
      longitud: -99.1332
    },
    direccion: "Calle Falsa 123",
    zona: "BajaCalifornia"
  },

  datosContacto: {
    telefono: "5551112233",
    correo: "juan@email.com",
    whatsapp: "5551112233"
  },

  condicionesMedicas: {
    padecimientos: ["pad_001"],
    medicamentos: ["med_001"]
  },

  fechaUltimaDonacion: new Date("2026-02-10"),

  elegibilidad: {
    estado: false,
    fechaProximaElegible: new Date("2026-04-07")
  },

  preferenciasNotificacion: ["whatsapp", "correo"],

  historialDonaciones: [
    {
      idDonacion: "don_1001",
      fecha: new Date("2026-02-10"),
      bancoReceptor: "banco_001",
      volumenExtraido: 450,
      estadoFinal: "aceptada"
    }
  ],

  fechaRegistro: new Date("2025-10-01")
})

// =======================
// BANCOS DE SANGRE
// =======================

db.bancos_sangre.insertOne({
  _id: "banco_001",
  nombreBanco: "Cruz Roja Central",
  tipoInstitucion: "cruz_roja",

  ubicacionGeografica: {
    coordenadas: {
      latitud: 19.4200,
      longitud: -99.1500
    },
    direccion: "Av Reforma 500",
    zona: "CDMX"
  },

  horarios: [
    {
      dia: "lunes",
      apertura: "08:00",
      cierre: "18:00"
    }
  ],

  capacidadAlmacenamiento: 5000,

  contactoEmergencia: {
    telefono: "5559998888",
    correo: "emergencias@banco.com"
  },

  personalResponsable: [
    {
      nombre: "Dr. Lopez",
      cargo: "Director",
      telefono: "5552223333",
      correo: "lopez@banco.com"
    }
  ],

  inventario: [
    {
      tipoSangre: "O+",
      unidadesDisponibles: 120,
      unidades: [
        {
          idUnidad: "u_001",
          fechaExtraccion: new Date("2026-03-01"),
          fechaCaducidad: new Date("2026-04-01"),
          estado: "disponible"
        }
      ]
    }
  ],

  umbralMinimo: {
    "O+": 50,
    "A+": 40
  }
})

// =======================
// HOSPITALES
// =======================

db.hospitales.insertOne({
  _id: "hospital_001",
  nombreHospital: "Hospital General",

  ubicacionGeografica: {
    coordenadas: {
      latitud: 19.40,
      longitud: -99.12
    },
    direccion: "Av Salud 200",
    zona: "CDMX"
  },

  datosContacto: {
    telefono: "5554443333",
    correo: "contacto@hospital.com"
  },

  medicos: [
    {
      idMedico: "med_001",
      nombre: "Dra. Ramirez",
      especialidad: "cirugia",
      telefono: "5556667777",
      correo: "ramirez@hospital.com"
    }
  ]
})

// =======================
// SOLICITUDES
// =======================

db.solicitudes_sangre.insertOne({
  _id: "sol_001",
  idHospital: "hospital_001",
  idMedicoSolicitante: "med_001",
  nombrePaciente: "Carlos Gomez",
  tipoSangreRequerido: "A+",
  unidadesNecesarias: 3,
  nivelUrgencia: "critica",
  fechaSolicitud: new Date("2026-04-10"),
  fechaLimite: new Date("2026-04-12"),
  motivoMedico: "cirugia urgente",
  tiposCompatibles: ["A+", "A-", "O+", "O-"],
  estadoSolicitud: "pendiente"
})

// =======================
// NOTIFICACIONES
// =======================

db.notificaciones.insertOne({
  _id: "notif_001",
  idDonante: "donante_001",
  idSolicitud: "sol_001",
  mensaje: "Se requiere donación urgente de sangre O+ cercana a tu ubicación",
  fechaEnvio: new Date("2026-04-14"),
  medio: "whatsapp",
  estado: "enviada"
})

// =======================
// DONACIONES
// =======================

db.donaciones.insertOne({
  _id: "don_1001",
  idDonante: "donante_001",
  idBanco: "banco_001",
  fecha: new Date("2026-02-10"),
  volumenExtraido: 450,

  resultadosPruebas: {
    hemoglobina: 14.2,
    presion: "120/80",
    hepatitis: "negativo",
    VIH: "negativo",
    sifilis: "negativo"
  },

  estadoFinal: "aceptada",
  idUnidadGenerada: "u_001"
})

// =======================
// COMPATIBILIDAD
// =======================

db.compatibilidad_sanguinea.insertMany([
  { tipo: "O-", compatibleCon: ["A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-"] },
  { tipo: "A+", compatibleCon: ["A+", "AB+"] },
  { tipo: "A-", compatibleCon: ["A-", "A+", "AB-", "AB+"] },
  { tipo: "B+", compatibleCon: ["B+", "AB+"] },
  { tipo: "B-", compatibleCon: ["B-", "B+", "AB-", "AB+"] },
  { tipo: "AB+", compatibleCon: ["AB+"] },
  { tipo: "AB-", compatibleCon: ["AB-", "AB+"] },
  { tipo: "O+", compatibleCon: ["O+", "A+", "B+", "AB+"] }
])

print("Base de datos cargada correctamente")