package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

const size = 1024

func main() {
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// Draw each pixel
	for y := range size {
		for x := range size {
			img.Set(x, y, pixelColor(float64(x), float64(y)))
		}
	}

	f, _ := os.Create("build/appicon.png")
	defer f.Close()
	png.Encode(f, img)
}

func pixelColor(x, y float64) color.Color {
	cx, cy := float64(size)/2, float64(size)/2

	// Rounded rectangle background
	margin := 40.0
	radius := 180.0
	if !inRoundedRect(x, y, margin, margin, float64(size)-margin, float64(size)-margin, radius) {
		return color.RGBA{0, 0, 0, 0}
	}

	// Background: dark gradient
	t := (y - margin) / (float64(size) - 2*margin)
	bgR := lerp(20, 14, t)
	bgG := lerp(20, 14, t)
	bgB := lerp(32, 24, t)

	// Stopwatch body (main circle)
	bodyRadius := 280.0
	bodyCx, bodyCy := cx, cy+30

	// Outer ring
	distBody := dist(x, y, bodyCx, bodyCy)
	ringWidth := 32.0

	// Stopwatch top button (small rectangle) - drawn outside the circle
	btnW, btnH := 24.0, 50.0
	btnX, btnY := cx, bodyCy-bodyRadius-ringWidth/2-btnH/2+5
	if math.Abs(x-btnX) < btnW/2 && math.Abs(y-btnY) < btnH/2 {
		return color.RGBA{200, 200, 220, 255}
	}

	// Small horizontal bar at top of stopwatch
	barW, barH := 60.0, 16.0
	barY := bodyCy - bodyRadius - ringWidth/2 - btnH + 2
	if math.Abs(x-cx) < barW/2 && math.Abs(y-barY) < barH/2 {
		return color.RGBA{200, 200, 220, 255}
	}

	// Draw outer ring of stopwatch
	if distBody < bodyRadius+ringWidth/2 && distBody > bodyRadius-ringWidth/2 {
		aa := smoothstep(bodyRadius-ringWidth/2-1, bodyRadius-ringWidth/2+1, distBody) *
			smoothstep(bodyRadius+ringWidth/2+1, bodyRadius+ringWidth/2-1, distBody)
		return blendOver(
			color.RGBA{uint8(bgR), uint8(bgG), uint8(bgB), 255},
			color.RGBA{200, 200, 220, uint8(aa * 255)},
		)
	}

	// Everything below is inside the circle - check interior details before fill
	if distBody < bodyRadius-ringWidth/2 {
		innerBg := color.RGBA{uint8(bgR + 8), uint8(bgG + 8), uint8(bgB + 10), 255}

		// Clock hands - minute hand pointing to ~11 o'clock
		handAngle := -math.Pi/2 + math.Pi/6
		handLen := 200.0
		handWidth := 14.0
		handEndX := bodyCx + math.Cos(handAngle)*handLen
		handEndY := bodyCy + math.Sin(handAngle)*handLen
		if distToSegment(x, y, bodyCx, bodyCy, handEndX, handEndY) < handWidth/2 {
			return color.RGBA{99, 102, 241, 255} // Indigo accent
		}

		// Second hand - pointing to ~4 o'clock
		secAngle := -math.Pi/2 + math.Pi*2*0.6
		secLen := 220.0
		secWidth := 6.0
		secEndX := bodyCx + math.Cos(secAngle)*secLen
		secEndY := bodyCy + math.Sin(secAngle)*secLen
		if distToSegment(x, y, bodyCx, bodyCy, secEndX, secEndY) < secWidth/2 {
			return color.RGBA{48, 209, 88, 255} // Green accent
		}

		// Center dot
		if dist(x, y, bodyCx, bodyCy) < 18.0 {
			return color.RGBA{230, 230, 240, 255}
		}

		// Tick marks around the circle
		for i := range 12 {
			angle := float64(i) * math.Pi / 6
			innerR := bodyRadius - ringWidth/2 - 8
			outerR := bodyRadius - ringWidth/2 - 28
			tickW := 4.0
			if i%3 == 0 {
				tickW = 7.0
				outerR = bodyRadius - ringWidth/2 - 38
			}

			tickOutX := bodyCx + math.Cos(angle)*innerR
			tickOutY := bodyCy + math.Sin(angle)*innerR
			tickInX := bodyCx + math.Cos(angle)*outerR
			tickInY := bodyCy + math.Sin(angle)*outerR

			if distToSegment(x, y, tickInX, tickInY, tickOutX, tickOutY) < tickW/2 {
				return color.RGBA{120, 120, 140, 255}
			}
		}

		// Lightning bolt below center
		if inLightningBolt(x, y, bodyCx, bodyCy+90, 0.65) {
			return color.RGBA{255, 214, 10, 255} // Gold
		}

		return innerBg
	}

	return color.RGBA{uint8(bgR), uint8(bgG), uint8(bgB), 255}
}

