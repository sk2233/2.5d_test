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
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	config.ObjectFactory.RegisterPointFactory(R.CLASS.PLAYER, createPlayer)
}

func createPlayer(data *model.ObjectData) model.IObject {
	res := &Player{Z: 0, drawOption: model.NewDrawOption(), dir: 1}
	res.PointObject = object.NewPointObject()
	factory.FillPointObject(data, res.PointObject)
	res.Control = res.GetBool(R.PROP.CONTROL, false)
	res.input = NewInput(res)
	res.actions = []func(){res.down, res.kick, res.skill, res.air, utils.EmptyFunc, res.fly, res.hurt, res.land}
	attack1 := CreateDynamicSprite(R.MAIN.PLAYER.ATTACK.ATTACK1)
	attack1.AddEventFrame(1)
	attack2 := CreateDynamicSprite(R.MAIN.PLAYER.ATTACK.ATTACK2)
	attack2.AddEventFrame(1)
	attack3 := CreateDynamicSprite(R.MAIN.PLAYER.ATTACK.ATTACK3)
	attack3.AddEventFrame(1)
	knife := CreateDynamicSprite(R.MAIN.PLAYER.ATTACK.KNIFE)
	knife.AddEventFrame(2)
	stick := CreateDynamicSprite(R.MAIN.PLAYER.ATTACK.STICK)
	stick.AddEventFrame(2)
	back := CreateDynamicSprite(R.MAIN.PLAYER.ATTACK.BACK)
	back.AddEventFrame(1)
	res.sprites = [][]model.ISprite{
		{CreateStaticSprite(R.MAIN.PLAYER.AIR.DOWN), CreateStaticSprite(R.MAIN.PLAYER.HIT.LAND)},
		{CreateStaticSprite(R.MAIN.PLAYER.AIR.FIST), CreateStaticSprite(R.MAIN.PLAYER.AIR.FOOT)},
		{CreateDynamicSprite(R.MAIN.PLAYER.AIR.SKILL)},
		{CreateStaticSprite(R.MAIN.PLAYER.AIR.UP)},
		{attack1, attack2, attack3, knife, stick, back},
		{CreateStaticSprite(R.MAIN.PLAYER.HIT.FLY)},
		{CreateFrameSprite(R.MAIN.PLAYER.HIT.NORMAL)},
		{CreateFrameSprite(R.MAIN.PLAYER.STAND.NORMAL), CreateFrameSprite(R.MAIN.PLAYER.STAND.KNIFE), CreateFrameSprite(R.MAIN.PLAYER.STAND.STICK)},
	}
	res.Pos = complex(real(res.Pos), (352-imag(res.Pos))/YScale)
	res.SetState(7, res.tool)
	return res
}

type Player struct {
	*object.PointObject
	Control bool
	Z       float64
	input   *Input
	actions []func()
	state   int // 0 落地  1 飞踢 2 技能  3 上跳/下落  4 攻击  5 击飞  6 挨打  7 走路
	type0   int // 不同情况下 不一样
	// 0 : 0 普通落地   1 击飞落地(不能重拳)
	// 1 : 0 普通飞踢   1 超级飞踢(踢的更远)
	// 4 : 0 拳1  1 拳2  2 上勾拳  3 飞刀  4 棍    5 脚
	// 7 : 0 普通走路   1  飞刀   2  棍
	sprites     [][]model.ISprite // 有的状态下 有3种  捡道具
	tool        int               //  0 无道具  1 飞刀   2  棍
	level       int               // 攻击等级
	sprite      model.ISprite
	drawOption  *model.DrawOption
	zSpeed      float64
	dir         float64
	startFrame  int
	attackFrame int
}

func (p *Player) GetMin() complex128 {
	return complex(real(p.Pos)-6, p.Z)
}

func (p *Player) GetMax() complex128 {
	return complex(real(p.Pos)+6, p.Z+40)
}

