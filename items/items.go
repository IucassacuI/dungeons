package items

import (
	"deepdungeons/hero"
	"deepdungeons/libtxt"
	"time"
)

func health(x, y int) {

	var h libtxt.Object = libtxt.Object{
		Char:   '♥',
		Width:  1,
		Height: 1,
		Static: true,
		X:      x,
		Y:      y,
	}

	h.Draw()

	hero.Health += 10
	hero.Damage(0)

	time.Sleep(500 * time.Millisecond)

	h.Destroy()
	libtxt.RenderText("+10 de vida!", 13, 0, 1000)

	libtxt.Updated = true
}

func sword(x, y int) {
	var s libtxt.Object = libtxt.Object{
		Char:   '!',
		Width:  1,
		Height: 1,
		Static: true,
		X:      x,
		Y:      y,
	}

	s.Draw()

	time.Sleep(500 * time.Millisecond)

	hero.HasSword = true

	s.Destroy()

	libtxt.RenderText("Pegou a espada!", 13, 0, 2000)
}

func bow(x, y int) {
	var b libtxt.Object = libtxt.Object{
		Char:   '>',
		Width:  1,
		Height: 1,
		Static: true,
		X:      x,
		Y:      y,
	}

	b.Draw()
	hero.HasBow = true

	time.Sleep(500 * time.Millisecond)

	b.Destroy()
	libtxt.RenderText("Pegou o arco!", 13, 0, 2000)
}

func leather(x, y int) {
	var ar libtxt.Object = libtxt.Object{
		Char:   'a',
		Width:  1,
		Height: 1,
		Static: true,
		X:      x,
		Y:      y,
	}

	ar.Draw()

	hero.Armor = 50

	time.Sleep(500 * time.Millisecond)

	ar.Destroy()

	hero.Damage(0)
	libtxt.RenderText("Pegou a armadura de couro!", 13, 0, 2000)
}

func steel(x, y int) {
	var ar libtxt.Object = libtxt.Object{
		Char:   'A',
		Width:  1,
		Height: 1,
		Static: true,
		X:      x,
		Y:      y,
	}

	ar.Draw()

	hero.Armor = 150

	time.Sleep(500 * time.Millisecond)

	ar.Destroy()

	hero.Damage(0)
	libtxt.RenderText("Pegou a armadura de aco!", 13, 0, 2000)
}

func arrow(x, y int) {
	var arrows libtxt.Object = libtxt.Object{
		Char:   '→',
		Width:  1,
		Height: 1,
		Static: true,
		X:      x,
		Y:      y,
	}

	arrows.Draw()

	if hero.Arrows+10 < 51 {

		hero.Arrows += 10

		time.Sleep(500 * time.Millisecond)

		arrows.Destroy()

		hero.Damage(0)
		libtxt.RenderText("Pegou 10 flechas!", 13, 0, 2000)

	} else {
		hero.Arrows = 50
		libtxt.RenderText("Limite de flechas atingido!", 13, 0, 3000)
		arrows.Destroy()
	}
}

func Chest(x, y int, item string) {
	var c libtxt.Object = libtxt.Object{
		Char:   'B',
		Width:  1,
		Height: 1,
		Static: true,
		X:      x,
		Y:      y,
	}

	c.Draw()

	go func() {
		var opened bool
		for !opened {
			col := c.Coll("left")
			if col == "@" {
				opened = true
			}

			if libtxt.Screen.S[c.Y][c.X+1] == 0 {
				return
			}

			time.Sleep(50 * time.Millisecond)
		}
		c.Destroy()

		switch item {
		case "sword":
			sword(c.X, c.Y)
		case "health":
			health(c.X, c.Y)
		case "bow":
			bow(c.X, c.Y)
		case "leather":
			leather(c.X, c.Y)
		case "steel":
			steel(c.X, c.Y)
		case "arrow":
			arrow(c.X, c.Y)
		}

	}()
}
