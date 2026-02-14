package cuzk

import "time"

// CadastralArea represents a cadastral territory (katastrální území).
type CadastralArea struct {
	Code int    `json:"kod"`
	Name string `json:"nazev"`
}

// ReferencePoint contains coordinates in S-JTSK.
// CUZK API returns positive values.
type ReferencePoint struct {
	X float64 `json:"souradniceX"`
	Y float64 `json:"souradniceY"`
}

// Parcel represents a cadastral parcel (parcela).
// NOTE: struct fields are approximations based on CUZK docs — adjust after testing with real API.
type Parcel struct {
	ID             int64           `json:"id"`
	BaseNumber     int             `json:"kmenoveCislo"`
	Subdivision    *int            `json:"poddeleni,omitempty"`
	NumberingType  string          `json:"druhCislovani"`
	CadastralArea  CadastralArea   `json:"katastralniUzemi"`
	Area           int             `json:"vymera"`
	LandType       *string         `json:"druhPozemku,omitempty"`
	UsageType      *string         `json:"zpusobVyuziti,omitempty"`
	OwnershipSheet *string         `json:"cisloLV,omitempty"`
	ReferencePoint *ReferencePoint `json:"definicniBod,omitempty"`
}

// Building represents a building object (stavba).
type Building struct {
	ID            int64         `json:"id"`
	DescriptiveNo *int          `json:"cisloPopisne,omitempty"`
	EvidenceNo    *int          `json:"cisloEvidencni,omitempty"`
	BuildingType  string        `json:"typStavby"`
	MunicipalPart *string       `json:"castObce,omitempty"`
	CadastralArea CadastralArea `json:"katastralniUzemi"`
	UsageType     *string       `json:"zpusobVyuziti,omitempty"`
	ParcelNumber  *string       `json:"parcelneCislo,omitempty"`
}

// Unit represents a property unit such as an apartment (jednotka).
type Unit struct {
	ID               int64  `json:"id"`
	UnitNumber       string `json:"cisloJednotky"`
	UnitType         string `json:"typJednotky"`
	CommonPartsShare string `json:"podilNaSpolecnychCastech"`
	BuildingID       *int64 `json:"stavbaId,omitempty"`
}

// Proceeding represents a cadastral proceeding (řízení).
type Proceeding struct {
	ID             int64      `json:"id"`
	SequenceNumber int        `json:"poradoveCislo"`
	Year           int        `json:"rok"`
	Office         string     `json:"pracoviste"`
	Status         string     `json:"stavRizeni"`
	Type           string     `json:"typRizeni"`
	FilingDate     *time.Time `json:"datumPodani,omitempty"`
}

// ParcelSearchResponse wraps a list of parcels from a search query.
type ParcelSearchResponse struct {
	Parcels []Parcel `json:"parcely"`
	Total   int      `json:"total"`
}

// BuildingSearchResponse wraps a list of buildings from a search query.
type BuildingSearchResponse struct {
	Buildings []Building `json:"stavby"`
	Total     int        `json:"total"`
}

// UnitSearchResponse wraps a list of units from a search query.
type UnitSearchResponse struct {
	Units []Unit `json:"jednotky"`
	Total int    `json:"total"`
}

// NeighborParcelsResponse contains a list of neighboring parcels.
type NeighborParcelsResponse struct {
	ParcelID  int64    `json:"parcelaId"`
	Neighbors []Parcel `json:"sousedniParcely"`
}
