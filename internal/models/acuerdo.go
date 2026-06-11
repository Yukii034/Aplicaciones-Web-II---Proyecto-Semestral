package models

//Struct de Acuerdo con los campos necesarios para la creación de un acuerdo

type Acuerdo struct {
	ID                       int    `json:"id"`
	PublicacionID            int    `json:"publicacion_id"`
	IDOfertante              int    `json:"id_ofertante"`
	IDPublicador             int    `json:"id_publicador"`
	Tipo                     string `json:"tipo"`
	Estado                   string `json:"estado"`
	Mensaje_Inicial          string `json:"mensaje_inicial"`
	Motivo_Cancelacion       string `json:"motivo_cancelacion"`
	Confirmacion_Solicitante bool   `json:"confirmacion_solicitante"`
	CreatedAt                string `json:"created_at"`
	UpdatedAt                string `json:"updated_at"`
}
