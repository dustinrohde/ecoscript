package main

import (
	"io/ioutil"
	"log"
	"strings"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Mapfile is the Settings object that a Mapfile will be marshalled into.
type Mapfile struct {
	Defaults struct {
		EmptyTile     string `mapstructure:"empty_tile"`
		DisplayLegend bool   `mapstructure:"display_legend"`
	} `mapstructure:"defaults"`

	Atlas struct {
		Map struct {
			grid [][]string
			Text string `mapstructure:"text"`
			File string `mapstructure:"file"`
		} `mapstructure:"map"`

		Legend map[string]string `mapstructure:"legend"`
	} `mapstructure:"atlas"`

	Ecology struct {
		Classes []Class `mapstructure:"classes"`
	} `mapstructure:"ecology"`

	Organisms map[string]OrganismData `mapstructure:"organisms"`
}

type OrganismData struct {
	*Organism
	Abilities []*Ability
}

// ParseMapfile reads and parses a Mapfile given a path.
func ParseMapfile(path string) (mapfile Mapfile) {
	v := viper.New()
	v.SetTypeByDefaultValue(true)

	v.SetDefault("defaults.empty_tile", ".")
	v.SetDefault("defaults.display_legend", false)

	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatal(errors.Wrapf(err, "error reading Mapfile '%s'", path))
	}

	err = v.Unmarshal(&mapfile)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error unmarshaling config"))
	}

	err = mapfile.Validate()
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (m Mapfile) Validate() (err error) {
	var mapText string
	mapRawGiven := len(m.Atlas.Map.Text) > 0
	mapLinkGiven := len(m.Atlas.Map.File) > 0
	if !(mapRawGiven || mapLinkGiven) {
		return errors.New("one of ``atlas.map.text`` or ``atlas.map.file`` must be present")
	}
	if mapRawGiven && mapLinkGiven {
		return errors.New("``atlas.map.text`` and ``atlas.map.file`` cannot both be present")
	}
	if mapRawGiven {
		mapText = m.Atlas.Map.Text
	} else {
		bytes, err := ioutil.ReadFile(m.Atlas.Map.File)
		if err != nil {
			return errors.Wrap(err, "error reading ``atlas.map.link``")
		}
		mapText = string(bytes)
	}
	m.Atlas.Map.grid = gridify(mapText)

	if len(m.Atlas.Legend) == 0 {
		return errors.New("``atlas.legend`` must have at least one entry")
	}
	err = m.validateMap(mapText)
	if err != nil {
		return
	}

	if len(m.Organisms) == 0 {
		return errors.New("``organisms`` must have at least one entry")
	}
	err = m.validateLegend()
	if err != nil {
		return
	}

	err = m.validateEcology()
	if err != nil {
		return
	}

	return
}

func gridify(s string) [][]string {
	rows := strings.Split(strings.TrimSpace(s), "\n")
	grid := make([][]string, len(rows))
	for y, row := range rows {
		grid[y] = strings.Split(strings.TrimSpace(row), "")
	}
	return grid
}

func (m Mapfile) validateMap(mapText string) error {
	for _, row := range m.Atlas.Map.grid {
		for _, char := range row {
			if char == m.Defaults.EmptyTile {
				continue
			}
			_, ok := m.Atlas.Legend[char]
			if !ok {
				return errors.Errorf("map symbol '%s' not found in ``atlas.legend``", char)
			}
		}

	}
	return nil
}

func (m Mapfile) validateLegend() error {
	for _, key := range m.Atlas.Legend {
		_, ok := m.Organisms[key]
		if !ok {
			return errors.Errorf("'%s' is referenced in ``atlas.legend``, but no entry is found in ``organisms``")
		}
	}
	return nil
}

func (m Mapfile) validateEcology() error {
	var result error
	classes := m.Ecology.Classes
	if classes != nil && len(classes) > 0 {
		for _, class := range classes {
			if err := vStringMinLen(string(class), 2, "classes"); err != nil {
				result = multierror.Append(result, err)
			}
		}
	}
	return result
}

func (m Mapfile) validateOrganisms() error {
	var result error
	for _, organism := range m.Organisms {
		if err := vStringMinLen(organism.Stats.Name, 2, "name"); err != nil {
			result = multierror.Append(result, err)
		}
		if err := vIntMinVal(organism.Stats.Energy, 1, "energy"); err != nil {
			result = multierror.Append(result, err)
		}
		if err := vIntMinVal(organism.Stats.Size, 1, "size"); err != nil {
			result = multierror.Append(result, err)
		}
		if err := vIntMinVal(organism.Stats.Mass, 1, "mass"); err != nil {
			result = multierror.Append(result, err)
		}
		if organism.Classes != nil && len(organism.Classes) > 0 {
			for _, orgClass := range organism.Classes {
				ok := false
				for _, ecoClass := range m.Ecology.Classes {
					if orgClass == ecoClass {
						ok = true
						break
					}
				}
				if !ok {
					err := errors.Errorf("class \"%s\" not found in ``ecology.classes``", orgClass)
					result = multierror.Append(result, err)
				}
			}
		}
	}
	return result
}

func vStringMinLen(val string, min int, key string) (err error) {
	if len(val) < min {
		err = errors.Errorf("organism attribute \"%s\" must have %d or more characters", key, min)
	}
	return
}

func vIntMinVal(val int, min int, key string) (err error) {
	if val < min {
		err = errors.Errorf("organism attribute \"%s\" must be %d or greater", key, min)
	}
	return
}

func (m Mapfile) ToWorld() *World {
	return nil
}
