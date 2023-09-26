package enemies

import (
	"bufio"
	"deepdungeons/hero"
	"deepdungeons/items"
	"deepdungeons/libtxt"
	"embed"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"sync"
)

//go:embed enms
var enms embed.FS

type Enemy struct {
	libtxt.Object
	health int
	speed  time.Duration
	damage int
	Transition bool
}

var EnemiesAlive []*Enemy
var EnemiesGroup sync.WaitGroup

func KillAll(){
	for _, enemy := range(EnemiesAlive) {
		enemy.Transition = true
		enemy.Kill()
	}

	EnemiesGroup.Wait()
}

func contains(enemy *Enemy) bool {
	for _, elm := range EnemiesAlive {
		if elm == enemy { return true }
	}

	return false
}

func (enemy *Enemy) GotoPlayerPos() string {

	var direction string

	if enemy.X < hero.GetPlayerPos('X') {
		direction = "right"
	} else if enemy.X > hero.GetPlayerPos('X') {
		direction = "left"
	}

	enemy.Move(direction)

	if enemy.Y < hero.GetPlayerPos('Y') {
		direction = "down"
	} else if enemy.Y > hero.GetPlayerPos('Y') {
		direction = "up"
	}

	if enemy.Y != hero.GetPlayerPos('Y') {
		enemy.Move(direction)
	}

	return direction
}

func (enemy *Enemy) Damage(damage int) bool {

	var col = enemy.Coll("")
	var char = enemy.Char
	var hit = false

	if strings.Contains(col, "-") || strings.Contains(col, "|") {
		hit = true

		enemy.Char = rune(48 + enemy.health - 1)

		enemy.Draw()
		time.Sleep(200 * time.Millisecond)
		enemy.Destroy()

		enemy.Char = char
		enemy.Draw()
	}

	if strings.Contains(col, "@") { hero.Damage(enemy.damage) }

	return hit
}

func (enemy *Enemy) Kill() {
	enemy.health = 0
	enemy.Destroy()
}

func (enemy *Enemy) DamageRoutine() {
	EnemiesGroup.Add(1)
	for enemy.health != 0 {
		if !contains(enemy) {
			EnemiesAlive = append(EnemiesAlive, enemy)
		}

		if enemy.Damage(enemy.damage) { enemy.health-- }

		time.Sleep(400 * time.Millisecond)
	}
	EnemiesGroup.Done()
	enemy.Destroy()
}

func (enemy *Enemy) Walk() {
	_ = enemy.GotoPlayerPos()
	time.Sleep(enemy.speed)
}

func (enemy *Enemy) Cycle() {
	go enemy.DamageRoutine()
	for enemy.health != 0 {
		if enemy.health != 0 { enemy.Walk() }
	}
}

func SpawnWalker(x, y int) {

	var walker Enemy = Enemy{
		Object: libtxt.Object{
			Char:   '$',
			Width:  1,
			Height: 1,
			X:      x,
			Y:      y,
		},
		health: 2,
		speed:  time.Second,
		damage: 1,
	}
	walker.Draw()

	go func() { walker.Cycle() }()
}

func SpawnRunner(x, y int) {

	var runner Enemy = Enemy{
		Object: libtxt.Object{
			Char:   '/',
			Width:  1,
			Height: 1,
			X:      x,
			Y:      y,
		},
		health: 3,
		speed:  300 * time.Millisecond,
		damage: 2,
	}
	runner.Draw()

	go func() { runner.Cycle() }()
}

func projectile(x, y int, direction string, magic bool) {
	var p Enemy = Enemy{
		Object: libtxt.Object{
			Char:   '-',
			Width:  1,
			Height: 1,
			X:      x,
			Y:      y,
		},
		damage: 5,
		speed:  100 * time.Millisecond,
		health: 1,
	}

	var offset = -2

	if direction == "right" { offset += 4 }
	p.X += offset

	if magic {
		p.damage = 10
		p.Char = '*'
	}

	EnemiesAlive = append(EnemiesAlive, &p)

	var collision string = p.Coll(direction)

	for collision == "" && p.health != 0 {
		collision = p.Coll(direction)
		p.Move(direction)
		time.Sleep(p.speed)
	}

	if strings.Contains(collision, "@") { hero.Damage(p.damage) }

	p.Destroy()
}

