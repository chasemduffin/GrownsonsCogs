package counting-report

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

type user struct {
	name string
	counts int
	mistakes int
	streak int
}

func main() {
	var input *csv.Reader
	if len(os.Args) > 1 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
		input = csv.NewReader(file)
	} else {
		input = csv.NewReader(os.Stdin)
	}
	firstLine := true
	users := make(map[string]user)
	last := 0
	var lastUser string
	streak := 1
	for {
		record, err := input.Read()
		if firstLine {
			firstLine = false
			continue
		}
		if err == io.EOF {
			break
		}
		if record != nil {
			mistake := false
			if len(record) < 4 {
				continue
			}
			userName := asciiName(record[1])
			date := record[2]
			countAttempt := record[3]
			current, err := strconv.Atoi(countAttempt)
			if userName != lastUser {
				streak = 1
			}
			if err != nil {
				mistake = true
				fmt.Printf("[%s] %s: '%s' is not a valid number\n", date, userName, countAttempt)
				streak = 1
			} else if current != last + 1 {
				mistake = true
				fmt.Printf("[%s] %s: %d does not follow %d\n", date, userName, current, last)
				streak = 1
				last = current
			} else {
				last = current
				if lastUser == userName {
					streak++
				}
				lastUser = userName
			}
			currentUser, ok := users[userName]
			if !ok {
				currentUser = user{
					name:     userName,
					counts:   0,
					mistakes: 0,
					streak:   1,
				}
			}
			currentUser.counts++
			if mistake {
				currentUser.mistakes++
			}
			if streak > currentUser.streak {
				currentUser.streak = streak
			}
			users[userName] = currentUser
		}
	}
	var totalCounts, totalMistakes, longestStreak int
	var names []string
	namePad, mistakePad, countPad, streakPad := 5, 8, 6, 6
	for name, chatter := range users {
		names = append(names, name)
		if len(chatter.name) > namePad {
			namePad = len(chatter.name)
		}
		if chatter.streak > longestStreak {
			longestStreak = chatter.streak
		}
		totalCounts += chatter.counts
		totalMistakes += chatter.mistakes
	}
	if totalMistakes > 99999999 {
		mistakePad = len(strconv.Itoa(totalMistakes))
	}
	if totalCounts > 999999 {
		countPad = len(strconv.Itoa(totalCounts))
	}
	sort.Strings(names)
	lineFormat := fmt.Sprintf("%%-%ds %%%dv %%%dv %%%dv\n", namePad, countPad, mistakePad, streakPad)
	fmt.Printf(lineFormat, "Name", "Counts", "Mistakes", "Streak")
	for _, name := range names {
		chatter := users[name]
		fmt.Printf(lineFormat, chatter.name, chatter.counts, chatter.mistakes, chatter.streak)
	}
	fmt.Printf(lineFormat, "Total", totalCounts, totalMistakes, longestStreak)
	os.Exit(0)
}

func asciiName(in string) string {
	var out string
	for _, runeVal := range in {
		letter, err := strconv.Unquote(strconv.QuoteRune(runeVal))
		if err != nil {
			continue
		}
		if letter == "#" {
			break
		}
		if runeVal > 126 || runeVal < 32 {
			continue
		}
		out += letter
	}
	return out
}
