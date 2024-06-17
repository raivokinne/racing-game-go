package main

import (
	"fmt"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 600
	playerSpeed  = 5.0
	playerAccel  = 0.1
)

type Player struct {
	Position     rl.Vector2
	Size         rl.Vector2
	Speed        float32
	Acceleration float32
	Color        rl.Color
	Score        int
	Lives        int
}

type Obstacle struct {
	Position rl.Vector2
	Size     rl.Vector2
	Color    rl.Color
}

type Track struct {
	Obstacles  []Obstacle
	Color      rl.Color
	Size       rl.Vector2
	FinishLine rl.Rectangle
	Length     float32
}

type Game struct {
	Player   Player
	Track    Track
	Over     bool
	Paused   bool
	LapCount int
	Objects  int
}

func main() {
	game := Game{}
	game.init()
	rl.InitWindow(screenWidth, screenHeight, "Racing Game")

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		game.update()
		game.draw()
	}

	rl.CloseWindow()
}

func (g *Game) init() {
	g.Player = Player{
		Position:     rl.Vector2{X: screenWidth / 2, Y: screenHeight - 50},
		Size:         rl.Vector2{X: 50, Y: 30},
		Speed:        playerSpeed,
		Acceleration: playerAccel,
		Color:        rl.Blue,
		Score:        0,
		Lives:        3,
	}

	g.Objects = 0

	g.Track = Track{
		Obstacles: []Obstacle{
			{Position: rl.Vector2{X: 100, Y: 300}, Size: rl.Vector2{X: 50, Y: 50}, Color: rl.Red},
			{Position: rl.Vector2{X: 300, Y: 200}, Size: rl.Vector2{X: 50, Y: 50}, Color: rl.Red},
		},
		Color:      rl.DarkGray,
		Size:       rl.Vector2{X: screenWidth, Y: screenHeight},
		FinishLine: rl.NewRectangle(screenWidth/2-150, 50, 300, 10),
	}

	g.Over = false
	g.Paused = false
	g.LapCount = 0
}

func (g *Game) update() {
	if !g.Over {
		if rl.IsKeyPressed(rl.KeySpace) {
			g.Paused = !g.Paused
		}

		if g.Paused {
			return
		}

		if rl.IsKeyDown(rl.KeyRight) {
			g.Player.Position.X += g.Player.Speed
		}
		if rl.IsKeyDown(rl.KeyLeft) {
			g.Player.Position.X -= g.Player.Speed
		}
		if rl.IsKeyDown(rl.KeyUp) {
			g.Player.Position.Y -= g.Player.Speed
		}
		if rl.IsKeyDown(rl.KeyDown) {
			g.Player.Position.Y += g.Player.Speed
		}

		for _, obstacle := range g.Track.Obstacles {
			if checkCollision(g.Player, obstacle) {
				g.Player.Lives--
				if g.Player.Lives <= 0 {
					g.Over = true
				}

				g.Player.Position = rl.Vector2{X: screenWidth / 2, Y: screenHeight - 50}
				g.Player.Score = 0
				g.LapCount = 0
				break
			}
		}

		if rl.CheckCollisionRecs(
			rl.NewRectangle(g.Player.Position.X, g.Player.Position.Y, g.Player.Size.X, g.Player.Size.Y),
			g.Track.FinishLine) {
			g.Objects++
			g.Track.Obstacles = generateRandomObstacles(g.Objects)
			g.LapCount++
			g.Player.Score += 100
			g.Player.Position = rl.Vector2{X: screenWidth / 2, Y: screenHeight - 50}
		}

	} else {
		if rl.IsKeyPressed(rl.KeyEnter) {
			g.init()
			g.Over = false
		}
	}
}

func (g Game) draw() {
	if !g.Over {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		drawTrack(g.Track)
		drawPlayer(g.Player)

		rl.DrawText("Score: "+fmt.Sprintf("%d", g.Player.Score), 10, 10, 20, rl.Black)
		rl.DrawText("Lives: "+fmt.Sprintf("%d", g.Player.Lives), 10, 40, 20, rl.Black)
		rl.DrawText("Laps: "+fmt.Sprintf("%d", g.LapCount), 10, 70, 20, rl.Black)

		if g.Paused {
			rl.DrawText("Paused", screenWidth/2-rl.MeasureText("Paused", 40)/2, screenHeight/2-20, 40, rl.Red)
		}

		rl.EndDrawing()
	} else {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.DrawText("Game Over", screenWidth/2-rl.MeasureText("Game Over", 40)/2, screenHeight/2-20, 40, rl.Red)
		rl.DrawText("Press Enter to restart", screenWidth/2-rl.MeasureText("Press Enter to restart", 20)/2, screenHeight/2+20, 20, rl.Black)
		rl.EndDrawing()
	}
}

func drawPlayer(player Player) {
	rl.DrawCircleV(player.Position, player.Size.X/2, rl.Blue)
}

func drawTrack(track Track) {
	rl.DrawRectangleV(rl.Vector2{X: 0, Y: 0}, track.Size, track.Color)
	for _, obstacle := range track.Obstacles {
		rl.DrawRectangleV(obstacle.Position, obstacle.Size, obstacle.Color)
	}
	rl.DrawRectangleRec(track.FinishLine, rl.Green)
}

func checkCollision(player Player, obstacle Obstacle) bool {
	playerRect := rl.NewRectangle(player.Position.X, player.Position.Y, player.Size.X, player.Size.Y)
	obstacleRect := rl.NewRectangle(obstacle.Position.X, obstacle.Position.Y, obstacle.Size.X, obstacle.Size.Y)

	return rl.CheckCollisionRecs(playerRect, obstacleRect)
}

func generateRandomObstacles(count int) []Obstacle {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	obstacles := make([]Obstacle, count)

	for i := range count {
		obstacles[i] = Obstacle{
			Position: rl.Vector2{
				X: float32(rand.Intn(screenWidth - 100)),
				Y: float32(rand.Intn(screenHeight-200) + 100),
			},
			Size:  rl.Vector2{X: 50, Y: 50},
			Color: rl.Red,
		}
	}
	return obstacles
}
