package main

import "fmt"

func userInput (a int){
	switch a {
	case 1 :
		skillsChoose()
	case 2:
		fmt.Println("恭喜您，逃跑成功")

	}
}

func skillsChoose(){
	var userInput2 int
	fmt.Println("请选择您想释放的技能:")
	fmt.Println("1.龙卷风摧毁停车场")
	fmt.Println("2.返回")
	fmt.Scan(&userInput2)
	switch userInput2 {
	case 1:
		ReleaseSkill("龙卷风摧毁停车场", func(skillName string) {
			fmt.Println("尝尝我的厉害吧！", skillName)
		})
	case 2:
		main()
	}
}

func main() {
	var inputInt int
	fmt.Println("请输入你想执行的操作:")
	fmt.Println("1.释放技能")
	fmt.Println("2.逃跑")
	fmt.Scan(&inputInt)
	userInput(inputInt)

}

func ReleaseSkill(skillNames string, releaseSkillFunc func(string)) {
releaseSkillFunc(skillNames)
}