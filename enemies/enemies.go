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

var mageh libtxt.Object = libtxt.Object{
	Char:   '^',
	Width:  1,
	Height: 1,
	X:      0,
	Y:      0,
}

var magel libtxt.Object = libtxt.Object{
	Char:   'W',
	Width:  1,
	Height: 1,
	X:      0,
	Y:      0,
}

var mager libtxt.Object = libtxt.Object{
	Char:   'o',
	Width:  1,
	Height: 1,
	X:      0,
	Y:      0,
}

var mager2 libtxt.Object = libtxt.Object{
	Char:   '|',
	Width:  1,
	Height: 1,
	X:      0,
	Y:      0,
}

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

	if col != "" && (col[0] == '-' || col[0] == '|') {
		enemy.Char = rune(48 + health-1)
		enemy.Draw()
		time.Sleep(200 * time.Millisecond)
		enemy.Destroy()
		enemy.Char = char
		enemy.Draw()
		return true
	}

	if col != "" && col[0] == '@' {
		hero.Damage(damage)
	}

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

			if damage(&bat, health, bat.Char, 0){
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

	var multiplier libtxt.Object = libtxt.Object {
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
	var spike libtxt.Object = libtxt.Object {
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

func SpawnMage(x, y int) {

	mageh.X = x
	mageh.Y = y
	magel.X = x
	magel.Y = y + 1
	mager.X = x + 1
	mager.Y = y
	mager2.X = x + 1
	mager2.Y = y + 1

	libtxt.Screen.S[y][x+5] = '♦'
	libtxt.Screen.S[y+1][x+5] = '♦'
	libtxt.Screen.S[y+2][x+5] = '♦'
	libtxt.Screen.S[y+3][x+5] = '♦'

	
	hero.CanMove = false

	magedraw("right")
	time.Sleep(3 * time.Second)
	magedraw("left")
	time.Sleep(time.Second)

	libtxt.RenderText("Ola, forasteiro.", x-10, y-1, 2000)
	libtxt.RenderText("Vejo que conseguiu", x-12, y-1, 1500)
	libtxt.RenderText("passar por todas as salas.", x-19, y-1, 3000)

	time.Sleep(2 * time.Second)

	libtxt.RenderText("Meus parabens, mas", x-15, y-1, 2000)
	libtxt.RenderText("infelizmente, este", x-15, y-1, 2000)
	libtxt.RenderText("e o fim da linha para voce", x-20, y-1, 3000)

	time.Sleep(time.Second)

	libtxt.RenderText("Ja parou para pensar o que", x-20, y-1, 2000)
	libtxt.RenderText("esses monstros fazem aqui", x-20, y-1, 2000)
	libtxt.RenderText("o dia inteiro?", x-10, y-1, 3000)
	libtxt.RenderText("Me protegem!", x-10, y-1, 2000)

	magedraw("right")
	time.Sleep(time.Second)

	libtxt.RenderText("Esta vendo essas joias aqui?", x-22, y-1, 3000)
	libtxt.RenderText("Sao as joias da terra, do fogo,", x-24, y-1, 3000)
	libtxt.RenderText("da agua e do ar.", x-15, y-1, 2000)
	libtxt.RenderText("As lendas dizem que aquele que", x-25, y-1, 3000)
	libtxt.RenderText("juntar todas elas e fizer um", x-25, y-1, 3000)
	libtxt.RenderText("feitico especifico se tornara um deus.", x-35, y-1, 3000)

	magedraw("left")
	time.Sleep(time.Second)

	libtxt.RenderText("Eu ainda nao descobri qual feitico", x-30, y-1, 3000)
	libtxt.RenderText("mas eu vou.", x-10, y-1, 3000)
	libtxt.RenderText("E voce nao vai me impedir", x-20, y-1, 2000)
	libtxt.RenderText("Fiz muito esforco para enfim", x-25, y-1, 3000)
	libtxt.RenderText("roubar todas elas.", x-15, y-1, 2000)
	libtxt.RenderText("Prepare-se para morrer.", x-20, y-1, 2000)

	hero.CanMove = true

	libtxt.Screen.S[y][x+5] = 0
	libtxt.Screen.S[y+1][x+5] = 0
	libtxt.Screen.S[y+2][x+5] = 0

	var health = 100
	var dead bool
	var visible bool = true

	go func() {
		var col string
		libtxt.RenderText("Quamoire: 100", 30, 1, 0)

		for health > 0 && !hero.Lost {

			if dead {
				break
			}

			if hero.Lost {
				KillAll = true
			}

			col = mageh.Coll("")
			if col != "" && visible {
				if col[0] == '-' {
					health -= 2
					libtxt.RenderText(fmt.Sprintf("Quamoire: %d", health), 30, 1, 0)
				} else if col[0] == '@' {
					hero.Damage(20)
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
		KillAll = true
	}()

	direct := ""

	for health > 70 && !hero.Lost {

		magedraw("right")
		visible = true

		for i := 0; i <= 3 && !hero.Lost; i++ {
			if hero.GetPlayerPos('X') < mageh.X {
				direct = "left"
				mageleft()
			} else if hero.GetPlayerPos('X') > mageh.X {
				direct = "right"
				mageright()
			}

			if hero.GetPlayerPos('Y') < magel.Y {
				mageup()
			} else if hero.GetPlayerPos('Y') > magel.Y {
				magedown(direct)
			}

			mageattack("magic", direct, mageh.X, mageh.Y)
			time.Sleep(2 * time.Second)
		}

		magedestroy()
		visible = false

		mageattack("spawn", "", hero.GetPlayerPos('X')-2, hero.GetPlayerPos('Y')-4)
	}

	magecenter()
	visible = true

	for i := 0; i <= 3 && !hero.Lost; i++ {
		mageattack("bullethell", "", 0, 0)
		time.Sleep(time.Second)
	}

	magedestroy()
	visible = false

	if !hero.Lost {
		mageattack("archers", "", 0, 0)
	}

	for health > 50 && !hero.Lost {
		magedraw("left")
		visible = true

		mageattack("spawn", "", hero.GetPlayerPos('X')-2, hero.GetPlayerPos('Y')-4)

		magedestroy()
		visible = false
		time.Sleep(3 * time.Second)
	}

	magedraw("right")
	visible = true

	if !hero.Lost && health > 0 {
		mageattack("spikes", "", 0, 0)
		mageattack("spikes", "", 0, 0)
	}

	for health > 0 && !hero.Lost {
		magedraw("right")
		visible = true

		mageattack("bullethell", "", 0, 0)
		time.Sleep(2 * time.Second)
		mageattack("spawn", "", hero.GetPlayerPos('X')-2, hero.GetPlayerPos('Y')-4)
		magedestroy()
		visible = false
		time.Sleep(2 * time.Second)
	}

	dead = true
	KillAll = true

	if hero.Lost {
		hero.CanMove = false
		magecenter()

		libtxt.RenderText("Foi bom te conhecer,", mageh.X-10, mageh.Y-1, 2000)
		libtxt.RenderText("mas preciso ir.", mageh.X-10, mageh.Y-1, 2000)
		libtxt.RenderText("Pouparei a sua vida pois", mageh.X-15, mageh.Y-1, 2000)
		libtxt.RenderText("voce foi forte o suficiente", mageh.X-20, mageh.Y-1, 2000)
		libtxt.RenderText("para chegar ate aqui.", mageh.X-20, mageh.Y-1, 4000)
		libtxt.RenderText("a saida esta ali.", mageh.X-10, mageh.Y-1, 1000)
		hero.CanMove = true
	}

	magedestroy()
	time.Sleep(2 * time.Second)
	libtxt.Screen.S[8][48] = 'D'
	libtxt.Update()
}

func mageleft() {

	mager.Destroy()
	mager2.Destroy()

	mageh.Move("left")
	mager.X = mageh.X - 1
	magel.Move("left")
	mager2.X = mageh.X - 1

	mager.Draw()
	mager2.Draw()
}

func mageright() {

	mager.Destroy()
	mager2.Destroy()

	mageh.Move("right")
	mager.X = mageh.X + 1
	magel.Move("right")
	mager2.X = mageh.X + 1

	mager.Draw()
	mager2.Draw()
}

func mageup() {
	mageh.Move("up")
	magel.Move("up")
	mager.Move("up")
	mager2.Move("up")
}

func magedown(direct string) {

	magedestroy()

	mager.Y++
	mager2.Y++
	mageh.Y++
	magel.Y++

	magedraw(direct)
}

func magedestroy() {
	mageh.Destroy()
	magel.Destroy()
	mager.Destroy()
	mager2.Destroy()
}

func magedraw(direction string) {

	mager.Destroy()
	mager2.Destroy()

	magel.X = mageh.X
	magel.Y = mageh.Y + 1

	if direction == "left" {
		mager.X = mageh.X - 1
		mager.Y = mageh.Y

		mager2.X = mageh.X - 1
		mager2.Y = mageh.Y + 1
	} else {
		mager.X = mageh.X + 1
		mager.Y = mageh.Y

		mager2.X = mageh.X + 1
		mager2.Y = mageh.Y + 1
	}

	mageh.Draw()
	magel.Draw()
	mager.Draw()
	mager2.Draw()
}

func magecenter() {
	magedestroy()

	mageh.X = 25
	mageh.Y = 8

	mageh.Draw()
	time.Sleep(1000 * time.Millisecond)

	magedraw("right")
}

func mageattack(magic string, direction string, x, y int) {

	rand.Seed(time.Now().UnixNano())

	switch magic {
	case "spawn":

		enms := []string{"walker", "runner", "multiplier"}

		spawn := enms[rand.Intn(2)]

		if spawn == "walker" {
			for i := 0; i < 5 && y+i < 14 && y+i > 2; i++ {
				SpawnWalker(x, y+i)
			}
		} else if spawn == "runner" {
			for i := 0; i < 5 && y+i < 14 && y+i > 2; i++ {
				SpawnRunner(x, y+i)
			}
		} else {
			for i := 0; i < 5 && y+i < 14 && y+i > 2; i++ {
				SpawnMultiplier(x, y+i, true)
			}
		}

		time.Sleep(20 * time.Second)
		KillAll = true
		time.Sleep(3 * time.Second)
		KillAll = false
	case "bullethell":
		for i := 13; i < 34; i += 2 {
			go projectile(i, 10, "down", true)
			go projectile(i, 6, "up", true)
		}

		for i := 4; i < 14; i += 2 {
			go projectile(12, i, "left", true)
			go projectile(34, i, "right", true)
		}
	case "magic":

		time.Sleep(time.Second)

		switch direction {
		case "up":
			go projectile(x, y-2, direction, true)
			go projectile(x-1, y-2, direction, true)
			go projectile(x+1, y-2, direction, true)
		case "down":
			go projectile(x, y+2, direction, true)
			go projectile(x+1, y+2, direction, true)
			go projectile(x-1, y+2, direction, true)
		case "left":
			go projectile(x-3, y, direction, true)
			go projectile(x-3, y+1, direction, true)
			go projectile(x-3, y-1, direction, true)
		case "right":
			go projectile(x+4, y, direction, true)
			go projectile(x+4, y+1, direction, true)
			go projectile(x+4, y-1, direction, true)
		}

	case "spikes":

		for i := 0; i <= 50; i++ {
			go TimedSpike(rand.Intn(48), 3+rand.Intn(11), 500)
			time.Sleep(500 * time.Millisecond)
		}
	case "archers":

		for i := 0; i < 11; i++ {
			go SpawnArcher(2, 3+i)
			go SpawnArcher(35, 3+i)
		}
		time.Sleep(20 * time.Second)
		KillAll = true
		time.Sleep(3 * time.Second)
		KillAll = false
	}

	switch magic {
	case "spawn", "bullethell", "spikes", "archers":

		mager.Destroy()
		mager2.Destroy()

		mager.Y--
		mager2.Y--
		mager.Draw()
		mager2.Draw()

		time.Sleep(400 * time.Millisecond)

		mager.Destroy()
		mager2.Destroy()

		mager.Y++
		mager2.Y++
		mager.Draw()
		mager2.Draw()
	case "magic":
		if direction == "left" {
			mager.Destroy()
			mager2.Destroy()
			mager.X--
			mager2.Char = '\\'
			mager.Draw()

			mager2.Draw()

			time.Sleep(300 * time.Millisecond)

			mager.Destroy()
			mager.X++
			mager2.Char = '|'
			mager.Draw()
			mager2.Draw()

		} else if direction == "right" {
			mager.Destroy()
			mager2.Destroy()
			mager.X++
			mager2.Char = '/'
			mager.Draw()
			mager2.Draw()

			time.Sleep(300 * time.Millisecond)

			mager.Destroy()
			mager.X--
			mager2.Char = '|'
			mager.Draw()
			mager2.Draw()
		}

	}
}
