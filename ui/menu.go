package ui

import (
	"fmt"
	"os"
)

func RunMenu() {
	var op int
	for {
		mainMenu(&op)
		switch op {
		case 0:
			os.Exit(0)
		}

	}
}

func mainMenu(op *int) {
	fmt.Println("============= 操作面板 =============")
	fmt.Println("       0.退出")
	fmt.Print("请输入：")
	fmt.Scanln(op)
}
