package enemies

import (
	"bufio"
	"deepdungeons/hero"
	"deepdungeons/items"
	"deepdungeons/libtxt"
	"embed"
	"fmt"
	"math/rand"
	"time"
)

var KillAll bool

//go:embed enms
var enms embed.FS

func gotoplayerpos(enemy *libtxt.Object) string {

	var direction string

	if enemy.X < hero.GetPlayerPos('X') {
		enemy.Move("right")
		direction = "right"
	} else if enemy.X > hero.GetPlayerPos('X') {
		enemy.Move("left")
		direction = "left"
	}

	if enemy.Y < hero.GetPlayerPos('Y') {
		enemy.Move("down")
	} else if enemy.Y > hero.GetPlayerPos('Y') {
		enemy.Move("up")
	}

	return direction
}

func damage(enemy *libtxt.Object, health int, char rune, damage int) bool {

	var col = enemy.Coll("")

	if col == "" {
		goto END
	}

	if col[0] == '-' || col[0] == '|' {
		enemy.Char = rune(48 + health - 1)

		enemy.Draw()
		time.Sleep(200 * time.Millisecond)
		enemy.Destroy()

		enemy.Char = char
		enemy.Draw()
		return true
	}

	if col[0] == '@' {
		hero.Damage(damage)
	}

END:
	return false
}

func SpawnWalker(x, y int) {

	var walker libtxt.Object = libtxt.Object{
		Char:   '$',
		Width:  1,
		Height: 1,
		X:      x,
		Y:      y,
	}
	walker.Draw()

	var health int = 2
	var speed = time.Second

	go func() {

		for health != 0 && !KillAll {

			if damage(&walker, health, walker.Char, 1) {
				speed = 400 * time.Millisecond
				health--
			}

			time.Sleep(400 * time.Millisecond)
		}
	}()

	go func() {

		time.Sleep(time.Second)

		for health != 0 && !KillAll {

			_ = gotoplayerpos(&walker)

			time.Sleep(speed)
		}
		walker.Destroy()
	}()
}

func SpawnRunner(x, y int) {
	var runner libtxt.Object = libtxt.Object{
		Char:   '/',
		Width:  1,
		Height: 1,
		X:      x,
		Y:      y,
	}
	runner.Draw()

	go func() {
		var health int = 3
		var speed = 400 * time.Millisecond

		time.Sleep(speed)

		for health != 0 && !KillAll {

			_ = gotoplayerpos(&runner)

			if damage(&runner, health, runner.Char, 2) {
				speed = 200 * time.Millisecond
				health--
			}

			time.Sleep(speed)
		}
		runner.Destroy()
	}()
}

func projectile(x, y int, direction string, magic bool) {
	var p libtxt.Object = libtxt.Object{
		Char:   '-',
		Width:  1,
		Height: 1,
		X:      x,
		Y:      y,
	}

	var speed time.Duration = 100

	if magic {
		p.Char = '*'
		speed = 50
	}

	for p.Coll(direction) == "" && !KillAll {

		p.Move(direction)
		time.Sleep(speed * time.Millisecond)
	}

	col := p.Coll(direction)

	if p.Y == hero.GetPlayerPos('Y') && col != "" && col[0] == '@' {
		if magic {
			hero.Damage(10)
		} else {
			hero.Damage(3)
		}
		p.Destroy()
	}

	p.Destroy()
}

