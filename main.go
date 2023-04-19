package main

import (
	"deepdungeons/enemies"
	"deepdungeons/hero"
	"deepdungeons/items"
	"deepdungeons/libtxt"
	"embed"
	"fmt"
	"math/rand"
	"os"
	"time"
)

//go:embed text
var text embed.FS

func savegame(leveln int) {

	var savedata []byte

	savedata = append(savedata, byte(leveln))
	savedata = append(savedata, byte(hero.Health))
	savedata = append(savedata, byte(hero.Armor))
	savedata = append(savedata, byte(hero.Arrows))

	if hero.HasSword {
		savedata = append(savedata, byte(1))
	} else {
		savedata = append(savedata, byte(0))
	}

	if hero.HasBow {
		savedata = append(savedata, byte(1))
	} else {
		savedata = append(savedata, byte(0))
	}

	err := os.WriteFile("save", savedata, 0644)

	if err != nil {
		panic(err)
	}

}

func loadlevel(leveln int) bool {

	savegame(leveln)
	
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
	case 9:
		items.Chest(32, 12, "bow")
	case 15:
		items.Chest(46, 5, "health")
		items.Chest(46, 4, "health")
	case 16:
		items.Chest(42, 6, "health")
		items.Chest(42, 7, "steel")
		items.Chest(42, 9, "health")
		items.Chest(42, 10, "arrow")
	case 18:
		items.Chest(25, 7, "health")
	case 19:
		items.Chest(24, 8, "health")
	case 20:
		libtxt.LoadScene("lvls/lvl20.dd")
		go enemies.Merto(40, 8)
	}

	if leveln != 20 {
		libtxt.Screen.S[8][48] = 'D'
	}

	hero.SetPlayerPos(5, 8)
	hero.DrawPlayer()
	hero.Damage(0)
	libtxt.Update()

	return !hero.Movement(leveln)
}

func main() {
	leveln := 1

	savedata, err := os.ReadFile("save")

	if err != nil {
		intro()
	} else {
		leveln = int(savedata[0])
		hero.Health = int(savedata[1])
		hero.Armor = int(savedata[2])
		hero.Arrows = int(savedata[3])
		hero.HasSword = savedata[4] == 1
		hero.HasBow = savedata[5] == 1
	}
	
	rand.Seed(time.Now().UnixNano())

	title, err := text.ReadFile("text/title")

	if err != nil {
		panic(err)
	}

	fmt.Printf("\033[H\033[2J")
	fmt.Println(string(title))

	_ = libtxt.Getkeystroke()

	libtxt.Init()

	go func() {
		for {
			time.Sleep(time.Millisecond * 50)
			if libtxt.Updated {
				libtxt.Update()
				libtxt.Updated = false
			}
		}
	}()

	for leveln <= 20 {
		wrongdoor := loadlevel(leveln)

		if wrongdoor {
			leveln = 1 + rand.Intn(9)
		}
		leveln++
	}

	end()
}

func end() {

	fmt.Printf("\033[H\033[2J")

	libtxt.Dialog("dialog/end")

	time.Sleep(time.Second)
	fmt.Printf("\033[H\033[2J")

	time.Sleep(2 * time.Second)

	fmt.Println("\tO Cemitério")
	time.Sleep(3 * time.Second)
	fmt.Printf("\033[H\033[2J")

	end, err := text.ReadFile("text/end")

	if err != nil {
		panic(err)
	}

	fmt.Print(string(end))

	time.Sleep(30 * time.Second)
}

func intro() {
	scenes := map[string]string{
		"text/forest": "dialog/intro1",
		"text/cave":   "dialog/intro2",
	}

	for key, val := range scenes {
		fmt.Printf("\033[H\033[2J")

		file, err := text.ReadFile(key)

		if err != nil {
			panic(err)
		}

		fmt.Println(string(file))

		time.Sleep(3 * time.Second)
		libtxt.Dialog(val)
	}
}
