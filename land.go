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
	config.ObjectFactory.RegisterRectFactory(R.CLASS.LAND, createLand)
}

func createLand(data *model.ObjectData) model.IObject {
	res := &Land{}
	res.RectObject = object.NewRectObject()
	factory.FillRectObject(data, res.RectObject)
	res.Pos, res.Size = To3DPosAndSize(res.Pos, res.Size)
	res.Z = res.GetFloat(R.PROP.HEIGHT, 0)
	res.Tag = TagWall
	return res
}

type Land struct {
	*object.RectObject
	Z float64
}

func (l *Land) GetZ(x, y float64) float64 {
	return l.Z
}

func (l *Land) Order() int {
	return -int(imag(l.Pos))
}

func (l *Land) Draw(screen *ebiten.Image) {
	x, y := utils.Vector2Float(l.Pos)
	w, h := utils.Vector2Float(l.Size)
	utils.FillRect(screen, To2DPos(x, y, l.Z), complex(w, l.Z*ZScale), ColorSide)
	utils.FillRect(screen, To2DPos(x, y+h, l.Z), complex(w, h*YScale), ColorLand)
}
