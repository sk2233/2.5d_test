/*
@author: sk
@date: 2023/4/10
*/
package main

import (
	"GameBase2/object"
	"GameBase2/utils"
	R "RikiKunio/res"

	"github.com/hajimehoshi/ebiten/v2"
)

type Item struct {
	*object.PointObject
	img     *ebiten.Image
	type0   int // 1 飞刀  2 棍
	z       float64
	targetZ float64
	zSpeed  float64
	timer   int
}

func (i *Item) Update() {
	i.timer--
	if i.z <= i.targetZ {
		return
	}
	i.zSpeed -= 0.3
	i.z += i.zSpeed
	if i.z < i.targetZ {
		i.z = i.targetZ
	}
}

func (i *Item) Order() int {
	return -int(imag(i.Pos))
}

func (i *Item) IsDie() bool {
	return i.timer < 0
}

func NewItem(pos complex128, type0 int) *Item {
	res := &Item{type0: type0, z: 256, timer: 900}
	res.PointObject = object.NewPointObject()
	res.Pos = pos
	res.targetZ = GetZ(utils.Vector2Float(pos))
	if type0 == 1 {
		res.img = GetImage(R.MAIN.ITEM.KNIFE)
	} else {
		res.img = GetImage(R.MAIN.ITEM.STICK)
	}
	utils.AddToLayer(R.LAYER.PLAYER, res)
	return res
}

var (
	anchors = []complex128{0, 4 + 3i, 8 + 2i}
)

func (i *Item) Draw(screen *ebiten.Image) {
	x, y := utils.Vector2Float(i.Pos)
	pos := To2DPos(x, y, i.z)
	utils.DrawImage(screen, i.img, pos-anchors[i.type0])
}