func SpawnArcher(x, y int) {
	var archer Enemy = Enemy{
		Object: libtxt.Object{
			Char:   '<',
			Width:  1,
			Height: 1,
			X:      x,
			Y:      y,
		},
		health: 3,
		speed:  500 * time.Millisecond,
		damage: 3,
	}

	archer.Draw()

	var herox int = hero.GetPlayerPos('X')
	var dir string

	go archer.DamageRoutine()

	go func() {

		for archer.health != 0 {

			herox = hero.GetPlayerPos('X')

			if herox-archer.X >= -5 {
				archer.Char = '>'

				for i := 0; i < 10 && archer.health != 0; i++ {
					time.Sleep(archer.speed)
					archer.Move("right")
				}

				time.Sleep(archer.speed)
			}

			herox = hero.GetPlayerPos('X')

			if archer.X-herox <= 0 && archer.X-herox >= -5 {
				archer.Char = '<'

				for i := 0; i < 10 && archer.health != 0; i++ {
					time.Sleep(archer.speed)
					archer.Move("left")
				}

				time.Sleep(archer.speed)
			}

			if archer.X > herox {
				archer.Char = '<'
				dir = "left"
			} else {
				archer.Char = '>'
				dir = "right"
			}
			archer.Draw()
			archer.Move(dir)

			for i := 0; i < 5 && archer.health != 0; i++ {
				projectile(archer.X, archer.Y, dir, false)
				time.Sleep(time.Second)
			}
		}
		archer.Kill()
		if !archer.Transition { items.Chest(archer.X, archer.Y, "arrow") }

	}()
}

func SpawnBat(x, y int) {
	var bat Enemy = Enemy{
		Object: libtxt.Object{
			Char:   '|',
			Width:  1,
			Height: 1,
			X:      x,
			Y:      y,
		},
		health: 1,
		speed:  200 * time.Millisecond,
		damage: 1,
	}

	bat.Draw()

	go bat.DamageRoutine()

	go func() {

		directions := []string{"up", "down", "left", "right"}

		direction := "left"
		rand.Seed(time.Now().UnixNano())

		hx := hero.GetPlayerPos('X')

		for hx+3 != bat.X && hx-3 != bat.X {
			hx = hero.GetPlayerPos('X')
			time.Sleep(50 * time.Millisecond)
		}

		bat.Char = '^'

		var previousx, previousy int

		for bat.health != 0 {
			for i := 0; i <= 10 && bat.health != 0; i++ {
				previousx, previousy = bat.X, bat.Y
				bat.Move(direction)

				if previousx == bat.X && previousy == bat.Y {
					direction = directions[rand.Intn(4)]
					bat.Move(direction)
				}

				time.Sleep(bat.speed)
			}
			direction = directions[rand.Intn(4)]
		}
		bat.Destroy()
	}()
}

func SpawnMultiplier(x, y int, canclone bool) {

	var multiplier Enemy = Enemy{
		Object: libtxt.Object{
			Char:   'o',
			Width:  1,
			Height: 1,
			X:      x,
			Y:      y,
		},
		health: 5,
		speed:  time.Second,
		damage: 5,
	}

	multiplier.Draw()

	if !canclone {
		multiplier.speed = 300 * time.Millisecond
	}

	go multiplier.Cycle()

	go func() {
		for multiplier.health != 0 {
			if multiplier.health != 5 && canclone {
				multiplier.Char = '%'
				multiplier.Draw()
				time.Sleep(500 * time.Millisecond)

				go SpawnMultiplier(multiplier.X, multiplier.Y+1, false)
				go SpawnMultiplier(multiplier.X, multiplier.Y-1, false)
				canclone = false

				multiplier.Char = 'o'
			}
			time.Sleep(multiplier.speed)
		}

	}()

}