func SpawnArcher(x, y int) {
	var archer libtxt.Object = libtxt.Object{
		Char:   '<',
		Width:  1,
		Height: 1,
		X:      x,
		Y:      y,
	}

	archer.Draw()

	var health int = 3
	var playerkill bool = true
	var speed = 500 * time.Millisecond

	go func() {

		for health != 0 {

			if KillAll {
				playerkill = false
				health = 0
			}

			if damage(&archer, health, archer.Char, 3) {
				speed = 250 * time.Millisecond
				health--
			}

			time.Sleep(400 * time.Millisecond)
		}
	}()

	go func() {

		time.Sleep(time.Second)

		for health != 0 && !KillAll {

			hx := hero.GetPlayerPos('X')

			for hx-archer.X < 15 && hx-archer.X > 0 && !KillAll && health != 0 {
				archer.Move("left")
				time.Sleep(500 * time.Millisecond)
			}

			for hx-archer.X < -5 && hx-archer.X > -25 && !KillAll && health != 0 {
				archer.Char = '>'
				archer.Move("right")
				time.Sleep(time.Second)
			}

			if hx < archer.X {
				archer.Char = '<'
				archer.Draw()
				for i := 0; i <= 3 && !KillAll && health != 0; i++ {
					time.Sleep(time.Second)
					projectile(archer.X-2, archer.Y, "left", false)
				}
			}

			if hx > archer.X {
				archer.Char = '>'
				archer.Draw()
				for i := 0; i <= 3 && !KillAll && health != 0; i++ {
					time.Sleep(speed)
					projectile(archer.X+2, archer.Y, "right", false)
				}

			}

			if KillAll {
				playerkill = false
				health = 0
			}

			time.Sleep(time.Second)
			archer.Char = '<'
		}
		archer.Destroy()

		if playerkill {
			items.Chest(archer.X, archer.Y, "arrow")
		}
	}()
}

func SpawnBat(x, y int) {
	var bat libtxt.Object = libtxt.Object{
		Char:   '|',
		Width:  1,
		Height: 1,
		X:      x,
		Y:      y,
	}

	bat.Draw()

	var health int = 1

	go func() {

		for health != 0 && !KillAll {

			if damage(&bat, health, bat.Char, 0) {
				health--
			}

			time.Sleep(10 * time.Millisecond)
		}

		bat.Destroy()
	}()

	go func() {

		directions := []string{"up", "down", "left", "right"}

		time.Sleep(time.Second)

		direction := "left"
		rand.Seed(time.Now().UnixNano())

		hx := hero.GetPlayerPos('X')

		for hx+3 != bat.X && hx-3 != bat.X && !KillAll {
			hx = hero.GetPlayerPos('X')
			time.Sleep(50 * time.Millisecond)
		}

		bat.Char = '^'

		for health != 0 && !KillAll {
			direction = directions[rand.Intn(3)]

			for i := 0; i <= 15 && !KillAll; i++ {

				if bat.X == 1 {
					direction = "right"
				} else if bat.X == 48 {
					direction = "left"
				}

				if bat.Y == 4 {
					direction = "down"
				} else if bat.Y == 13 {
					direction = "up"
				}

				bat.Move(direction)
				time.Sleep(200 * time.Millisecond)
			}
		}
		bat.Destroy()
	}()
}

func SpawnMultiplier(x, y int, canclone bool) {

	var multiplier libtxt.Object = libtxt.Object{
		Char:   'o',
		Width:  1,
		Height: 1,
		X:      x,
		Y:      y,
	}

	multiplier.Draw()

	var health int = 5
	var speed = time.Second

	if !canclone {
		speed = 300 * time.Millisecond
	}

	go func() {
		for health != 0 && !KillAll {

			if damage(&multiplier, health, multiplier.Char, 5) {
				speed = 300 * time.Millisecond
				health--
			}

			if health != 5 && canclone {
				multiplier.Char = '%'
				multiplier.Draw()
				time.Sleep(500 * time.Millisecond)

				go SpawnMultiplier(multiplier.X, multiplier.Y+1, false)
				go SpawnMultiplier(multiplier.X, multiplier.Y-1, false)
				canclone = false

				multiplier.Char = 'o'
			}
			time.Sleep(400 * time.Millisecond)
		}

	}()

	go func() {

		time.Sleep(time.Second)

		for health != 0 && !KillAll {
			_ = gotoplayerpos(&multiplier)
			time.Sleep(speed)
		}

		multiplier.Destroy()

	}()

}

func SensingSpike(x, y int) {
	var spike libtxt.Object = libtxt.Object{
		Char:   '|',
		Width:  1,
		Height: 1,
		Static: true,
		X:      x,
		Y:      y,
	}

	go func() {
		for !KillAll {

			if spike.X == hero.GetPlayerPos('X') && spike.Y == hero.GetPlayerPos('Y') {
				time.Sleep(200 * time.Millisecond)
				spike.Draw()

				if spike.X == hero.GetPlayerPos('X') && spike.Y == hero.GetPlayerPos('Y') {
					hero.Damage(7)
				}
				time.Sleep(500 * time.Millisecond)
				spike.Destroy()
			}
			time.Sleep(50 * time.Millisecond)
		}
	}()
}

