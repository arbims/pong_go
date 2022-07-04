package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

const WINDOW_WIDTH = 800
const WINDOW_HEIGHT = 600
const RECT_SIZE = 25
const RECT_SPEED = 200
const FPS = 60
const DELTA_TIME_SEC = 2.0 / FPS
const BAR_LEN = 100
const TARGET_WIDTH = BAR_LEN
const TARGET_PADD = 20
const BAR_THIKNESS = 20
const BAR_Y = WINDOW_HEIGHT - BAR_THIKNESS - 100
const BAR_SPEED = RECT_SPEED

func create_rect(x float32, y float32, width int32, height int32) sdl.Rect {
	return sdl.Rect{int32(x), int32(y), width, height}
}

type Target struct {
	x    int32
	y    int32
	dead bool
}

type RectTarget struct {
	target *Target
	rect   sdl.Rect
}

func create_target_rect() [5][5]Target {

	target_pool := [5][5]Target{}
	for j := 1; j < 5; j++ {
		for i := 0; i < 5; i++ {
			target_pool[i][j] = Target{int32(100 + (TARGET_WIDTH+TARGET_PADD)*i), int32(50 * j), false}
		}
	}

	return target_pool
}

func draw_target_rect(target_pool *[5][5]Target, renderer *sdl.Renderer) []RectTarget {
	targets_rect := []RectTarget{}
	for j := 1; j < 5; j++ {
		for i := 0; i < 5; i++ {
			if target_pool[i][j].dead == false {
				target_rect := create_rect(float32(target_pool[i][j].x), float32(target_pool[i][j].y), TARGET_WIDTH, BAR_THIKNESS)
				renderer.SetDrawColor(0x00, 0xFF, 0x00, 0xFF)
				renderer.FillRect(&target_rect)
				target_with_rect := RectTarget{&target_pool[i][j], target_rect}
				targets_rect = append(targets_rect, target_with_rect)
			}
		}
	}
	return targets_rect
}

func main() {
	target_pool := create_target_rect()

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Pong game", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WINDOW_WIDTH, WINDOW_HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
	}

	var rect_dx float32 = 1
	var rect_dy float32 = 1
	var bar_y float32 = BAR_Y - BAR_THIKNESS
	var bar_x float32 = 0
	var rect_x float32 = bar_x + RECT_SIZE
	var rect_y float32 = BAR_Y - BAR_THIKNESS/2 - RECT_SIZE

	// var bar_dx float32 = 0
	Keyboard := sdl.GetKeyboardState()

	running := true

	for running {

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.KeyboardEvent:
				keyCode := t.Keysym.Sym
				switch keyCode {
				case 'q':
					running = false
				}
			}
		}
		if Keyboard[sdl.SCANCODE_RIGHT] != 0 && bar_x < WINDOW_WIDTH-BAR_LEN {
			bar_x += 10
		}

		if Keyboard[sdl.SCANCODE_LEFT] != 0 && bar_x > 0 {
			bar_x -= 10
		}

		renderer.SetDrawColor(0x18, 0x18, 0x18, 0xFF)
		renderer.Clear()

		targets_rect := draw_target_rect(&target_pool, renderer)

		var rect = create_rect(rect_x, rect_y, RECT_SIZE, RECT_SIZE)
		renderer.SetDrawColor(0xFF, 0xFF, 0xFF, 0xFF)
		renderer.FillRect(&rect)

		var bar_rect = create_rect(bar_x, bar_y, BAR_LEN, BAR_THIKNESS)
		renderer.SetDrawColor(0xFF, 0x00, 0x00, 0x00)
		renderer.FillRect(&bar_rect)

		rect_nx := rect_x + rect_dx*RECT_SPEED*DELTA_TIME_SEC

		rect = create_rect(rect_nx, rect_y, RECT_SIZE, RECT_SIZE)

		for i := 0; i < len(targets_rect); i++ {
			if targets_rect[i].rect.HasIntersection(&rect) {
				targets_rect[i].target.dead = true
				rect_dx *= -1
				rect_nx = rect_x + rect_dx*RECT_SPEED*DELTA_TIME_SEC
			}
		}

		if rect_nx < 0 || rect_nx+RECT_SIZE > WINDOW_WIDTH || bar_rect.HasIntersection(&rect) {
			rect_dx *= -1
			rect_nx = rect_x + rect_dx*RECT_SPEED*DELTA_TIME_SEC
		}
		rect_x = rect_nx

		rect_ny := rect_y + rect_dy*RECT_SPEED*DELTA_TIME_SEC
		rect = create_rect(rect_x, rect_ny, RECT_SIZE, RECT_SIZE)
		for i := 0; i < len(targets_rect); i++ {
			if targets_rect[i].rect.HasIntersection(&rect) {
				targets_rect[i].target.dead = true

				rect_dy *= -1
				rect_ny = rect_y + rect_dy*RECT_SPEED*DELTA_TIME_SEC
			}
		}
		if rect_ny < 0 || rect_ny+RECT_SIZE > WINDOW_HEIGHT || bar_rect.HasIntersection(&rect) {
			rect_dy *= -1
			rect_ny = rect_y + rect_dy*RECT_SPEED*DELTA_TIME_SEC
		}
		rect_y = rect_ny
		renderer.Present()
		sdl.Delay(1000 / FPS)
	}
	renderer.Destroy()
	window.Destroy()
}
