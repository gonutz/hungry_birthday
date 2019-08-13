package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/gonutz/prototype/draw"
)

func main() {
	const (
		windowW, windowH   = 1080, 720
		background         = "grass_3.png"
		hero               = "bug_small.png"
		heroSize           = 100
		heroShadow         = "bug_small_shadow.png"
		heroHighDx         = 50
		heroHighDy         = 70
		dRotation          = 3.5
		acceleration       = 0.5
		maxSpeed           = 8
		rockDelay          = 5
		rockImage          = "rock_small.png"
		rockShadow         = "rock_small_shadow.png"
		rockSize           = 42
		frogImage          = "pepe_small.png"
		frogSize           = 135
		frogHurt           = "pepe_hurt.png"
		frogHurtSize       = 223
		tongue             = "tongue.png"
		tongueH            = 16
		frogMouthY         = 20 // y-offset from frog center to mouth
		frogReactionTime   = 30
		frogAttackDist     = 300
		frogTongueCooldown = 90
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
		frogs[i].tongueTimer = -99999
	}
	lastHeroPositions := make([][2]float64, frogReactionTime)
	for i := range lastHeroPositions {
		lastHeroPositions[i] = [2]float64{9999999, 9999999}
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
			speed += acceleration
			if speed > maxSpeed {
				speed = maxSpeed
			}
		} else {
			speed -= acceleration * 0.5
			if speed < 0 {
				speed = 0
			}
		}

		if window.WasKeyPressed(draw.KeyLeftControl) || window.WasKeyPressed(draw.KeyRightControl) {
			if rockTimer <= 0 {
				rockTimer = rockDelay
				rocks = append(rocks, rock{heroX, heroY, rand.Intn(360), 1.0})
			}
		}
		rockTimer--
		for i := 0; i < len(rocks); i++ {
			if rocks[i].height <= 0 {
				continue
			}
			rocks[i].height -= 0.05
			if rocks[i].height <= 0 {
				rocks[i].height = 0
				for f, frog := range frogs {
					dx := frog.x - rocks[i].x
					dy := frog.y - rocks[i].y
					if math.Hypot(dx, dy) < frogSize/2 {
						rocks = append(rocks[:i], rocks[i+1:]...)
						i--
						frogs[f].life--
						frogs[f].hurtTimer = 5
						if frogs[f].life <= 0 {
							frogs = append(frogs[:f], frogs[f+1:]...)
						}
						break
					}
				}
			} else {
				rocks[i].rotation += 2
			}
		}

		for i := range frogs {
			frogs[i].hurtTimer--
			frogs[i].tongueTimer--
			if frogs[i].tongueTimer < -frogTongueCooldown {
				frogs[i].tongueTimer = -frogTongueCooldown
			}
			if frogs[i].tongueTimer == -frogTongueCooldown {
				dx := frogs[i].x - lastHeroPositions[0][0]
				dy := frogs[i].y - lastHeroPositions[0][1]
				if math.Hypot(dx, dy) < frogAttackDist {
					frogs[i].tongueTimer = 5
					frogs[i].tongueX = lastHeroPositions[0][0]
					frogs[i].tongueY = lastHeroPositions[0][1]
				}
			}
		}

		if speed != 0 {
			dy, dx := math.Sincos(rotation / 180 * math.Pi)
			heroX += dx * speed
			heroY += dy * speed
		}
		heroOffset += 0.1

		copy(lastHeroPositions[0:], lastHeroPositions[1:])
		lastHeroPositions[len(lastHeroPositions)-1] = [2]float64{heroX, heroY}

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
			x, y := hx-(heroX-f.x), hy-(heroY-f.y)-frogSize/4
			img := frogImage
			ix, iy := round(x), round(y)
			if f.hurtTimer > 0 {
				img = frogHurt
				ix += (frogSize - frogHurtSize) / 2
				iy += (frogSize - frogHurtSize) / 2
			}
			window.DrawImageFile(img, ix, iy)

			// draw tongue
			if f.tongueTimer > 0 {
				window.DrawImageFileTo(
					tongue,
					ix+frogSize/2,
					iy+frogMouthY+(frogSize-tongueH)/2,
					500,
					tongueH,
					0,
				)
			}
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
			window.DrawImageFileRotated(rockImage, round(x-heroHighDx*r.height), round(y-heroHighDy*r.height), r.rotation)
		}
		window.DrawImageFileRotated(hero, round(hx+hdx-heroHighDx), round(hy+hdy-heroHighDy), round(rotation))
		window.FillRect(0, 0, 180, 30, draw.White)
		window.DrawText(fmt.Sprintf("%d Frogs Remaining", len(frogs)), 12, 7, draw.Black)
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
	x, y             float64
	life             int
	hurtTimer        int
	tongueTimer      int
	tongueX, tongueY float64
}
