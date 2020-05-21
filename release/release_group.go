package release

import (
	"encoding/json"
	"io/ioutil"
)

type ReleaseGroup struct {
	Name     string     `json:"name"`
	Source   string     `json:"source"`
	Releases []*Release `json:"-"`
}

func (g *ReleaseGroup) Load() error {
	data, err := ioutil.ReadFile(g.Source)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &g.Releases)
}