func SensingSpike(x, y int) {
	var spike Enemy = Enemy{
		Object: libtxt.Object{
			Char:   '|',
			Width:  1,
			Height: 1,
			Static: true,
			X:      x,
			Y:      y,
		},
		health: 1,
		speed:  50 * time.Millisecond,
		damage: 7,
	}
	EnemiesAlive = append(EnemiesAlive, &spike)

	go func() {
		for spike.health == 1 {

			if spike.X == hero.GetPlayerPos('X') && spike.Y == hero.GetPlayerPos('Y') {
				time.Sleep(200 * time.Millisecond)
				spike.Draw()

				if spike.X == hero.GetPlayerPos('X') && spike.Y == hero.GetPlayerPos('Y') {
					hero.Damage(spike.damage)
				}

				time.Sleep(500 * time.Millisecond)
				spike.Destroy()
			}
			time.Sleep(50 * time.Millisecond)
		}
	}()
}

func TimedSpike(x, y int, ms time.Duration) {
	var spike Enemy = Enemy{
		Object: libtxt.Object{
			Char:   '|',
			Width:  1,
			Height: 1,
			Static: true,
			X:      x,
			Y:      y,
		},
		health: 1,
		speed:  ms * time.Millisecond,
		damage: 7,
	}
	EnemiesAlive = append(EnemiesAlive, &spike)

	go func() {
		for spike.health == 1 {
			hx := hero.GetPlayerPos('X')
			hy := hero.GetPlayerPos('Y')

			spike.Draw()

			if spike.X == hx && spike.Y == hy {
				hero.Damage(spike.damage)
			}

			time.Sleep(200 * time.Millisecond)
			spike.Destroy()

			time.Sleep(ms * time.Millisecond)
		}
	}()
}

func SpawnStatue(x, y int) {
	var statue Enemy = Enemy{
		Object: libtxt.Object{
			Char:   '&',
			Width:  1,
			Height: 1,
			Static: true,
			X:      x,
			Y:      y,
		},
		health: 8,
		speed:  200 * time.Millisecond,
		damage: 7,
	}

	statue.Draw()

	go statue.DamageRoutine()

	go func() {
		for statue.health != 0 {
			if statue.health != 8 {
				statue.Static = false
				statue.Walk()
			}
		}
	}()

}

func LoadEnemies(filename string) {
	file, err := enms.Open(filename)

	if err != nil { panic(err) }

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var enemy string
	var enemyx, enemyy, ms int

	for scanner.Scan() {

		fmt.Sscanf(scanner.Text(), "%s %d %d %d", &enemy, &enemyx, &enemyy, &ms)

		if err != nil { panic(err) }

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
	var mage Enemy = Enemy{
		Object: libtxt.Object {
			X:      x,
			Y:      y,
			Char:   '§',
			Static: false,
			Width:  1,
			Height: 1,
		},
		speed: 500 * time.Millisecond,
	}

	mage.Draw()

	hero.CanMove = false
	time.Sleep(3 * time.Second)

	libtxt.Dialog("dialog/mage1")

	libtxt.RenderText("Merto: 500", 30, 1, 0)
	hero.CanMove = true

	go func() {
		health := 500
		for hero.Health > 10 {
			time.Sleep(mage.speed)
			health++

			if mage.speed > 250 { mage.speed -= 3 }

			libtxt.RenderText(fmt.Sprintf("Merto: %d", health), 30, 1, 0)
		}
	}()

	for hero.Health > 10 {
		dir := mage.GotoPlayerPos()
		go projectile(mage.X, mage.Y, dir, true)

		time.Sleep(mage.speed)
	}

	KillAll()
	EnemiesGroup.Wait()

	hero.CanMove = false

	time.Sleep(time.Second)
	libtxt.Dialog("dialog/mage2")
	time.Sleep(time.Second)

	hero.CanMove = true
	mage.Destroy()
	libtxt.RenderText("            ", 30, 1, 1)
	time.Sleep(time.Second)

	libtxt.Screen.S[8][48] = 'D'
	libtxt.Update()
}
