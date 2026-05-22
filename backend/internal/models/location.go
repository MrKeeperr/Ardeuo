package models

type Pais struct {
	PaisID    int    `json:"pais_id"`
	Nombre   string `json:"nombre_pais"`
}

type Departamento struct {
	DepartamentoID int    `json:"departamento_id"`
	Nombre         string `json:"nombre_departamento"`
	PaisID         int    `json:"pais_id"`
}

type Municipio struct {
	MunicipioID    int    `json:"municipio_id"`
	Nombre         string `json:"nombre_municipio"`
	DepartamentoID int    `json:"departamento_id"`
}