func (p *Player) AnimEvent(frame int) {
	if p.state == 2 {
		return
	} // 4
	switch p.type0 {
	case 0, 1: // 绝对值 3
		PlayerManager.Attack(real(p.Pos), real(p.Pos)+18*p.dir, imag(p.Pos), p.Z+28, false, p, 3)
	case 2:
		PlayerManager.Attack(real(p.Pos), real(p.Pos)+18*p.dir, imag(p.Pos), p.Z+28, true, p, 3)
	case 5:
		PlayerManager.Attack(real(p.Pos), real(p.Pos)-20*p.dir, imag(p.Pos), p.Z+20, false, p, 3)
	case 3:
		p.tool = 0
		NewKnife(p.Pos, p.Z+30, p.dir)
	case 4:
		PlayerManager.Attack(real(p.Pos), real(p.Pos)+16*p.dir, imag(p.Pos), p.Z+18, true, p, 16)
	}
}

func (p *Player) AnimEnd() {
	switch p.state {
	case 4:
		p.SetState(7, p.tool)
	case 2:
		p.SetState(3, 0)
	}
}

func (p *Player) Order() int {
	return -int(imag(p.Pos))
}

func (p *Player) Draw(screen *ebiten.Image) {
	p.drawOption.Image.GeoM.Reset()
	p.drawOption.Image.GeoM.Translate(p.sprite.GetOffset())
	p.drawOption.Image.GeoM.Scale(p.dir, 1)
	x, y := utils.Vector2Float(p.Pos)
	p.drawOption.Image.GeoM.Translate(x, 352-y*YScale-p.Z*ZScale)
	p.sprite.Draw(screen, p.drawOption)
}

func (p *Player) Update() {
	utils.InvokeUpdate(p.sprite)
	p.input.UpdateInput()
	p.actions[p.state]()
	if p.input.XDir != 0 {
		p.dir = p.input.XDir
	}
}

func (p *Player) SetState(state int, type0 int) {
	p.state = state
	p.type0 = type0
	p.sprite = p.sprites[state][type0]
	p.sprite.SetTarget(p)
	utils.InvokeAnimReset(p.sprite)
}

//===================action===================

func (p *Player) down() {
	if config.Frame-p.startFrame > 15 {
		p.SetState(7, p.tool)
	}
	if p.type0 == 0 {
		p.checkLandAttack(true)
	}
}

func (p *Player) land() {
	p.checkMove(true)
	p.checkLand()
	p.checkJump()
	p.checkLandAttack(false)
	if p.input.KeyIDown {
		tool := ItemManager.CheckItem(p.Pos)
		if tool > 0 {
			p.tool = tool
			p.startFrame = config.Frame
			p.SetState(0, 0)
		}
	}
}

func (p *Player) kick() {
	p.checkDown(false)
	x, y := utils.Vector2Float(p.Pos)
	if GetZ(x+p.dir*6, y) > p.Z {
		x = float64(int(x - p.dir))
		for i := 0; i < 3; i++ {
			if GetZ(x+p.dir, y) > p.Z {
				break
			}
			x += p.dir
		}
	} else {
		x += p.dir * 6
	}
	p.Pos = complex(x, y)
	if config.Frame%5 == 0 {
		if p.type0 == 0 {
			PlayerManager.Attack(real(p.Pos), real(p.Pos)+22*p.dir, imag(p.Pos), p.Z+8, false, p, 3)
		} else {
			PlayerManager.Attack(real(p.Pos), real(p.Pos)+12*p.dir, imag(p.Pos), p.Z+20, true, p, 3)
		}
	}
}

func (p *Player) skill() {
	if config.Frame%5 == 0 {
		PlayerManager.Attack(real(p.Pos)-18, real(p.Pos)+18, imag(p.Pos), p.Z+18, true, p, 18)
	}
}

func (p *Player) air() {
	p.checkMove(false)
	p.checkDown(false)
	if p.input.KeyJDown {
		p.SetState(1, 0)
	}
	if p.input.KeyKDown {
		p.SetState(1, 1)
	}
	if p.input.KeyIDown {
		p.SetState(2, 0)
	}
}

func (p *Player) fly() {
	x, y := utils.Vector2Float(p.Pos)
	if GetZ(x-p.dir*3, y) > p.Z {
		x = float64(int(x + p.dir))
		for i := 0; i < 3; i++ {
			if GetZ(x-p.dir, y) > p.Z {
				break
			}
			x -= p.dir
		}
	} else {
		x -= p.dir * 3
	}
	p.Pos = complex(x, y)
	p.checkDown(true)
}

