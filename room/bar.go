package room

import (
	"math/rand"
	"sort"

	"github.com/Holmqvist1990/WARF2/entity"
	"github.com/Holmqvist1990/WARF2/globals"
	m "github.com/Holmqvist1990/WARF2/worldmap"
)

var barAutoID = 0

type Bar struct {
	ID    int
	tiles []int
}

func NewBar(mp *m.Map, x, y int) (*Bar, bool) {
	b := &Bar{}
	tiles := mp.FloodFillRoom(x, y, func() int { return m.BarFloor })
	if len(tiles) == 0 {
		return nil, false
	}
	sort.Ints(tiles)
	hasPlacedBar := false
	for _, idx := range tiles {
		tile := &mp.Tiles[idx]
		tile.Room = b
		if !hasPlacedBar {
			hasPlacedBar = placeBar(mp, tiles, idx)
		}
	}
	if !hasPlacedBar {
		for _, idx := range tiles {
			mp.Tiles[idx].Sprite = m.Ground
			mp.Tiles[idx].Room = nil
		}
		return nil, false
	}
	b.ID = barAutoID
	barAutoID++
	return b, true
}

func (b *Bar) GetID() int {
	return b.ID
}

func (b *Bar) String() string {
	return "Bar"
}

func (b *Bar) Update(mp *m.Map) {

}

func (b *Bar) Tiles() []int {
	return b.tiles
}

var bars = [][]int{
	{
		entity.BarDrinksLeft, entity.BarDrinksRight, entity.NoItem, entity.NoItem, entity.NoItem, entity.NoItem,
		entity.NoItem, entity.NoItem, entity.NoItem, entity.BarV, entity.BarStool, entity.NoItem,
		entity.BarLeft, entity.BarH, entity.BarH, entity.BarRight, entity.BarStool, entity.NoItem,
		entity.BarStool, entity.BarStool, entity.BarStool, entity.BarStool, entity.BarStool, entity.NoItem,
		entity.NoItem, entity.NoItem, entity.NoItem, entity.NoItem, entity.NoItem, entity.NoItem,
	},
	{
		entity.NoItem, entity.NoItem, entity.NoItem, entity.NoItem, entity.BarDrinksLeft, entity.BarDrinksRight,
		entity.NoItem, entity.BarStool, entity.BarV, entity.NoItem, entity.NoItem, entity.NoItem,
		entity.NoItem, entity.BarStool, entity.BarLeft, entity.BarH, entity.BarH, entity.BarRight,
		entity.NoItem, entity.BarStool, entity.BarStool, entity.BarStool, entity.BarStool, entity.BarStool,
		entity.NoItem, entity.NoItem, entity.NoItem, entity.NoItem, entity.NoItem, entity.NoItem,
	},
}

func placeBar(mp *m.Map, tiles []int, idx int) bool {
	////////////////
	// TODO
	// This is crap.
	////////////////
	placements := []int{}
	idxX, idxY := globals.IdxToXY(idx)
	width, height := 6, 5
	for y := idxY; y < idxY+height; y++ {
		for x := idxX; x < idxX+width; x++ {
			curr := globals.XYToIdx(x, y)
			if mp.Items[curr].Sprite != entity.NoItem {
				return false
			}
			if m.IsAnyWall(mp.Tiles[curr].Sprite) {
				return false
			}
			placements = append(placements, curr)
		}
	}
	randomBar := bars[rand.Intn(len(bars))]
	for i := 0; i < len(placements); i++ {
		mp.Items[placements[i]].Sprite = randomBar[i]
	}
	return true
}
