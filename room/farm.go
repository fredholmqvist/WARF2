package room

import (
	"sort"

	"github.com/Holmqvist1990/WARF2/item"
	m "github.com/Holmqvist1990/WARF2/worldmap"
)

type Farm struct {
	ID       int     // First tile.
	tileIdxs []int   // To be indexed against Worldmap.
	farmTile *m.Tile // Knows when farm has reached maturity.
}

func NewFarm(mp *m.Map, x, y int) *Farm {
	f := &Farm{}
	tiles := mp.FloodFillRoom(x, y, func() int { return m.None })
	if len(tiles) == 0 {
		return nil
	}
	sort.Sort(tiles)
	f.ID = tiles[0].Idx
	f.tileIdxs = tiles.ToIdxs()
	for _, t := range tiles {
		f.placeFarm(mp, t)
		if f.farmTile != nil {
			continue
		}
		if !item.IsFarm(mp.Items[t.Idx].Sprite) {
			continue
		}
		f.farmTile = &mp.Items[t.Idx]
	}
	return f
}

func (f *Farm) Update(mp *m.Map) {
	for _, tIdx := range f.tileIdxs {
		tile := &mp.Items[tIdx]
		if tile.Sprite == 0 {
			continue
		}
		switch tile.Sprite {
		case item.FarmSingleEmpty:
			tile.Sprite = item.FarmSingleWheat1
		case item.FarmLeftEmpty:
			tile.Sprite = item.FarmLeftWheat1
		case item.FarmMiddleEmpty:
			tile.Sprite = item.FarmMiddleWheat1
		case item.FarmRightEmpty:
			tile.Sprite = item.FarmRightWheat1

		case item.FarmSingleWheat1:
			tile.Sprite = item.FarmSingleWheat2
		case item.FarmLeftWheat1:
			tile.Sprite = item.FarmLeftWheat2
		case item.FarmMiddleWheat1:
			tile.Sprite = item.FarmMiddleWheat2
		case item.FarmRightWheat1:
			tile.Sprite = item.FarmRightWheat2

		case item.FarmSingleWheat2:
			tile.Sprite = item.FarmSingleWheat3
		case item.FarmLeftWheat2:
			tile.Sprite = item.FarmLeftWheat3
		case item.FarmMiddleWheat2:
			tile.Sprite = item.FarmMiddleWheat3
		case item.FarmRightWheat2:
			tile.Sprite = item.FarmRightWheat3

		case item.FarmSingleWheat3:
			tile.Sprite = item.FarmSingleWheat4
		case item.FarmLeftWheat3:
			tile.Sprite = item.FarmLeftWheat4
		case item.FarmMiddleWheat3:
			tile.Sprite = item.FarmMiddleWheat4
		case item.FarmRightWheat3:
			tile.Sprite = item.FarmRightWheat4
		}
	}
}

func (f *Farm) ShouldHarvest(mp *m.Map) ([]int, bool) {
	if !item.IsFarmHarvestable(f.farmTile.Sprite) {
		return nil, false
	}
	return f.GetHarvestIdxs(mp), true
}

func (f *Farm) GetHarvestIdxs(mp *m.Map) []int {
	idxs := []int{}
	for _, tIdx := range f.tileIdxs {
		if !item.IsFarm(mp.Items[tIdx].Sprite) {
			continue
		}
		idxs = append(idxs, tIdx)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(idxs)))
	return idxs
}

func (f *Farm) placeFarm(mp *m.Map, t m.Tile) {
	skips := []bool{
		m.IsAnyWall(mp.OneTileLeft(t.Idx).Sprite),
		m.IsAnyWall(mp.OneTileRight(t.Idx).Sprite),
		m.IsAnyWall(mp.OneTileUp(t.Idx).Sprite),
		m.IsAnyWall(mp.OneTileUpLeft(t.Idx).Sprite),
		m.IsAnyWall(mp.OneTileUpRight(t.Idx).Sprite),
		m.IsAnyWall(mp.OneTileDown(t.Idx).Sprite),
		m.IsAnyWall(mp.OneTileDownLeft(t.Idx).Sprite),
		m.IsAnyWall(mp.OneTileDownRight(t.Idx).Sprite),
	}
	for _, skip := range skips {
		if skip {
			return
		}
	}
	if item.IsFarmSingle(mp.Items[m.OneTileLeft(t.Idx)].Sprite) {
		mp.Items[t.Idx-1].Sprite = item.FarmLeftEmpty
		item.Place(mp, t.X, t.Y, item.FarmRightEmpty)
		return
	}
	if item.IsFarmRight(mp.Items[m.OneTileLeft(t.Idx)].Sprite) {
		mp.Items[t.Idx-1].Sprite = item.FarmMiddleEmpty
		item.Place(mp, t.X, t.Y, item.FarmRightEmpty)
		return
	}
	item.Place(mp, t.X, t.Y, item.FarmSingleEmpty)
	return
}