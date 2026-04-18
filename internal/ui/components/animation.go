package components

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// AnimTickMsg signals the animation should advance one frame.
type AnimTickMsg time.Time

// AnimTickCmd returns a tea.Cmd that fires an AnimTickMsg after the given
// interval, creating a smooth frame-by-frame animation loop.
func AnimTickCmd(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return AnimTickMsg(t)
	})
}

// Animation cycles through ASCII frames of a sheep jumping over a fence.
type Animation struct {
	frame  int
	frames []string
}

// NewAnimation constructs the animation with all pre-rendered frames.
func NewAnimation() Animation {
	return Animation{frames: buildFrames()}
}

// Frame returns the current animation frame string.
func (a *Animation) Frame() string {
	if len(a.frames) == 0 {
		return ""
	}
	return a.frames[a.frame%len(a.frames)]
}

// NextFrame advances to the next frame in the loop.
func (a *Animation) NextFrame() {
	if len(a.frames) > 0 {
		a.frame = (a.frame + 1) % len(a.frames)
	}
}

// ─── Programmatic Frame Builder ─────────────────────────────────────────────

const (
	sceneW  = 58
	sceneH  = 10
	fenceX  = 30
	nFrames = 12
)

type point struct{ x, y int }

// buildFrames generates 12 smooth animation frames of a sheep leaping
// over a fence under a starry night sky with twinkling stars.
func buildFrames() []string {
	path := []point{
		{5, 0}, {11, 0}, {17, 0}, // walking toward fence
		{22, 1}, {25, 2}, {28, 3}, // jumping up
		{31, 3},          // peak (over fence)
		{34, 2}, {38, 1}, // descending
		{43, 0}, {50, 0}, // walking away
		{-1, 0}, // pause (empty scene)
	}

	// Two star-field patterns that alternate each frame for twinkling.
	skyA := [2]string{
		"   .  *      .     +      C    .    +   .  *   .     ",
		"     +   .    *    .      .     *  .    +    .   .   ",
	}
	skyB := [2]string{
		"     .    *    .  +    C     .   .  * .    .    +    ",
		"   +    .   *      .    .       +    .  *   .  .    ",
	}

	frames := make([]string, nFrames)

	for i, p := range path {
		// Blank canvas.
		var canvas [sceneH][]rune
		for r := range canvas {
			canvas[r] = make([]rune, sceneW)
			for c := range canvas[r] {
				canvas[r][c] = ' '
			}
		}

		// 1 ── Sky (rows 0-1), alternate for twinkle.
		sky := skyA
		if i%2 == 1 {
			sky = skyB
		}
		placeRunes(canvas[:], 0, 0, sky[0])
		placeRunes(canvas[:], 1, 0, sky[1])

		// 2 ── Fence (rows 7-8).
		placeRunes(canvas[:], 7, fenceX-1, "|-|")
		placeRunes(canvas[:], 8, fenceX-1, "| |")

		// 3 ── Ground (row 9).
		for c := range canvas[9] {
			canvas[9][c] = '-'
		}
		for _, t := range []int{3, 10, 17, 24, 37, 44, 52} {
			if t < sceneW {
				canvas[9][t] = 'v'
			}
		}
		for _, t := range []int{7, 21, 41, 55} {
			if t < sceneW {
				canvas[9][t] = '*'
			}
		}

		// 4 ── Sheep.
		if p.x >= 0 {
			baseRow := 8 - p.y // hooves row

			var wool, face, body, feet string
			switch {
			case p.y == 0:
				wool = " ,@@@. "
				face = "( o.o )"
				body = " /|  |\\"
				feet = "  d  b "
			case p.y >= 3:
				wool = " ,@@@. "
				face = "( >w< )"
				body = "  /  \\" + " "
				feet = "       "
			default:
				wool = " ,@@@. "
				face = "( o^o )"
				body = "  /  \\" + " "
				feet = "       "
			}

			parts := []string{wool, face, body, feet}
			for j, part := range parts {
				row := baseRow - 3 + j
				placeRunes(canvas[:], row, p.x-3, part)
			}

			// Sparkle near sheep at peak.
			if p.y >= 3 {
				placeRunes(canvas[:], baseRow-4, p.x+5, "*")
			}
			// Dust on landing.
			if p.y == 1 && i > 6 {
				placeRunes(canvas[:], baseRow+1, p.x+4, "~")
			}
		}

		// Render canvas to string.
		var b strings.Builder
		for r, row := range canvas {
			b.WriteString(string(row))
			if r < sceneH-1 {
				b.WriteByte('\n')
			}
		}
		frames[i] = b.String()
	}

	return frames
}

// placeRunes writes a string onto a canvas row at the given column,
// clipping characters that fall outside the canvas bounds.
func placeRunes(canvas [][]rune, row, col int, s string) {
	if row < 0 || row >= len(canvas) {
		return
	}
	for i, r := range []rune(s) {
		c := col + i
		if c >= 0 && c < len(canvas[row]) {
			canvas[row][c] = r
		}
	}
}
