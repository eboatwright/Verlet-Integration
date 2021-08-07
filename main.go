package main


import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/eboatwright/tomohawk"

	"fmt"
)


type RopePoint struct {
	position    tomohawk.Vector2
	oldPosition tomohawk.Vector2
	radius      float64
	gravity     float64
}
func (r *RopePoint) ID() string { return "ropePoint" }

type Stick struct {
	length float64
	a      *RopePoint
	b      *RopePoint
}
func (s *Stick) ID() string { return "stick" }


type RopeSystem struct {
	ropePointImage *ebiten.Image
}

func (rs *RopeSystem) Update(entity tomohawk.Entity) {
	rp := entity.GetComponent("ropePoint").(*RopePoint)
	if entity.ID == "point1" {
		cursorX, cursorY := ebiten.CursorPosition()
		rp.position = tomohawk.Vector2 { float64(cursorX) - rp.radius, float64(cursorY) - rp.radius }
	} else {
		velocity := tomohawk.Vector2SubtractV(rp.position, rp.oldPosition)
		velocity.Y += rp.gravity
		rp.oldPosition = rp.position
		rp.position.AddV(velocity)
	}
}

func (rs *RopeSystem) Draw(entity tomohawk.Entity, screen *ebiten.Image, options *ebiten.DrawImageOptions) {
	rp := entity.GetComponent("ropePoint").(*RopePoint)
	tomohawk.SetDrawPosition(options, rp.position)
	screen.DrawImage(rs.ropePointImage, options)
}

func (rs *RopeSystem) GetRequirements() []string {
	return []string { "ropePoint" }
}

type StickSystem struct {}

func (ss *StickSystem) Update(entity tomohawk.Entity) {
	s := entity.GetComponent("stick").(*Stick)

	positionDifference := tomohawk.Vector2SubtractV(s.a.position, s.b.position)
	d := positionDifference.Magnitude()

	difference := (s.length - d) / d

	translatePosition := tomohawk.Vector2MultiplyF64(tomohawk.Vector2MultiplyF64(positionDifference, 0.55), difference)

	s.a.position.AddV(translatePosition)
	s.b.position.SubtractV(translatePosition)
}

func (ss *StickSystem) Draw(entity tomohawk.Entity, screen *ebiten.Image, options *ebiten.DrawImageOptions) {}

func (ss *StickSystem) GetRequirements() []string {
	return []string { "stick" }
}


func onStart() {
	gameScene := &tomohawk.Scene { ID: "game" }
	tomohawk.LoadScene(gameScene)

	numberOfPoints := 10

	points := []tomohawk.Entity {  }
	for i := 0; i < numberOfPoints; i++ {
		point := tomohawk.Entity { ID: "point" + fmt.Sprint(i + 1) }
		point.AddComponent(&RopePoint {
			position: tomohawk.Vector2 { 0, 0 },
			radius:   3,
			gravity:  0.5,
		})
		points = append(points, point)
		gameScene.AddEntity(point)
	}

	for j := 1; j < numberOfPoints; j++ {
		stick := tomohawk.Entity { ID: "stick" }
		stick.AddComponent(&Stick {
			length: 30,
			a:      points[j - 1].GetComponent("ropePoint").(*RopePoint),
			b:      points[j].GetComponent("ropePoint").(*RopePoint),
		})
		gameScene.AddEntity(stick)
	}

	gameScene.AddSystem(&RopeSystem { ropePointImage: tomohawk.LoadImage("data/img/circle.png") })
	gameScene.AddSystem(&StickSystem {  })
}


func main() {
	tomohawk.Start(960, 600, 1, "Verlet Integration", onStart)
}