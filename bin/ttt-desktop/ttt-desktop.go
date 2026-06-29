package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:web-build
var MyAssets embed.FS

func main() {
	var e error
	e=wails.Run(&options.App{
        Title: "Task Time Tracker Desktop",
		Width:  1024,
		Height: 768,
        AssetServer: &assetserver.Options{
            Assets: MyAssets,
        },
    })

    if e!=nil {
        panic(e)
    }
}