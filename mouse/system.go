package mouse

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"

	m "projects/games/warf2/worldmap"
)

// System for handling
// all functionality by mouse.
type System struct {
	Mode Mode
}

// Mode enum for managing mouse action state.
type Mode int

// Mode enum.
const (
	Normal Mode = iota
	FloorTiles
	ResetFloor
)

// Handle all the mouse interactivity.
func (s *System) Handle(mp *m.Map) {
	s.mouseHover(mp)

	idx := mousePos()

	if idx < 0 || idx > m.TilesT {
		return
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		s.mouseClick(mp, idx)
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		endPoint = idx
		s.mouseUp(mp)
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		s.Mode = Normal
	}
}

func (s *System) mouseClick(mp *m.Map, currentMousePos int) {
	switch s.Mode {

	case Normal:
		noneMode(mp, currentMousePos)

	case FloorTiles:
		floorTileMode(mp, currentMousePos)

	case ResetFloor:
		resetFloorMode(mp, currentMousePos)

	default:
		fmt.Println("mouseClick: unknown MouseMode:", s.Mode)
	}
}

func (s *System) mouseUp(mp *m.Map) {
	if startPoint >= 0 {
		mouseRange(mp, startPoint, endPoint, []func(*m.Map, int, int){mouseUpSetWalls})
	}

	unsetHasClicked()
}

func (s *System) mouseHover(mp *m.Map) {
	switch s.Mode {
	default:
	}
}
