package job

import (
	"fmt"

	"github.com/Holmqvist1990/WARF2/dwarf"
	"github.com/Holmqvist1990/WARF2/globals"
	"github.com/Holmqvist1990/WARF2/resource"
	"github.com/Holmqvist1990/WARF2/room"
	m "github.com/Holmqvist1990/WARF2/worldmap"
)

type Carrying struct {
	resource        resource.Resource
	dwarf           *dwarf.Dwarf
	destinations    []int
	goalDestination int
	storageIdx      int
	sprite          int
	path            []int
	prev            int
}

func NewCarrying(destinations []int, r resource.Resource, storageIdx int, goalDestination, sprite int) *Carrying {
	return &Carrying{
		resource:        r,
		dwarf:           nil,
		destinations:    destinations,
		goalDestination: goalDestination,
		storageIdx:      storageIdx,
		sprite:          sprite,
		path:            nil,
		prev:            0,
	}
}

func (c *Carrying) NeedsToBeRemoved(mp *m.Map, r *room.Service) bool {
	return c.path != nil && len(c.path) == 0
}

func (c *Carrying) Finish(mp *m.Map, s *room.Service) {
	if c.dwarf == nil {
		return
	}
	defer func() {
		c.dwarf.SetToAvailable()
		c.dwarf = nil
	}()
	// Storage was deleted by user.
	if c.storageIdx > len(s.Storages)-1 {
		return
	}
	if c.sprite == m.None {
		return
	}
	dropIdx, ok := s.Storages[c.storageIdx].AddItem(c.dwarf.Idx, 1, c.resource)
	if !ok {
		////////
		// TODO
		// Yeah.
		////////
		fmt.Println("Carrying: Finish: Couldn't find storage tile.",
			"Ignoring item (forever lost!).")
		return
	}
	mp.Items[dropIdx].Sprite = c.sprite
}

func (c *Carrying) PerformWork(mp *m.Map, dwarves []*dwarf.Dwarf, rs *room.Service) bool {
	if storageMissingOrFull(c, rs) {
		// Try again with
		// new storage.
		return finished
	}
	if c.path == nil {
		///////////////////////////////////
		// TODO
		// Item is no longer there, abort.
		// What should we actually do here?
		///////////////////////////////////
		if !globals.IsCarriable(mp.Items[c.dwarf.Idx].Sprite) {
			c.path = []int{}
			return finished
		}
		setupPath(c, mp)
		return unfinished
	}
	if len(c.path) == 0 {
		return finished
	}
	moveAlongPath(c, mp)
	return unfinished
}

func (c *Carrying) GetWorker() *dwarf.Dwarf {
	return c.dwarf
}

func (c *Carrying) SetWorker(dw *dwarf.Dwarf) {
	c.dwarf = dw
}

func (c *Carrying) GetDestinations() []int {
	return c.destinations
}

func (c *Carrying) HasInternalMove() bool {
	return false
}

func (c *Carrying) String() string {
	return "Carrying"
}

func setupPath(c *Carrying, mp *m.Map) {
	mp.Items[c.dwarf.Idx].Sprite = 0
	mp.Items[c.dwarf.Idx].Resource = 0
	c.prev = c.dwarf.Idx
	c.destinations[0] = c.dwarf.Idx
	path, ok := c.dwarf.CreatePath(
		&mp.Tiles[c.dwarf.Idx],
		&mp.Tiles[c.goalDestination],
	)
	if !ok {
		return
	}
	c.path = path
}

func moveAlongPath(c *Carrying, mp *m.Map) {
	// Move indexes to current path index.
	c.dwarf.Idx = c.path[0]
	c.destinations[0] = c.path[0]
	c.prev = c.path[0]
	// Iterate path.
	c.path = c.path[1:]
}

func storageMissingOrFull(c *Carrying, rs *room.Service) bool {
	if len(rs.Storages)-1 < c.storageIdx {
		c.path = []int{}
		return true
	}
	storage := rs.Storages[c.storageIdx]
	if !storage.HasSpace(c.resource) {
		c.path = []int{}
		return true
	}
	return false
}
