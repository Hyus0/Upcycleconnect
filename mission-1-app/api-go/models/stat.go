package models

type PlatformStats struct {
	Co2Evite       float64 `json:"co2_evite"`     
	ObjetsUpcycles int     `json:"objets_upcycles"`
	ArtisansActifs int     `json:"artisans_actifs"`
	SitesActifs    int     `json:"sites_actifs"`
}