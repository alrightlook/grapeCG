// CrossGeta主执行类
// version 1.0 beta
// by koangel
// email: jackliu100@gmail.com
// 2017/7/26
package main

import (
	"github.com/fatih/color"
	logger "github.com/koangel/grapeNet/Logger"
)

func main() {
	color.Blue(`
	_______ _______ _______ ___ ___ 
	|     __|   |   |     __|   |   |
	|    |  |       |__     |   |   |
	|_______|__|_|__|_______|\_____/    
	`)

	color.Green("==================================================")
	color.Green("	grapeCG Gmsv 1.0 beta")
	color.Green("	CrossGate Emulator for Golang")
	color.Green("	Author:Koangel")
	color.Green("	github - github.com/koangel/grapeCG")
	color.Green("==================================================")

	RunDir := logger.GetCurrentDirectory()
	logger.BuildLogger(RunDir+"/logs", "main.log")

	logger.INFO("============ CG System Init ============")
	logger.INFO("start load conf,wait...")

}
