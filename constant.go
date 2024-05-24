/*
@author: sk
@date: 2023/4/9
*/
package main

import "GameBase2/utils"

const (
	YScale = 0.5
	ZScale = 1.0
)

var (
	ColorLand = utils.RGBA(255, 255, 0, 128)
	ColorSide = utils.RGBA(0, 255, 255, 128)
)

const (
	TagWall = 1 << iota
)
