package ui

import (
	"fmt"
	"os"
	"ribin-game-robot/robot"
)

func RunMenu() {
	var op int
	for {
		mainMenu(&op)
		switch op {
		case 0:
			test()
		case -1:
			os.Exit(0)
		}

	}
}

func mainMenu(op *int) {
	fmt.Println("============= 操作面板 =============")
	fmt.Println("       0.测试")
	fmt.Println("       -1.退出")
	fmt.Print("请输入：")
	fmt.Scanln(op)
}

func test() {
	r := robot.NewRobot("ribinaco")
	r.EnterRoom("home")
}
