package release

import (
	"encoding/json"
	"errors"
	"github.com/gosimple/slug"
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
	StableRelease *Release   `json:"-"`
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

	for _, release := range g.Releases {
		release.Slug = slug.Make(release.Version)
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
		release.Group = g

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
		return errors.New("could not find latest version")
	}

	g.StableRelease = g.LatestRelease
	for i := len(versions) - 1; i >= 0; i-- {
		v := versions[i].Original()[1:]
		release := g.findByVersion(v)
		if release == nil {
			return errors.New("could not find version")
		}
		if release.Type == "release" {
			g.StableRelease = release
			break
		}
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
