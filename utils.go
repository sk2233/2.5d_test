/*
@author: sk
@date: 2023/4/9
*/
package main

import (
	"GameBase2/config"
	"GameBase2/model"
	"GameBase2/utils"
	R "RikiKunio/res"

	"github.com/hajimehoshi/ebiten/v2"
)

func To3DPosAndSize(pos complex128, size complex128) (complex128, complex128) {
	x, y := utils.Vector2Float(pos)
	w, h := utils.Vector2Float(size)
	return complex(x, (352-y-h)/YScale), complex(w, h/YScale)
}

func To2DPos(x, y, z float64) complex128 {
	return complex(x, 352-y*YScale-z*ZScale)
}

func CreateStaticSprite(name string) model.ISprite {
	return config.SpriteFactory.CreateStaticSprite(R.SPRITE.MAIN, name)
}

func CreateDynamicSprite(name string) model.IDynamicSprite {
	return config.SpriteFactory.CreateDynamicSprite(R.SPRITE.MAIN, name)
}

func CreateFrameSprite(name string) model.IFrameSprite {
	return config.SpriteFactory.CreateFrameSprite(R.SPRITE.MAIN, name)
}

func GetZ(x, y float64) float64 {
	res := utils.CollisionPoint(R.LAYER.COLLISION, TagWall, complex(x, y))
	return InvokeGetZ(res, x, y)
}

func InvokeGetZ(src any, x, y float64) float64 {
	if tar, ok := src.(IGetZ); ok {
		return tar.GetZ(x, y)
	}
	return 0
}

func GetImage(name string) *ebiten.Image {
	return config.SpritesLoader.LoadStaticSprite(R.SPRITE.MAIN, name).Image
}
