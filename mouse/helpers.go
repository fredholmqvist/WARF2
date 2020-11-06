package mouse

import (
	m "projects/games/warf2/worldmap"

	"github.com/hajimehoshi/ebiten"
)

// Most mouse modes share the same functionality:
// 	1. Handle the first click.
// 	2. Handle the dragging of mouse to select some range.

// These two functions below wrap this functionality
// and use lambdas to inject specific behaviour.

func firstClick(mp *m.Map, currentMousePos int, clickF func(), dragFs []func(*m.Map, int, int)) {
	if !hasClicked {
		clickF()
		setHasClicked(currentMousePos)
	}

	if startPoint >= 0 {
		mouseRange(mp, currentMousePos, startPoint, dragFs)
	}
}

func mouseRange(mp *m.Map, start, end int, fs []func(*m.Map, int, int)) {
	x1, y1, x2, y2 := tileRange(start, end)

	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			for _, f := range fs {
				f(mp, x, y)
			}
		}
	}

	previousStartPoint = startPoint
	previousEndPoint = end
}

func mousePos() int {
	mx, my := ebiten.CursorPosition()
	mx, my = mx/m.TileSize, my/m.TileSize

	return mx + (my * m.TilesW)
}

func setHasClicked(currentMousePos int) {
	startPoint = currentMousePos
	hasClicked = true
}

func unsetHasClicked() {
	startPoint = -1
	hasClicked = false
}