func (p *Player) hurt() {
	if config.Frame-p.startFrame > 15 {
		p.SetState(7, p.tool)
	}
}

//====================check=====================

func (p *Player) checkMove(land bool) { // 检查移动  空中 陆地 不太相同
	x, y := utils.Vector2Float(p.Pos)
	xDir := p.input.XDir
	if xDir != 0 {
		offset := 0.0
		if land {
			offset = 1
		}
		if GetZ(x+xDir*2, y) > p.Z+offset {
			x = float64(int(x - xDir)) // 防止 误差 进入上层
			for i := 0; i < 2; i++ {
				if GetZ(x+xDir, y) > p.Z+offset {
					break
				}
				x += xDir
				p.Z += offset
			}
		} else {
			x += xDir * 2
			p.Z += offset * 2
		}
		z := GetZ(x, y)
		if p.Z-z < offset*2*2 {
			p.Z = z
		}
	}
	yDir := p.input.YDir
	if yDir != 0 {
		if GetZ(x, y+yDir*2) > p.Z {
			y = float64(int(y - yDir))
			for i := 0; i < 2; i++ {
				if GetZ(x, y+yDir) > p.Z {
					break
				}
				y += yDir
			}
		} else {
			y += yDir * 2
		}
	}
	if xDir == 0 && yDir == 0 {
		return
	}
	p.Pos = complex(x, y)
	if config.Frame%6 == 0 {
		utils.InvokeNext(p.sprite)
	}
}

func (p *Player) checkLand() { // 检查在地面上
	if GetZ(utils.Vector2Float(p.Pos)) < p.Z {
		p.zSpeed = 0
		p.SetState(3, 0)
	}
}

func (p *Player) checkDown(hurt bool) { // 检测落地  击飞与普通落地不同
	p.zSpeed -= 0.3
	if p.zSpeed < 0 { // 只管下落
		z := GetZ(utils.Vector2Float(p.Pos))
		if p.Z+p.zSpeed < z {
			p.Z = z
			p.startFrame = config.Frame
			if hurt {
				p.SetState(0, 1)
			} else {
				p.SetState(0, 0)
			}
		} else {
			p.Z += p.zSpeed
		}
	} else {
		p.Z += p.zSpeed
	}
}

func (p *Player) checkJump() {
	if p.input.KeyUDown {
		p.zSpeed = 6
		p.SetState(3, 0)
	}
}

func (p *Player) checkLandAttack(down bool) { // 攻击检测  落地时(非击落)攻击 必定上勾拳
	if p.input.KeyJDown { // 多种拳  脚不是这个按键
		if p.tool > 0 { // 优先使用工具
			p.SetState(4, p.tool+2)
			return
		}
		if down { //落地时(非击落)攻击 必定上勾拳
			p.SetState(4, 2)
			return
		}
		if config.Frame-p.attackFrame < 15 {
			p.level = (p.level + 1) % 3
		} else {
			p.level = 0
		}
		p.attackFrame = config.Frame
		p.SetState(4, p.level)
		return
	}
	if p.input.KeyKDown { // 脚
		p.SetState(4, 5)
	}
}

func (p *Player) Collision(x1 float64, x2 float64, y float64, z float64, offset float64) bool {
	if p.state == 2 { // 技能期间无敌
		return false
	}
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	return utils.HorizontalCollision(p, x1, x2, z) && math.Abs(imag(p.Pos)-y) < offset
}

func (p *Player) Hurt(fly bool, src model.IPos) {
	if p.state == 1 || p.state == 3 || p.state == 5 { // 击飞处理
		p.SetState(5, 0)
		p.dir = utils.Sign(real(src.GetPos()) - real(p.Pos))
	} else if fly { // 地面上击飞处理
		p.zSpeed = 6
		p.SetState(5, 0)
		p.dir = utils.Sign(real(src.GetPos()) - real(p.Pos))
	} else { // 普通挨打
		p.SetState(6, 0)
		utils.InvokeNext(p.sprite)
		p.startFrame = config.Frame
	}
}
