package ucell_tpl_match

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	patternKeyDigit = 0
	patternKeyWord  = 1
)

var keyLetter = map[int]string{
	patternKeyWord:  "w",
	patternKeyDigit: "d",
}

var wordRangeEqual = regexp.MustCompile("^%w{1,([0-9]+)}$")
var digitRangeEqual = regexp.MustCompile("^%d{1,([0-9]+)}$")

var wordRangeContains = regexp.MustCompile("%w{1,([0-9]+)}")
var digitRangeContains = regexp.MustCompile("%d{1,([0-9]+)}")

var isDigit = regexp.MustCompile("^[0-9]+$")

type patternItem struct {
	max      int
	prefix   string
	suffix   string
	children *templateItem
}

type templateItem struct {
	children        map[string]*templateItem
	patternChildren map[int][]*patternItem

	isEnd bool
}

type ucellTemplate struct {
	items *templateItem
}

func NewUcellTemplate(tpls ...string) *ucellTemplate {
	ut := &ucellTemplate{
		items: newTemplateItem(true),
	}

	for _, tpl := range tpls {
		ut.Add(tpl)
	}

	return ut
}

func (t *ucellTemplate) Add(tpl string) {
	cleanTpl := cleanTemplate(tpl)

	tplWords := strings.Split(strings.ToLower(cleanTpl), " ")

	current := t.items

	for idx, word := range tplWords {
		isLastWord := idx == len(tplWords)-1

		switch {
		case digitRangeEqual.MatchString(word):
			n, _ := strconv.Atoi(digitRangeEqual.FindStringSubmatch(word)[1])
			current = addStaticPattern(current, patternKeyDigit, n, isLastWord)
		case digitRangeContains.MatchString(word):
			n, _ := strconv.Atoi(digitRangeContains.FindStringSubmatch(word)[1])
			current = addRegexPattern(current, patternKeyDigit, word, n, isLastWord)
		case wordRangeEqual.MatchString(word):
			n, _ := strconv.Atoi(wordRangeEqual.FindStringSubmatch(word)[1])
			current = addStaticPattern(current, patternKeyWord, n, isLastWord)
		case wordRangeContains.MatchString(word):
			n, _ := strconv.Atoi(wordRangeContains.FindStringSubmatch(word)[1])
			current = addRegexPattern(current, patternKeyWord, word, n, isLastWord)
		default:
			if c, ok := current.children[word]; ok {
				current = c
			} else {
				item := newTemplateItem(isLastWord)
				current.children[word] = item
				current = item
			}
		}
	}

	//salom %d dunyo
	//salom %d
	//kabi ikkita shablon kelsa, birinchi holatda %d isEnd=false bo'ladi
	//shuning uchun ushbu kod, isEnd=false qilib yaratilganni true ga
	//yani %d ni true aylantirish uchun kerak
	current.isEnd = true
}

func (t *ucellTemplate) IsMatch(message string) bool {
	cleanMsg := strings.ToLower(cleanMessage(message))

	if len(cleanMsg) == 0 {
		//Agar empty string kelsa
		//shablonlar ham yo'qligiga ishonch hosil qilish kerak
		return t.items == nil
	}

	words := strings.Split(cleanMsg, " ")
	return isMatch(t.items, words)
}

func isMatch(current *templateItem, words []string) bool {
	if len(words) == 0 {
		return current.isEnd
	}

	//Agar suffix kelsa, demak %d yoki %w shablon bo'yicha tekshirish boshlanga
	//Shu so'zning o'zi borligiga tekshiramiz
	if c, ok := current.children[words[0]]; ok {
		if isMatch(c, words[1:]) {
			return true
		}
	}

	if matchDigit(current, words) {
		return true
	}

	if matchWord(current, words) {
		return true
	}

	return false
}

