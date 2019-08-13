package main

import (
	"math"
	"math/rand"

	"github.com/gonutz/prototype/draw"
)

func main() {
	const (
		windowW, windowH = 1080, 720
		background       = "grass_3.png"
		hero             = "bug_small.png"
		heroSize         = 100
		heroShadow       = "bug_small_shadow.png"
		dRotation        = 5
		acceleration     = 8
		rockDelay        = 5
		rockImage        = "rock_small.png"
		rockShadow       = "rock_small_shadow.png"
		rockSize         = 42
		frogImage        = "pepe_small.png"
		frogSize         = 135
	)

	heroX, heroY := 0.0, 0.0
	heroOffset := 0.0
	rotation := 0.0
	speed := 0.0
	rockTimer := 0
	var rocks []rock
	var frogs []frog
	frogs = make([]frog, 3)
	for i := range frogs {
		frogs[i].x = float64(rand.Intn(2*windowW) - windowW)
		frogs[i].y = float64(rand.Intn(2*windowH) - windowH)
		frogs[i].life = 3
	}

	draw.RunWindow("Hungry Birthday!", windowW, windowH, func(window draw.Window) {
		if window.WasKeyPressed(draw.KeyEscape) {
			window.Close()
		}

		if window.IsKeyDown(draw.KeyLeft) {
			rotation -= dRotation
		}
		if window.IsKeyDown(draw.KeyRight) {
			rotation += dRotation
		}

		if window.IsKeyDown(draw.KeyUp) {
			speed = acceleration
		} else {
			speed = 0
		}

		if window.WasKeyPressed(draw.KeyLeftControl) {
			if rockTimer <= 0 {
				rockTimer = rockDelay
				rocks = append(rocks, rock{heroX, heroY, rand.Intn(360), 1.0})
			}
		}
		rockTimer--
		for i := 0; i < len(rocks); i++ {
			rocks[i].height -= 0.05
			if rocks[i].height <= 0 {
				rocks[i].height = 0
			} else {
				rocks[i].rotation += 2
			}
		}

		if speed != 0 {
			dy, dx := math.Sincos(rotation / 180 * math.Pi)
			heroX += dx * speed
			heroY += dy * speed
		}
		heroOffset += 0.1

		// draw background in modulo space
		dx, dy := round(-heroX)%windowW, round(-heroY)%windowH
		if dx < 0 {
			dx += windowW
		}
		if dy < 0 {
			dy += windowH
		}
		window.DrawImageFile(background, dx, dy)
		window.DrawImageFile(background, dx-windowW, dy)
		window.DrawImageFile(background, dx-windowW, dy-windowH)
		window.DrawImageFile(background, dx, dy-windowH)
		// draw frogs
		const hx, hy = (windowW - heroSize) / 2, (windowH - heroSize) / 2
		for _, f := range frogs {
			x, y := hx-(heroX-f.x)+frogSize/2, hy-(heroY-f.y)+frogSize/2
			window.DrawImageFile(frogImage, round(x), round(y))
		}
		// draw shadows
		for _, r := range rocks {
			x, y := hx-(heroX-r.x)+rockSize/2, hy-(heroY-r.y)+rockSize/2
			window.DrawImageFileRotated(rockShadow, round(x), round(y), r.rotation)
		}
		hdy, hdx := math.Sincos(heroOffset)
		hdy *= 4
		hdx *= 4
		window.DrawImageFileRotated(heroShadow, round(hx+hdx), round(hy+hdy), round(rotation))
		// draw objects
		for _, r := range rocks {
			x, y := hx-(heroX-r.x)+rockSize/2, hy-(heroY-r.y)+rockSize/2
			window.DrawImageFileRotated(rockImage, round(x-50*r.height), round(y-70*r.height), r.rotation)
		}
		window.DrawImageFileRotated(hero, round(hx+hdx-50), round(hy+hdy-70), round(rotation))
	})
}

func round(x float64) int {
	if x < 0 {
		return int(x - 0.5)
	}
	return int(x + 0.5)
}

type rock struct {
	x, y     float64
	rotation int
	height   float64
}

type frog struct {
	x, y float64
	life int
}
