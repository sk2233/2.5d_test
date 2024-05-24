/*
@author: sk
@date: 2023/4/11
*/
package main

import (
	"GameBase2/object"
	"GameBase2/utils"
	R "RikiKunio/res"

	"github.com/hajimehoshi/ebiten/v2"
)

type Knife struct {
	*object.PointObject
	z   float64
	img *ebiten.Image
	dir float64
	die bool
}

func (k *Knife) Draw(screen *ebiten.Image) {
	x, y := utils.Vector2Float(k.Pos)
	pos := To2DPos(x, y, k.z)
	utils.DrawScaleImage(screen, k.img, pos-4-3i, -4-3i, complex(k.dir, 1))
}

func (k *Knife) IsDie() bool {
	return k.die
}

func (k *Knife) Update() {
	x, y := utils.Vector2Float(k.Pos)
	x += k.dir * 6
	if x < 0 || x > 992 || GetZ(x, y) > k.z {
		k.die = true
	}
	k.Pos = complex(x, y)
	if PlayerManager.Attack(x-4, x+4, y, k.z, true, k, 4) {
		k.die = true
	}
}

func NewKnife(pos complex128, z float64, dir float64) *Knife {
	res := &Knife{z: z, img: GetImage(R.MAIN.ITEM.KNIFE), dir: dir, die: false}
	res.PointObject = object.NewPointObject()
	res.Pos = pos
	utils.AddToLayer(R.LAYER.PLAYER, res)
	return res
}
