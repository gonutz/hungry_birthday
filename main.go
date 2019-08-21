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
		heroDead           = "bug_small_dead.png"
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
		tongueOutTime      = 10
		frogMouthY         = 20 // y-offset from frog center to mouth
		frogAttackDist     = 400
		frogTongueCooldown = 70
	)

	// frogCountInCurrentLevel is increased whenever the player kills all frogs
	frogCountInCurrentLevel := 0
	// frogReactionTime is decreased for every won level
	frogReactionTime := 20

	heroX, heroY := 0.0, 0.0
	heroOffset := 0.0
	rotation := 0.0
	speed := 0.0
	rockTimer := 0
	var rocks []rock
	var frogs []frog
	lastHeroPositions := make([][2]float64, frogReactionTime)
	dead := false
	var markerTimer int

	newGame := func() {
		if len(frogs) == 0 {
			// player won the last game
			frogCountInCurrentLevel++
			frogReactionTime--
			if frogReactionTime < 3 {
				frogReactionTime = 3
			}
		}
		heroX, heroY = 0.0, 0.0
		heroOffset = 0.0
		rotation = 0.0
		speed = 0.0
		rockTimer = 0
		rocks = nil
		frogs = make([]frog, frogCountInCurrentLevel)
		for i := range frogs {
			frogs[i].x = float64(rand.Intn(2*windowW) - windowW)
			frogs[i].y = float64(rand.Intn(2*windowH) - windowH)
			frogs[i].life = 3
			frogs[i].tongueTimer = 0
		}
		for i := range lastHeroPositions {
			lastHeroPositions[i] = [2]float64{9999999, 9999999}
		}
		dead = false
		markerTimer = 0
	}
	newGame()

	draw.RunWindow("Hungry Birthday! Drop Rocks with CTRL", windowW, windowH, func(window draw.Window) {
		if window.WasKeyPressed(draw.KeyEscape) {
			window.Close()
		}

		if window.WasKeyPressed(draw.KeyF2) {
			newGame()
			return
		}

		if !dead {
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
					// add 30 to the frog's x center, this works better
					dx := frog.x + 30 - rocks[i].x
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
					frogs[i].tongueTimer = tongueOutTime
					frogs[i].tongueX = lastHeroPositions[0][0]
					frogs[i].tongueY = lastHeroPositions[0][1]
				}
			}
			if frogs[i].tongueTimer > 0 {
				dx := frogs[i].tongueX - heroX
				dy := frogs[i].tongueY - heroY
				if math.Hypot(dx, dy) < 20 {
					dead = true
					speed = 0
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

		markerTimer++

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
				// the tongue forms a line from (x,y) to (toX,toY)
				x := ix + frogSize/2
				y := iy + frogSize/2 + frogMouthY
				toX := x + round(f.tongueX-f.x) - heroHighDx
				toY := y + round(f.tongueY-f.y) - heroHighDy
				tongueW := round(math.Hypot(float64(toY-y), float64(toX-x)))
				angle := math.Atan2(float64(toY-y), float64(toX-x))
				deg := round(angle / math.Pi * 180)
				window.DrawImageFileTo(
					tongue,
					(x+toX-tongueW)/2,
					(y+toY-tongueH)/2,
					tongueW,
					tongueH,
					deg,
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
		if dead {
			hdx, hdy = 0, 0
		}
		if !dead {
			window.DrawImageFileRotated(heroShadow, round(hx+hdx), round(hy+hdy), round(rotation))
		}
		// draw objects
		for _, r := range rocks {
			x, y := hx-(heroX-r.x)+rockSize/2, hy-(heroY-r.y)+rockSize/2
			window.DrawImageFileRotated(rockImage, round(x-heroHighDx*r.height), round(y-heroHighDy*r.height), r.rotation)
		}
		heroImage := hero
		if dead {
			heroImage = heroDead
		}
		window.DrawImageFileRotated(heroImage, round(hx+hdx-heroHighDx), round(hy+hdy-heroHighDy), round(rotation))
		window.FillRect(0, 0, 220, 30, draw.White)
		if dead {
			window.DrawText("You're dead! Press F2", 15, 7, draw.Black)
		} else if len(frogs) == 0 {
			window.DrawText("You win! Press F2", 32, 7, draw.Black)
		} else {
			frogText := "Frogs"
			if len(frogs) == 1 {
				frogText = "Frog"
			}
			window.DrawText(fmt.Sprintf("%d %s Remaining", len(frogs), frogText), 32, 7, draw.Black)
		}

		// draw arrow indicating closest frog
		closest := -1
		minSqrDist := 0.0
		for i, f := range frogs {
			dist := (heroX-f.x)*(heroX-f.x) + (heroY-f.y)*(heroY-f.y)
			if closest == -1 || dist < minSqrDist {
				closest = i
				minSqrDist = dist
			}
		}
		if closest != -1 {
			f := frogs[closest]
			dx, dy := f.x-heroX, f.y-heroY
			if dx*dx+dy*dy > windowH*windowH/3 {
				scale := 1.0 / math.Hypot(dx, dy)
				dx *= scale
				dy *= scale
				move := 20 * math.Sin(float64(markerTimer)*0.1)
				x := round(windowW/2 - heroHighDx + dx*(70+move))
				y := round(windowH/2 - heroHighDy + dy*(70+move))
				window.DrawImageFileTo(frogImage, x, y, 20, 20, 0)
			}
		}
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
