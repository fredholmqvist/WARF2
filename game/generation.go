package game

import (
	"fmt"
	"math/rand"
	d "projects/games/warf2/dwarf"
	j "projects/games/warf2/jobservice"
	"projects/games/warf2/mouse"
	rail "projects/games/warf2/railservice"
	"projects/games/warf2/ui"
	u "projects/games/warf2/ui"
	m "projects/games/warf2/worldmap"
)

func GenerateGame(dwarves int, worldmap *m.Map) Game {
	game := Game{
		WorldMap:     *worldmap,
		JobService:   j.JobService{Map: worldmap},
		DwarfService: d.NewService(),
		RailService:  rail.RailService{Map: worldmap},

		time:        Time{Frame: 1},
		mouseSystem: mouse.System{},
		ui: u.UI{
			MouseMode: u.NewMouseOverlay(),
			MainMenu:  ui.NewMainMenu(),
		},
	}
	for i := 0; i < dwarves; i++ {
		addDwarfToGame(&game, game.DwarfService.RandomName())
	}
	return game
}

func emptyMap() *m.Map {
	return m.New()
}

func placeNewDwarf(mp m.Map, name string) (*d.Dwarf, bool) {
	var availableSpots []int
	for i := range mp.Tiles {
		if m.IsGround(mp.Tiles[i].Sprite) {
			availableSpots = append(availableSpots, mp.Tiles[i].Idx)
		}
	}
	if len(availableSpots) == 0 {
		fmt.Println("generation.go:placeNewDwarf: no available spaces")
		return nil, false
	}
	startingPosition := availableSpots[rand.Intn(len(availableSpots))]
	return d.New(startingPosition, name), true
}

func addDwarfToGame(g *Game, name string) {
	dwarf, ok := placeNewDwarf(g.WorldMap, name)
	if !ok {
		fmt.Println("generation.go:addDwarfToGame: dwarf was nil")
		return
	}
	g.JobService.Workers = append(g.JobService.Workers, dwarf)
}
