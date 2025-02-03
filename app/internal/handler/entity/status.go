package entity

type Status struct {
	Status    string `json:"status"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	Build     string `json:"build"`
	GoVersion string `json:"goVersion"`
}
