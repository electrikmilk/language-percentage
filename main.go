package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

var types map[string]float64
var percentages = map[string]string{}
var fileCount float64
var languages map[string]string
var initialPath string

func main() {
	languages = make(map[string]string)
	types = make(map[string]float64)
	percentages = make(map[string]string)
	mapLanguages()
	if len(os.Args) < 2 {
		cwd, err := os.Getwd()
		handle(err)
		initialPath = cwd
	} else {
		initialPath = os.Args[1]
	}
	getFiles(initialPath)
	sortedTypes := sortByValue(types)
	for typ, count := range sortedTypes {
		var percentage float64 = (count / fileCount) * 100
		var percentString string = fmt.Sprintf("%f", percentage)
		if percentage > 0.5 {
			var percent string = percentString[0:4] + "%"
			if percentage < 10 {
				percent = Style(percent, DIM)
			}
			if percentage > 10 && percentage < 25 {
				percent = Style(percent, CYAN)
			}
			if percentage > 10 && percentage < 25 {
				percent = Style(percent, BLUE)
			}
			if percentage > 25 && percentage < 50 {
				percent = Style(percent, YELLOW)
			}
			if percentage > 50 {
				percent = Style(percent, GREEN)
			}

			percentages[typ] = percent
		}
	}
	if len(os.Args) > 2 {
		fmt.Printf("%s:\n", initialPath)
		banner()
		fmt.Printf("\n")
	}
	for lang, percentage := range percentages {
		fmt.Printf("%s %s\n", percentage, languageName(lang))
	}
}

func banner() {
	bytes, err := exec.Command("tput", "cols").Output()
	if err != nil {
		panic(err)
	}
	cols, _ := strconv.Atoi(strings.TrimSuffix(string(bytes), "\n"))
	for i := 0; i < cols; i++ {
		fmt.Printf("-")
	}
}

func mapLanguages() {
	languages["js"] = "JavaScript"
	languages["ts"] = "TypeScript"
	languages["py"] = "Python"
	languages["cpp"] = "C++"
	languages["rs"] = "Rust"
	languages["md"] = "Markdown"
	languages["txt"] = "Plain Text"
	languages["go"] = "Go"
	languages["ttf"] = "TrueType Font"
	languages["woff"] = "WOFF Font"
	languages["eot"] = "Embedded OpenType"
	languages["less"] = "Less"
}

func languageName(language string) (name string) {
	if fullName, ok := languages[language]; ok {
		name = fullName
	} else {
		name = strings.ToUpper(language)
	}
	return
}

func getFiles(path string) {
	items, dirErr := os.ReadDir(path)
	handle(dirErr)
	for _, item := range items {
		if !item.IsDir() {
			if strings.Contains(item.Name(), ".") {
				var nameChars []string = strings.Split(item.Name(), "")
				if nameChars[0] != "." {
					var ext []string = strings.Split(item.Name(), ".")
					if len(ext) > 2 {
						increment(ext[2])
					} else {
						increment(ext[1])
					}
					fileCount++
				}
			}
		} else {
			getFiles(fmt.Sprintf("%s/%s", path, item.Name()))
		}
	}
	return
}

func increment(language string) {
	if count, ok := types[language]; ok {
		count++
		types[language] = count
	} else {
		types[language] = 1
	}
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func sortByValue(m map[string]float64) (sortedMap map[string]float64) {
	pl := make(PairList, len(m))
	i := 0
	for k, v := range m {
		pl[i] = Pair{k, int(v)}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	sortedMap = make(map[string]float64)
	for _, kv := range pl {
		sortedMap[kv.Key] = float64(kv.Value)
	}
	return
}

func handle(e error) {
	if e != nil {
		panic(e)
	}
}
