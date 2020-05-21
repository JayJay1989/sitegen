package release

import (
	"encoding/json"
	"github.com/hashicorp/go-version"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"sort"
	"time"
)

type ReleaseGroup struct {
	Name          string     `json:"name"`
	Source        string     `json:"source"`
	Date          *time.Time `json:"date"`
	Releases      []*Release `json:"-"`
	LatestRelease *Release   `json:"-"`
	Featured      bool       `json:"featured"`
}

func (g *ReleaseGroup) Load() error {
	data, err := ioutil.ReadFile(g.Source)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &g.Releases)
	if err != nil {
		return err
	}

	err = g.loadLatest()
	if err != nil {
		return err
	}

	return nil
}

func (g *ReleaseGroup) loadLatest() error {
	if len(g.Releases) == 0 {
		log.WithField("group", g.Name).Warn("Release group has no latest version")
		return nil
	}

	versions := make([]*version.Version, len(g.Releases))
	for i, release := range g.Releases {
		v, err := version.NewVersion("v" + release.Version)
		if err != nil {
			return err
		}

		versions[i] = v
	}

	sort.Sort(version.Collection(versions))

	latest := versions[len(versions)-1].Original()[1:]

	g.LatestRelease = g.findByVersion(latest)
	if g.LatestRelease == nil {
		log.WithField("group", g.Name).WithField("version", latest).Warn("Could not find version")
	}

	return nil
}

func (g *ReleaseGroup) findByVersion(version string) *Release {
	for _, release := range g.Releases {
		if release.Version == version {
			return release
		}
	}
	return nil
}
