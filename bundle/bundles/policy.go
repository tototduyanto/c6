package bundles

type Policy struct {
	UID      string   `json:"uid"`
	Name     string   `json:"name"`
	Filename string   `json:"filename"`
	Sequence int      `json:"sequence"`
	Tags     []string `json:"tags"`
}
