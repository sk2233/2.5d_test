/*
@author: sk
@date: 2023/4/9
*/
package main

import (
	"GameBase2/config"
	"GameBase2/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Input struct {
	control    bool
	XDir, YDir float64
	KeyJDown   bool
	KeyKDown   bool
	KeyUDown   bool
	KeyIDown   bool
}

func (i *Input) UpdateInput() {
	if !i.control {
		return
	}
	i.XDir = utils.GetAxis(ebiten.KeyA, ebiten.KeyD)
	i.YDir = utils.GetAxis(ebiten.KeyS, ebiten.KeyW)
	i.KeyJDown = inpututil.IsKeyJustPressed(ebiten.KeyJ)
	//i.KeyJPress = ebiten.IsKeyPressed(ebiten.KeyJ)
	i.KeyKDown = inpututil.IsKeyJustPressed(ebiten.KeyK)
	i.KeyUDown = inpututil.IsKeyJustPressed(ebiten.KeyU)
	i.KeyIDown = inpututil.IsKeyJustPressed(ebiten.KeyI)
	//i.KeyKPress = ebiten.IsKeyPressed(ebiten.KeyK)
}

//func (i *Input) JKPress() bool {
//	return (i.KeyJPress && i.KeyKDown) || (i.KeyJDown && i.KeyKPress)
//}

func NewInput(player *Player) *Input {
	if player.Control {
		config.Camera.SetTarget(player)
	}
	return &Input{control: player.Control}
}
