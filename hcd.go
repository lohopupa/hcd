package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var CD_HISTORY_PATH = "/home/lohopupa/.cd_history"

func startTui(paths []string) {
	app := tview.NewApplication()

	list := tview.NewList().
		ShowSecondaryText(false).
		SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
			app.Stop()
			println("Selected:", mainText)
		})
	for _, path := range paths {
		list.AddItem(path, "", 0, nil)
	}

	inputField := tview.NewInputField().
		SetLabel("Search: ").
		SetFieldWidth(0).
		SetChangedFunc(func(text string) {
			list.Clear()
			for _, item := range sortByLevenshteinDistance(paths, text) {
				list.AddItem(item, "", 0, nil)
			}
			if list.GetItemCount() > 0 {
				list.SetCurrentItem(0)
			}
		})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(inputField, 3, 1, true).
		AddItem(list, 0, 8, false)

	app.SetRoot(flex, true).SetFocus(inputField)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			app.Stop()
		case tcell.KeyUp:
			if list.GetCurrentItem() > 0 {
				list.SetCurrentItem(list.GetCurrentItem() - 1)
			}
		case tcell.KeyDown:
			if list.GetCurrentItem() < list.GetItemCount()-1 {
				list.SetCurrentItem(list.GetCurrentItem() + 1)
			}
		case tcell.KeyEnter:
			app.Stop()
			p, _ := list.GetItemText(list.GetCurrentItem())
			fmt.Println(p)
			cd(p)
		}
		
		return event
	})

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func cd(path string) error {
	err := os.Chdir(path)
	if err != nil {
		return err
	}
	return nil
}

func loadPaths(filepath string) ([]string, error) {
    content, err := os.ReadFile(filepath)
    if err != nil {
        return nil, err
    }

    paths := strings.Split(strings.TrimSpace(string(content)), "\n")

    uniquePaths := make(map[string]struct{})
    for _, path := range paths {
        uniquePaths[path] = struct{}{}
    }

    result := make([]string, 0, len(uniquePaths))
    for path := range uniquePaths {
        result = append(result, path)
    }

    return result, nil
}

func main() {
	paths, err := loadPaths(CD_HISTORY_PATH)
	if err != nil {
		log.Fatalf("Could not read file %s: %s", CD_HISTORY_PATH, err)
	}
	startTui(paths)
}


func levenshtein(a, b string) int {
	lenA := len(a)
	lenB := len(b)

	dp := make([][]int, lenA+1)
	for i := range dp {
		dp[i] = make([]int, lenB+1)
	}

	for i := 0; i <= lenA; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= lenB; j++ {
		dp[0][j] = j
	}

	for i := 1; i <= lenA; i++ {
		for j := 1; j <= lenB; j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			dp[i][j] = min(dp[i-1][j]+1, min(dp[i][j-1]+1, dp[i-1][j-1]+cost))
		}
	}

	return dp[lenA][lenB]
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func sortByLevenshteinDistance(stringsList []string, query string) []string {
	type pair struct {
		str      string
		distance int
	}
	pairs := make([]pair, len(stringsList))

	for i, str := range stringsList {
		pairs[i] = pair{str, levenshtein(str, query)}
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].distance < pairs[j].distance
	})

	sortedList := make([]string, len(stringsList))
	for i, p := range pairs {
		sortedList[i] = p.str
	}

	return sortedList
}
