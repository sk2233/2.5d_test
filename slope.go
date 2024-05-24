/*
@author: sk
@date: 2023/4/9
*/
package main

import (
	"GameBase2/config"
	"GameBase2/factory"
	"GameBase2/model"
	"GameBase2/object"
	"GameBase2/utils"
	R "RikiKunio/res"

	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	config.ObjectFactory.RegisterRectFactory(R.CLASS.SLOPE, createSlope)
}

func createSlope(data *model.ObjectData) model.IObject {
	res := &Slope{}
	res.RectObject = object.NewRectObject()
	factory.FillRectObject(data, res.RectObject)
	res.Pos, res.Size = To3DPosAndSize(res.Pos, res.Size)
	res.Left = res.GetFloat(R.PROP.LEFT, 0)
	res.Right = res.GetFloat(R.PROP.RIGHT, 0)
	res.Tag = TagWall
	return res
}

type Slope struct {
	*object.RectObject
	Left, Right float64
}

func (s *Slope) GetZ(x, y float64) float64 {
	return s.Left + (x-real(s.Pos))/real(s.Size)*(s.Right-s.Left)
}

func (s *Slope) Draw(screen *ebiten.Image) {
	x, y := utils.Vector2Float(s.Pos)
	w, h := utils.Vector2Float(s.Size)
	vs := make([]ebiten.Vertex, 3)
	vs[0] = utils.NewVertex(x, 352-y*YScale)
	vs[1] = utils.NewVertex(x+w, 352-y*YScale)
	if s.Left > 0 {
		vs[2] = utils.NewVertex(x, 352-y*YScale-s.Left*ZScale)
	} else {
		vs[2] = utils.NewVertex(x+w, 352-y*YScale-s.Right*ZScale)
	}
	utils.FillTriangle(screen, vs, ColorSide)
	vs = make([]ebiten.Vertex, 4)
	vs[0] = utils.NewVertex(x, 352-(y+h)*YScale-s.Left*ZScale)
	vs[1] = utils.NewVertex(x+w, 352-(y+h)*YScale-s.Right*ZScale)
	vs[2] = utils.NewVertex(x+w, 352-y*YScale-s.Right*ZScale)
	vs[3] = utils.NewVertex(x, 352-y*YScale-s.Left*ZScale)
	utils.FillFan(screen, vs, ColorLand)
}

func (s *Slope) Order() int {
	return -int(imag(s.Pos))
}
