/*
 _____               _____
|_   _| __ _   _  __|_   _|__ _ __ ___  _ __ ___
  | || '__| | | |/ _ \| |/ _ \ '_ ` _ \| '_ ` _ \
  | || |  | |_| |  __/| |  __/ | | | | | | | | | |
  |_||_|   \__,_|\___||_|\___|_| |_| |_|_| |_| |_|

 _____ __  __       ____                               __ _
|_   _|  \/  |     |  _ \ _ __ __ _  __ _  ___  _ __  / _| |_   _
  | | | |\/| |_____| | | | '__/ _` |/ _` |/ _ \| '_ \| |_| | | | |
  | | | |  | |_____| |_| | | | (_| | (_| | (_) | | | |  _| | |_| |
  |_| |_|  |_|     |____/|_|  \__,_|\__, |\___/|_| |_|_| |_|\__, |
                                    |___/                   |___/

@author TrueTemm
@link   https://github.com/TrueTemm
TM-Dragonfly Project
*/

package world

import (
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
)

var (
	Overworld overworld

	Nether nether

	End end
)

var dimensionReg = newDimensionRegistry(map[int]Dimension{
	0: Overworld,
	1: Nether,
	2: End,
})

func DimensionByID(id int) (Dimension, bool) {
	return dimensionReg.Lookup(id)
}

func DimensionID(dim Dimension) (int, bool) {
	return dimensionReg.LookupID(dim)
}

type dimensionRegistry struct {
	dimensions map[int]Dimension
	ids        map[Dimension]int
	custom     []DimensionRegistration
}

type DimensionRegistration struct {
	ID        int
	Name      string
	Dimension Dimension
}

func newDimensionRegistry(dim map[int]Dimension) *dimensionRegistry {
	ids := make(map[Dimension]int, len(dim))
	for k, v := range dim {
		ids[v] = k
	}
	return &dimensionRegistry{dimensions: dim, ids: ids}
}

func (reg *dimensionRegistry) Lookup(id int) (Dimension, bool) {
	dim, ok := reg.dimensions[id]
	if !ok {
		dim = Overworld
	}
	return dim, ok
}

func (reg *dimensionRegistry) LookupID(dim Dimension) (int, bool) {
	id, ok := reg.ids[dim]
	return id, ok
}

func (reg *dimensionRegistry) RegisterDimension(id int, name string, dim Dimension) error {
	if id < 1000 || id > math.MaxUint16 {
		return fmt.Errorf("custom dimension ID must be between 1000 and %d", math.MaxUint16)
	}
	if name == "" {
		return fmt.Errorf("custom dimension name must not be empty")
	}
	if dim == nil {
		return fmt.Errorf("custom dimension must not be nil")
	}
	r := dim.Range()
	if r.Min() > r.Max() {
		return fmt.Errorf("custom dimension range must not be empty")
	}
	if r.Min()%16 != 0 || (r.Max()+1)%16 != 0 {
		return fmt.Errorf("custom dimension range must align with 16-block sub-chunks")
	}
	if r.Min() < math.MinInt16 || r.Max() > math.MaxInt16 {
		return fmt.Errorf("custom dimension range must be between %d and %d", math.MinInt16, math.MaxInt16)
	}
	if _, ok := reg.dimensions[id]; ok {
		return fmt.Errorf("dimension ID %d is already registered", id)
	}
	if existing, ok := reg.ids[dim]; ok {
		return fmt.Errorf("dimension is already registered with ID %d", existing)
	}
	for _, existing := range reg.custom {
		if existing.Name == name {
			return fmt.Errorf("dimension name %q is already registered", name)
		}
	}
	reg.dimensions[id] = dim
	reg.ids[dim] = id
	reg.custom = append(reg.custom, DimensionRegistration{ID: id, Name: name, Dimension: dim})
	return nil
}

func RegisterDimension(id int, name string, dim Dimension) error {
	return dimensionReg.RegisterDimension(id, name, dim)
}

func CustomDimensions() []DimensionRegistration {
	return slices.Clone(dimensionReg.custom)
}

type (
	Dimension interface {
		Range() cube.Range
		WaterEvaporates() bool
		LavaSpreadDuration() time.Duration
		WeatherCycle() bool
		TimeCycle() bool
	}
	overworld struct{}
	nether    struct{}
	end       struct{}
)

func (overworld) Range() cube.Range                 { return cube.Range{-64, 319} }
func (overworld) WaterEvaporates() bool             { return false }
func (overworld) LavaSpreadDuration() time.Duration { return time.Second * 3 / 2 }
func (overworld) WeatherCycle() bool                { return true }
func (overworld) TimeCycle() bool                   { return true }
func (overworld) String() string                    { return "Overworld" }

func (nether) Range() cube.Range                 { return cube.Range{0, 127} }
func (nether) WaterEvaporates() bool             { return true }
func (nether) LavaSpreadDuration() time.Duration { return time.Second / 4 }
func (nether) WeatherCycle() bool                { return false }
func (nether) TimeCycle() bool                   { return false }
func (nether) String() string                    { return "Nether" }

func (end) Range() cube.Range                 { return cube.Range{0, 255} }
func (end) WaterEvaporates() bool             { return false }
func (end) LavaSpreadDuration() time.Duration { return time.Second * 3 / 2 }
func (end) WeatherCycle() bool                { return false }
func (end) TimeCycle() bool                   { return false }
func (end) String() string                    { return "End" }
