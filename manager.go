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
	"reflect"
)

//==================playerManager=======================

type playerManager struct {
	players []*Player
}

func (p *playerManager) Init() {
	temps := utils.GetObjectLayer(R.LAYER.PLAYER).GetObjectsByType(reflect.TypeOf(&Player{}))
	for _, temp := range temps {
		if player := temp.(*Player); !player.Control {
			p.players = append(p.players, player)
		}
	}
}

func (p *playerManager) Attack(x1 float64, x2 float64, y float64, z float64, fly bool, src model.IPos, offset float64) bool {
	for _, player := range p.players {
		if player != src && player.Collision(x1, x2, y, z, offset) {
			player.Hurt(fly, src)
			return true
		}
	}
	return false
}

func NewPlayerManager() *playerManager {
	return &playerManager{}
}

//=================itemManager=================

type itemManager struct {
	items []*Item
}

func (i *itemManager) Update() {
	if config.Frame%300 == 0 {
		pos := complex(utils.RandomFloat(184, 992-184), utils.RandomFloat(0, 320))
		i.items = append(i.items, NewItem(pos, utils.RandomItem(1, 2)))
	}
}

func (i *itemManager) CheckItem(pos complex128) int {
	i.items = utils.FilterSlice(i.items, func(item *Item) bool {
		return !item.IsDie()
	})
	for _, item := range i.items {
		if utils.VectorLen2(item.Pos-pos) < 64 {
			item.timer = -2233
			return item.type0
		}
	}
	return -1
}

func NewItemManager() *itemManager {
	return &itemManager{items: make([]*Item, 0)}
}