func inLightningBolt(x, y, cx, cy, scale float64) bool {
	// Translate to bolt-local coords
	lx := (x - cx) / scale
	ly := (y - cy) / scale

	// Simple lightning bolt shape defined as a polygon
	// Top part going down-right, then jagging left, then down-right to point
	bolt := [][2]float64{
		{-8, -45},
		{15, -45},
		{2, -8},
		{18, -8},
		{-10, 45},
		{2, 5},
		{-15, 5},
	}

	return pointInPolygon(lx, ly, bolt)
}

func pointInPolygon(x, y float64, poly [][2]float64) bool {
	n := len(poly)
	inside := false

	j := n - 1
	for i := range n {
		yi, yj := poly[i][1], poly[j][1]
		xi, xj := poly[i][0], poly[j][0]

		if (yi > y) != (yj > y) && x < (xj-xi)*(y-yi)/(yj-yi)+xi {
			inside = !inside
		}

		j = i
	}

	return inside
}

func dist(x1, y1, x2, y2 float64) float64 {
	dx := x1 - x2
	dy := y1 - y2

	return math.Sqrt(dx*dx + dy*dy)
}

func distToSegment(px, py, x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	lenSq := dx*dx + dy*dy

	if lenSq == 0 {
		return dist(px, py, x1, y1)
	}

	t := ((px-x1)*dx + (py-y1)*dy) / lenSq
	t = math.Max(0, math.Min(1, t))

	projX := x1 + t*dx
	projY := y1 + t*dy

	return dist(px, py, projX, projY)
}

func inRoundedRect(x, y, left, top, right, bottom, radius float64) bool {
	// Check if point is in rounded rectangle
	if x < left || x > right || y < top || y > bottom {
		return false
	}

	// Check corners
	corners := [][2]float64{
		{left + radius, top + radius},
		{right - radius, top + radius},
		{left + radius, bottom - radius},
		{right - radius, bottom - radius},
	}

	for _, c := range corners {
		dx := math.Abs(x - c[0])
		dy := math.Abs(y - c[1])

		if x < left+radius && y < top+radius && c[0] == left+radius && c[1] == top+radius {
			if dx*dx+dy*dy > radius*radius {
				return false
			}
		}
		if x > right-radius && y < top+radius && c[0] == right-radius && c[1] == top+radius {
			if dx*dx+dy*dy > radius*radius {
				return false
			}
		}
		if x < left+radius && y > bottom-radius && c[0] == left+radius && c[1] == bottom-radius {
			if dx*dx+dy*dy > radius*radius {
				return false
			}
		}
		if x > right-radius && y > bottom-radius && c[0] == right-radius && c[1] == bottom-radius {
			if dx*dx+dy*dy > radius*radius {
				return false
			}
		}
	}

	return true
}

func smoothstep(edge0, edge1, x float64) float64 {
	t := (x - edge0) / (edge1 - edge0)
	t = math.Max(0, math.Min(1, t))

	return t * t * (3 - 2*t)
}

func lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

func blendOver(bg, fg color.RGBA) color.RGBA {
	a := float64(fg.A) / 255
	return color.RGBA{
		R: uint8(float64(fg.R)*a + float64(bg.R)*(1-a)),
		G: uint8(float64(fg.G)*a + float64(bg.G)*(1-a)),
		B: uint8(float64(fg.B)*a + float64(bg.B)*(1-a)),
		A: 255,
	}
}
