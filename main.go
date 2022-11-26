package main

import (
	"deepdungeons/enemies"
	"deepdungeons/hero"
	"deepdungeons/items"
	"deepdungeons/libtxt"
	_ "embed"
	"fmt"
	"strings"
	"time"
)

//go:embed title.dd
var title string

//go:embed intro/intro.dd
var intro string

//go:embed intro/cave.dd
var cave string

//go:embed intro/introtext.dd
var introtext string

//go:embed intro/introtext2.dd
var introtext2 string

//go:embed intro/gend.dd
var goodending string

//go:embed intro/bend.dd
var badending string

//go:embed intro/bendtext.dd
var bendtext string

//go:embed intro/gendtext.dd
var gendtext string

var textfiles = map[string]string{
	"good":   gendtext,
	"bad":    bendtext,
	"intro":  introtext,
	"intro2": introtext2,
}

func loadlevel(leveln int) {

	enemies.KillAll = true
	time.Sleep(3 * time.Second)
	enemies.KillAll = false

	if leveln < 10 || leveln > 14 {

		libtxt.LoadScene(fmt.Sprintf("lvls/lvl%d.dd", leveln))

	} else {
		libtxt.LoadScene("lvls/lvl1.dd")
		libtxt.Screen.S[6][48] = 'D'
		libtxt.Screen.S[4][48] = 'D'
		libtxt.Screen.S[10][48] = 'D'
		libtxt.Screen.S[12][48] = 'D'
	}

	if leveln != 20 {
		enemies.LoadEnemies(fmt.Sprintf("enms/enemies%d.dd", leveln))
	}

	switch leveln {
	case 2:
		items.Chest(22, 7, "sword")
	case 7:
		items.Chest(35, 8, "leather")
	case 8:
		items.Chest(24, 7, "health")
	case 15:
		items.Chest(46, 5, "health")
		items.Chest(46, 4, "health")
	case 16:
		items.Chest(42, 6, "health")
		items.Chest(42, 7, "bow")
		items.Chest(42, 8, "steel")
		items.Chest(42, 9, "health")
		items.Chest(42, 10, "arrow")
	case 18:
		items.Chest(21, 7, "health")
	case 19:
		items.Chest(24, 8, "health")
	case 20:
		libtxt.LoadScene("lvls/lvl20.dd")
		go enemies.SpawnMage(40, 7)
	}

	if leveln != 20 {
		libtxt.Screen.S[8][48] = 'D'
	}

	hero.SetPlayerPos(5, 8)
	hero.DrawPlayer()
	hero.Damage(0)
	libtxt.Update()

	if !hero.Movement(leveln) {
		loadlevel(leveln)
	}
}

func main() {

	var skip bool
	var timeover bool

	libtxt.Init()

	go func() {
		for !timeover {
			if libtxt.Getkeystroke() != 0 {
				skip = true
			}
		}
	}()

	time.Sleep(time.Second)
	
	timeover = true

	if !skip {
		fmt.Println(intro)

		time.Sleep(time.Second)

		text("intro")

		time.Sleep(2 * time.Second)

		fmt.Printf("\033[H\033[2J")

		time.Sleep(time.Second)

		fmt.Println(cave)

		time.Sleep(time.Second)

		text("intro2")

		time.Sleep(5 * time.Second)

		fmt.Printf("\033[H\033[2J")
	}

	fmt.Print(title)

	for libtxt.Getkeystroke() != 0 {
		time.Sleep(10 * time.Millisecond)
	}
	
	go func() {
		for {
			time.Sleep(time.Millisecond * 50)
			if libtxt.Updated {
				libtxt.Update()
				libtxt.Updated = false
			}
		}
	}()

	hero.HasSword = true
	
	for i := 19; i <= 20; i++ {
		loadlevel(i)
	}

	time.Sleep(time.Second)
	fmt.Printf("\033[H\033[2J")

	if hero.Lost {
		lose()
	} else {
		win()
	}
}

func lose() {
	text("bad")

	time.Sleep(time.Second)
	fmt.Printf("\033[H\033[2J")

	time.Sleep(2 * time.Second)

	fmt.Println("\tO Cemitério")
	time.Sleep(3 * time.Second)
	fmt.Printf("\033[H\033[2J")

	fmt.Print(badending)

	time.Sleep(30 * time.Second)
}

func win() {

	fmt.Printf("\033[H\033[2J")

	time.Sleep(time.Second)

	fmt.Println(goodending)
	text("good")

	time.Sleep(30 * time.Second)
}

func text(file string) {

	for _, line := range strings.Split(textfiles[file], "\n") {

		for _, char := range strings.Split(line, "") {
			time.Sleep(100 * time.Millisecond)
			fmt.Print(char)
		}
		fmt.Println()
		time.Sleep(2 * time.Second)
	}
}
