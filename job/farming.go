package job

import (
	"github.com/Holmqvist1990/WARF2/dwarf"
	"github.com/Holmqvist1990/WARF2/item"
	"github.com/Holmqvist1990/WARF2/room"
	m "github.com/Holmqvist1990/WARF2/worldmap"
)

type Farming struct {
	FarmID       int
	dwarf        *dwarf.Dwarf
	destinations []int
	path         []int
}

func NewFarming(farmID int, destinations []int) *Farming {
	return &Farming{farmID, nil, destinations, nil}
}

func (f *Farming) NeedsToBeRemoved(mp *m.Map) bool {
	return len(f.destinations) == 0 && f.path == nil
}

func (d *Farming) Finish(*m.Map, *room.Service) {
	if d.dwarf == nil {
		return
	}
	d.dwarf.SetToAvailable()
	d.dwarf = nil
}

// Ran on arrival.
func (f *Farming) PerformWork(mp *m.Map, dwarves []*dwarf.Dwarf) bool {
	if len(f.destinations) == 0 {
		return finished
	}
	if f.dwarf == nil {
		return unfinished
	}
	return f.moveDwarf(mp)
}

func (f *Farming) Priority() int {
	return 1
}

func (f *Farming) GetWorker() *dwarf.Dwarf {
	return f.dwarf
}

func (f *Farming) SetWorker(dw *dwarf.Dwarf) {
	f.dwarf = dw
}

func (f *Farming) GetDestinations() []int {
	return f.destinations
}

func (f *Farming) String() string {
	return "Farming"
}

func (f *Farming) moveDwarf(mp *m.Map) bool {
	currentIdx := f.destinations[len(f.destinations)-1]
	if f.dwarf.Idx == currentIdx {
		mp.Items[currentIdx].Sprite = item.Wheat
		f.destinations = f.destinations[:len(f.destinations)-1]
	}
	if f.NeedsToBeRemoved(mp) {
		return finished
	}
	if f.path != nil {
		f.moveAlongPath()
		return unfinished
	}
	nextIdx := f.getNextIdx()
	if nextIdx-f.dwarf.Idx == 1 {
		f.dwarf.Idx = nextIdx // Adjecent
	} else {
		f.getPath(mp, nextIdx) // Elsewhere
	}
	return unfinished
}

func (f *Farming) getPath(mp *m.Map, next int) {
	path, ok := f.dwarf.CreatePath(
		&mp.Tiles[f.dwarf.Idx],
		&mp.Tiles[next],
	)
	if !ok {
		return
	}
	f.path = path
}

func (f *Farming) moveAlongPath() {
	if len(f.path) == 0 {
		f.path = nil
		return
	}
	// Move indexes to current path index.
	f.dwarf.Idx = f.path[0]
	// Iterate path.
	f.path = f.path[1:]
}

func (f *Farming) getNextIdx() int {
	if len(f.destinations) == 1 {
		return f.destinations[0]
	}
	return f.destinations[len(f.destinations)-1]
}