func matchDigit(current *templateItem, words []string) bool {
	pd, ok := current.patternChildren[patternKeyDigit]
	if !ok {
		return false
	}

	word := words[0]

	for _, pi := range pd {
		if pi.prefix != "" {
			if word == pi.prefix || !strings.HasPrefix(word, pi.prefix) {
				continue
			}
		}

		if pi.suffix != "" && pi.max > 1 {
			//Agar %d{1,n}suffix ko'rinishida kelsa
			//suffix so'zi ajratilgan bo'lishi kerak
			//va birinchida kelgan so'z prefix%d bo'lishi shart
			if !isDigit.MatchString(word[len(pi.prefix):]) {
				continue
			}

			for i, w := range words[1:] {
				if i >= pi.max {
					break
				}

				if w != pi.suffix {
					//birinchi so'z allaqachon tekshirilgan
					//prefix orqali
					if !isDigit.MatchString(w) {
						break
					}

					continue
				}

				if isMatch(pi.children, words[i+2:]) {
					return true
				}
			}

			continue
		}

		size := len(words)

		if pi.max != 0 {
			if pi.max < size {
				size = pi.max
			}
		}

		for i := 0; i < size; i++ {
			currentWord := words[i]

			if pi.suffix != "" {
				if pi.suffix == currentWord || !strings.HasSuffix(currentWord, pi.suffix) {
					continue
				}

				currentWord = currentWord[:len(currentWord)-len(pi.suffix)]
			}

			if i == 0 && pi.prefix != "" {
				currentWord = currentWord[len(pi.prefix):]
			}

			if !isDigit.MatchString(currentWord) {
				continue
			}

			if isMatch(pi.children, words[i+1:]) {
				return true
			}
		}
	}

	return false
}

func matchWord(current *templateItem, words []string) bool {
	pd, ok := current.patternChildren[patternKeyWord]
	if !ok {
		return false
	}

	word := words[0]

	for _, pi := range pd {
		if pi.prefix != "" {
			if word == pi.prefix || !strings.HasPrefix(word, pi.prefix) {
				continue
			}
		}

		if pi.suffix != "" && pi.max > 1 {

			//Agar %w{1,n}suffix ko'rinishida kelsa
			//suffix so'zi ajratilgan bo'lishi kerak
			//to shu ajratilgan so'zgacha qidirib
			//qolganini oddiy tekshiramiz
			for i, w := range words[1:] {
				if i >= pi.max {
					return false
				}

				if w != pi.suffix {
					continue
				}

				if isMatch(pi.children, words[i+2:]) {
					return true
				}
			}

			continue
		}

		size := len(words)

		if pi.max != 0 {
			if pi.max < size {
				size = pi.max
			}
		}

		for i := 0; i < size; i++ {
			currentWord := words[i]
			if pi.suffix != "" {
				if pi.suffix == currentWord || !strings.HasSuffix(currentWord, pi.suffix) {
					continue
				}
			}

			if isMatch(pi.children, words[i+1:]) {
				return true
			}
		}
	}

	return false
}

func addStaticPattern(current *templateItem, key, n int, isEnd bool) *templateItem {
	//agar prefix va suffix mavjud bo'lmasa
	for _, pi := range current.patternChildren[key] {
		if pi.max == n && pi.prefix == "" && pi.suffix == "" {
			return pi.children
		}
	}

	item := newTemplateItem(isEnd)
	pitem := &patternItem{
		max:      n,
		children: item,
	}

	current.patternChildren[key] = append(current.patternChildren[key], pitem)
	return item
}

func addRegexPattern(current *templateItem, key int, word string, n int, isEnd bool) *templateItem {
	keyword := fmt.Sprintf("%%%s{1,%d}", keyLetter[key], n)
	parts := strings.Split(word, keyword)

	//agar prefix va suffix mavjud bo'lmasa
	for _, pi := range current.patternChildren[key] {
		if pi.max == n && pi.prefix == parts[0] && pi.suffix == parts[1] {
			return pi.children
		}
	}

	item := newTemplateItem(isEnd)
	pitem := &patternItem{
		max:      n,
		prefix:   parts[0],
		suffix:   parts[1],
		children: item,
	}

	current.patternChildren[key] = append(current.patternChildren[key], pitem)
	return item
}

func newTemplateItem(isEnd bool) *templateItem {
	return &templateItem{
		children:        make(map[string]*templateItem),
		patternChildren: make(map[int][]*patternItem),
		isEnd:           isEnd,
	}
}
