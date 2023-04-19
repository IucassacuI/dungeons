package hero

import (
	"deepdungeons/libtxt"
	"fmt"
	"os"
	"time"
)

var player libtxt.Object = libtxt.Object{
	Char:   '@',
	Width:  1,
	Height: 1,
	X:      4,
	Y:      7,
}

var sword libtxt.Object = libtxt.Object{
	Char:   '-',
	Width:  2,
	Height: 1,
	Static: true,
	X:      0,
	Y:      0,
}

var (
	Pdirection string
	Health     int = 100
	Armor      int
	HasSword   bool
	HasBow     bool
	Arrows     int
	CanMove    bool = true
	Died       bool
)

func swordattack(direction string) {

	switch direction {
	case "up":
		sword.X = player.X
		sword.Y = player.Y - 1
		sword.Char = '|'
		sword.Width = 1
	case "down":
		sword.X = player.X
		sword.Y = player.Y + 1
		sword.Char = '|'
		sword.Width = 1
	case "left":
		sword.X = player.X - 2
		sword.Y = player.Y
	case "right":
		sword.X = player.X + 1
		sword.Y = player.Y
	}

	sword.Draw()
	time.Sleep(300 * time.Millisecond)
	sword.Destroy()

	sword.X = 0
	sword.Y = 0
	sword.Char = '-'
	sword.Width = 2
}

func bow() {

	var x, y int
	var c rune = '-'
	var direct = Pdirection

	switch direct {
	case "right":
		x = player.X + 2
		y = player.Y
	case "left":
		x = player.X - 2
		y = player.Y
	case "up":
		c = '|'
		x = player.X
		y = player.Y - 2
	case "down":
		c = '|'
		x = player.X
		y = player.Y + 2
	}

	var arrow libtxt.Object = libtxt.Object{
		Char:   c,
		Width:  1,
		Height: 1,
		X:      x,
		Y:      y,
	}

	arrow.Draw()

	go func() {

		for arrow.Coll(direct) == "" {
			arrow.Move(direct)
			time.Sleep(100 * time.Millisecond)
		}

		time.Sleep(300 * time.Millisecond)
		arrow.Destroy()
	}()

	Arrows--
	Damage(0)
}

func GetPlayerPos(axis byte) int {
	switch axis {
	case 'X':
		return player.X
	case 'Y':
		return player.Y
	}
	return 0
}

func SetPlayerPos(x, y int) {
	player.Destroy()

	player.X = x
	player.Y = y

	player.Draw()
}

func DrawPlayer() {
	player.Draw()
}

func Movement(level int) bool {
	var key rune

	for {
		time.Sleep(120 * time.Millisecond)

		if !CanMove {
			continue
		}

		if Health <= 0 {
			break
		}

		key = libtxt.Getkeystroke()

		switch key {
		case 'w':
			Pdirection = "up"
			player.Move("up")
		case 'a':
			Pdirection = "left"
			player.Move("left")
		case 's':
			Pdirection = "down"
			player.Move("down")
		case 'd':
			Pdirection = "right"
			player.Move("right")
		case 'm':
			if HasSword {
				swordattack(Pdirection)
			}
		case 'n':
			if HasBow && Arrows > 0 {
				bow()
			}
		case 'x':
			os.Exit(0)
		}

		col := player.Coll("")

		if col != "" && col[0] == 'D' {
			goto END

		}
	}

END:

	switch level {
	case 10:
		return player.Y == 10
	case 11:
		return player.Y == 6
	case 12:
		return player.Y == 12
	case 13:
		return player.Y == 4
	case 14:
		return player.Y == 10
	default:
		return true
	}

	return false
}

func Damage(amount int) {

	if Armor > 0 {
		Armor -= amount
	} else {
		Health -= amount
	}

	libtxt.RenderText(fmt.Sprintf("Saude: %d", Health), 0, 0, 0)
	libtxt.RenderText(fmt.Sprintf("Armadura: %d", Armor), 0, 1, 0)
	libtxt.RenderText(fmt.Sprintf("Flechas: %d", Arrows), 15, 1, 0)

	libtxt.Update()

	if Health <= 0 {
		player.Char = '+'
		player.Draw()
		libtxt.Update()
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}

	if Armor < 0 {
		Health += Armor
		Armor = 0
	}
}
