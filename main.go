package main

import (
	"GameBase2/app"
	"GameBase2/config"
	R "RikiKunio/res"
	"embed"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	//go:embed res
	files         embed.FS
	PlayerManager = NewPlayerManager()
	ItemManager   = NewItemManager()
)

func main() {
	config.ViewSize = complex(624, 352)
	config.Debug = true
	config.ShowFps = true
	config.Files = &files // 先使用内部资源 ，不存在  再寻找外部资源文件
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMinimum)
	app.Run(NewMainApp(), 1248, 704)
}

type MainApp struct {
	*app.App
}

// Init 必须先传入实例  初始化使用该方法
func NewMainApp() *MainApp {
	res := &MainApp{}
	res.App = app.NewApp()
	temp := config.RoomFactory.LoadAndCreate(R.MAP.LEVEL1)
	temp.AddManager(PlayerManager)
	temp.AddManager(ItemManager)
	res.PushRoom(temp)
	return res
}