func TimedSpike(x, y int, ms time.Duration) {
	var spike libtxt.Object = libtxt.Object{
		Char:   '|',
		Width:  1,
		Height: 1,
		Static: true,
		X:      x,
		Y:      y,
	}

	go func() {
		for !KillAll {

			hx := hero.GetPlayerPos('X')
			hy := hero.GetPlayerPos('Y')

			spike.Draw()

			if spike.X == hx && spike.Y == hy {
				hero.Damage(7)
			}

			time.Sleep(200 * time.Millisecond)
			spike.Destroy()

			time.Sleep(ms * time.Millisecond)
		}
	}()
}

func SpawnStatue(x, y int) {
	var statue libtxt.Object = libtxt.Object{
		Char:   '&',
		Width:  1,
		Height: 1,
		Static: true,
		X:      x,
		Y:      y,
	}

	statue.Draw()

	var health int = 8
	var speed = 200 * time.Millisecond

	go func() {
		for health != 0 && !KillAll {

			if damage(&statue, health, statue.Char, 5) {
				health--
			}

			time.Sleep(400 * time.Millisecond)
		}
		statue.Destroy()
	}()

	go func() {

		for health != 0 && !KillAll {

			if health != 8 {
				statue.Static = false
				_ = gotoplayerpos(&statue)
			}
			time.Sleep(speed)
		}
	}()

}

func LoadEnemies(filename string) {
	file, err := enms.Open(filename)

	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var enemy string
	var enemyx, enemyy, ms int

	for scanner.Scan() {

		fmt.Sscanf(scanner.Text(), "%s %d %d %d", &enemy, &enemyx, &enemyy, &ms)

		if err != nil {
			panic(err)
		}

		switch enemy {
		case "walker":
			SpawnWalker(enemyx, enemyy)
		case "runner":
			SpawnRunner(enemyx, enemyy)
		case "archer":
			SpawnArcher(enemyx, enemyy)
		case "bat":
			SpawnBat(enemyx, enemyy)
		case "multiplier":
			SpawnMultiplier(enemyx, enemyy, true)
		case "statue":
			SpawnStatue(enemyx, enemyy)
		case "Sspike":
			SensingSpike(enemyx, enemyy)
		case "Tspike":
			TimedSpike(enemyx, enemyy, time.Duration(ms))
		default:
			panic(fmt.Sprintf("MISSINGNO.: %s", enemy))
		}
	}

	file.Close()
}

func Merto(x, y int) {
	var mage libtxt.Object = libtxt.Object{
		X:      x,
		Y:      y,
		Char:   '§',
		Static: false,
		Width:  1,
		Height: 1,
	}

	var delay int = 500

	mage.Draw()

	hero.CanMove = false

	time.Sleep(3 * time.Second)

	libtxt.Dialog("dialog/mage1")

	libtxt.RenderText("Merto: 500", 30, 1, 0)

	hero.CanMove = true

	var dir string

	go func() {
		health := 500
		for hero.Health > 10 {
			time.Sleep(time.Duration(delay) * time.Millisecond)
			health++

			if delay > 250 {
				delay -= 3
			}
			
			libtxt.RenderText(fmt.Sprintf("Merto: %d", health), 30, 1, 0)
		}
	}()

	for hero.Health > 10 {
		dir = gotoplayerpos(&mage)

		if dir == "left" {
			go projectile(mage.X-2, mage.Y, "left", true)
		} else if dir == "right" {
			go projectile(mage.X+2, mage.Y, "right", true)
		}

		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	KillAll = true

	hero.CanMove = false

	time.Sleep(time.Second)

	libtxt.Dialog("dialog/mage2")

	time.Sleep(time.Second)

	hero.CanMove = true

	mage.Destroy()

	time.Sleep(time.Second)

	libtxt.Screen.S[8][48] = 'D'

	libtxt.Update()
}
