package libtxt

import (
	"bufio"
	"embed"
	"fmt"
	"github.com/mattn/go-tty"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type scr struct {
	m sync.Mutex
	S [15][]rune
}

var Updated bool
var Screen scr

//go:embed lvls
var lvls embed.FS

//go:embed dialog
var dia embed.FS

func Init() {

	if runtime.GOOS == "windows" {
		regcmd := exec.Command(os.ExpandEnv("$windir\\System32\\reg.exe"), "add", "HKCU\\Console", "/v", "VirtualTerminalLevel", "/t", "REG_DWORD", "/d", "1", "/f")

		regcmd.Run()
	}

	for i := range Screen.S {
		Screen.S[i] = make([]rune, 50)
	}
}

type Object struct {
	Char   rune
	Width  int
	Height int
	Static bool
	X      int
	Y      int
}

func (o *Object) Move(direction string) {

	if o.Static {
		return
	}

	o.Destroy()

	switch direction {
	case "up":
		if o.Coll("up") == "" {
			o.Y--
		}
	case "down":
		if o.Coll("down") == "" {
			o.Y++
		}
	case "left":
		if o.Coll("left") == "" {
			o.X--
		}
	case "right":
		if o.Coll("right") == "" {
			o.X++
		}
	default:
		panic("???")
	}

	o.Draw()
}

func offlimits(o Object) bool {
	testx := o.X < 0 || o.X > 49
	testy := o.Y > 15 || o.Y < 0

	return testx || testy
}

func (o Object) Draw() {

	if offlimits(o) {
		return
	}

	Screen.m.Lock()

	for i := 1; i <= o.Height-1; i++ {
		Screen.S[o.Y+i][o.X+1] = o.Char
	}

	for i := 1; i <= o.Width; i++ {
		Screen.S[o.Y][o.X+i] = o.Char
	}

	Screen.m.Unlock()

	Updated = true
}

func (o Object) Destroy() {

	if offlimits(o) {
		return
	}

	Screen.m.Lock()

	for i := 0; i < o.Height; i++ {
		Screen.S[o.Y+i][o.X+1] = 0
	}

	for i := 0; i < o.Width; i++ {
		Screen.S[o.Y][o.X+i+1] = 0
	}

	Screen.m.Unlock()

	Updated = true
}

func (o Object) Coll(direction string) string {

	var left = Screen.S[o.Y][o.X-o.Width+1]
	var right = Screen.S[o.Y][o.X+o.Width+1]
	var up = Screen.S[o.Y-o.Height][o.X+1]
	var down = Screen.S[o.Y+o.Height][o.X+1]

	if direction == "" {
		if right != 0 {
			return fmt.Sprintf("%c/right", right)
		} else if left != 0 {
			return fmt.Sprintf("%c/left", left)
		}

		if down != 0 {
			return fmt.Sprintf("%c/down", down)
		} else if up != 0 {
			return fmt.Sprintf("%c/up", up)
		} else {
			return ""
		}
	}

	switch direction {
	case "up":
		if o.Y < 3 {
			return "barrier"
		} else if up != 0 {
			return string(up)
		} else {
			return ""
		}
	case "down":
		if o.Y > len(Screen.S)-3 {
			return "barrier"
		} else if down != 0 {
			return string(down)
		} else {
			return ""
		}
	case "left":
		if o.X == 0 {
			return "barrier"
		} else if left != 0 {
			return string(left)
		} else {
			return ""
		}
	case "right":
		if o.X > len(Screen.S[0])-4 {
			return "barrier/right"
		} else if right != 0 {
			return string(right)
		} else {
			return ""
		}

	}

	return ""
}

func Getkeystroke() rune {
	tty, err := tty.Open()
	if err != nil {
		panic(err)
	}
	defer tty.Close()

	key, err := tty.ReadRune()

	if err != nil {
		panic(err)
	}

	return key
}

func Update() {

	fmt.Printf("\033[H\033[2J")

	for _, x := range Screen.S {
		for _, c := range x {

			if c == 0 {
				fmt.Print(" ")
			} else {
				fmt.Printf("%c", c)
			}
		}
		fmt.Print("\n")
	}
}

func RenderText(text string, x, y int, ms time.Duration) {
	for i, char := range text {
		Screen.S[y][x+i+1] = char
	}

	Screen.S[y][x+len(text)+1] = 0

	Updated = true

	time.Sleep(ms * time.Millisecond)

	for i := 0; i <= len(text) && ms != 0; i++ {
		Screen.S[y][x+i+1] = 0
	}

	Updated = true
}

func LoadScene(filename string) {

	for y := 0; y < 15; y++ {
		for x := 0; x < 50; x++ {
			Screen.S[y][x] = 0
		}
	}

	var scene, err = lvls.Open(filename)

	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(scene)
	scanner.Split(bufio.ScanLines)

	wx := ""
	wy := ""

	var repeatx, repeaty bool

	var ix, iy int

	for scanner.Scan() {

		fmt.Sscanf(scanner.Text(), "%s %s", &wx, &wy)

		repeatx = strings.Contains(wx, "r")
		repeaty = strings.Contains(wy, "r")

		if repeatx {
			wx = strings.Replace(wx, "r", "", 1)
			startx := strings.Split(wx, "-")[0]

			ix, err = strconv.Atoi(startx)
			wx = strings.Split(wx, "-")[1]
		}

		if repeaty {
			wy = strings.Replace(wy, "r", "", 1)
			starty := strings.Split(wy, "-")[0]

			iy, err = strconv.Atoi(starty)
			wy = strings.Split(wy, "-")[1]
		}

		if err != nil {
			panic(err)
		}

		x, err := strconv.Atoi(wx)

		if err != nil {
			panic(err)
		}

		y, err := strconv.Atoi(wy)

		if err != nil {
			panic(err)
		}

		for repeatx && ix < x {
			Screen.S[y][ix] = '#'
			ix++
		}

		for repeaty && iy < y {
			Screen.S[iy][x] = '#'
			iy++
		}
		if !(repeatx || repeaty) {
			Screen.S[y][x] = '#'
		}
	}

	scene.Close()

}

func indexof(barray []byte, byt byte) int {
	for index, b := range barray {
		if b == byt {
			return index
		}
	}
	return -1
}

func Dialog(filename string) {

	file, err := dia.Open(filename)

	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)

	var line string

	for scanner.Scan() {
		line = scanner.Text()

		for _, char := range strings.Split(line, "") {
			time.Sleep(90 * time.Millisecond)
			fmt.Print(char)
		}
		fmt.Println()
		time.Sleep(2 * time.Second)
	}
}
