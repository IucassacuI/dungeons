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

var Pdirection string

var Health int = 100
var Armor int

var HasSword bool
var HasBow bool
var Arrows int
var CanMove bool = true
var Lost bool
var magefight bool

func swordattack(creature libtxt.Object, direction string) {

	switch direction {
	case "up":
		sword.X = creature.X
		sword.Y = creature.Y - 1
		sword.Char = '|'
		sword.Width = 1
	case "down":
		sword.X = creature.X
		sword.Y = creature.Y + 1
		sword.Char = '|'
		sword.Width = 1
	case "left":
		sword.X = creature.X - 2
		sword.Y = creature.Y
	case "right":
		sword.X = creature.X + 1
		sword.Y = creature.Y
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

	var a libtxt.Object = libtxt.Object{
		Char:   c,
		Width:  1,
		Height: 1,
		X:      x,
		Y:      y,
	}

	a.Draw()

	go func() {
		var col string = a.Coll(direct)

		for col == "" {
			col = a.Coll(direct)
			a.Move(direct)
			time.Sleep(100 * time.Millisecond)
		}

		col = a.Coll(direct)

		if col != "" {

			time.Sleep(300 * time.Millisecond)
			a.Destroy()

		}

		a.Destroy()
	}()

	Arrows -= 1
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

	if level == 20 {
		magefight = true
	}

	for {

		if !CanMove {
			continue
		}

		key = libtxt.Getkeystroke()

		if Health <= 0 {
			break
		}

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
				swordattack(player, Pdirection)
			}
		case 'n':
			if HasBow && Arrows > 0 {
				bow()
			}
		case 'x':
			os.Exit(0)
		}

		if player.Coll("") != "" {
			if player.Coll("")[0] == 'D' {
				goto END
			}
		}
		time.Sleep(120 * time.Millisecond)
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

	if magefight && Health < 10 {
		Lost = true
		return
	}

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
