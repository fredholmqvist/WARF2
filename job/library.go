package job

import (
	"projects/games/warf2/dwarf"
	"projects/games/warf2/item"
	"projects/games/warf2/room"
	"projects/games/warf2/worldmap"
	m "projects/games/warf2/worldmap"
)

type LibraryRead struct {
	dwarf        *dwarf.Dwarf
	destinations []int
	readingTime  int
}

func NewLibraryRead(destinations []int, readingTime int) *LibraryRead {
	return &LibraryRead{nil, destinations, readingTime}
}

func (l *LibraryRead) NeedsToBeRemoved(*m.Map) bool {
	return l.readingTime <= 0 || l.dwarf == nil
}

func (l *LibraryRead) Finish(*m.Map, *room.Service) {
	l.dwarf = nil
}

func (l *LibraryRead) PerformWork(m *m.Map, dwarves []*dwarf.Dwarf) bool {
	if shouldGetChair(m, l) {
		return getChair(m, l, dwarves)
	}
	// Still reading.
	if l.readingTime > 1 {
		l.readingTime--
		return unfinished
	}
	// Done!
	l.readingTime = 0
	return finished
}

func (l *LibraryRead) Priority() int {
	return 0
}

func (l *LibraryRead) GetWorker() *dwarf.Dwarf {
	return l.dwarf
}

func (l *LibraryRead) SetWorker(dw *dwarf.Dwarf) {
	l.dwarf = dw
}

func (l *LibraryRead) GetDestinations() []int {
	return l.destinations
}

func (l *LibraryRead) String() string {
	return "Library"
}

func shouldGetChair(m *worldmap.Map, l *LibraryRead) bool {
	return !item.IsChair(m.Items[l.dwarf.Idx].Sprite) &&
		l.dwarf.State != dwarf.WorkerMoving
}

func getChair(m *worldmap.Map, l *LibraryRead, dwarves []*dwarf.Dwarf) bool {
	dsts, ok := item.FindNearestChairs(m, l.destinations[0])
	if !ok {
		return unfinished
	}
	// Don't sit in occupied chair.
	target := -1
	for _, dst := range dsts {
		for _, dwarf := range dwarves {
			if dwarf.Idx == dst {
				continue
			}
		}
		target = dst
		break
	}
	// No available chairs.
	if target == -1 {
		return unfinished
	}
	l.destinations[0] = target
	l.dwarf.MoveTo(target, m)
	return unfinished
}
