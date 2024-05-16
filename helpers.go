package ucell_tpl_match

import (
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
// patternW = "[a-zа-я0-9]+"
// patternD = "[0-9]+"
)

var whiteSpaceRegex = regexp.MustCompile(`\p{Z}+`)
var nonLatinCyrillicRegex = regexp.MustCompile(`[^a-zA-Zа-яА-Я0-9 ]+`)
var spaceReplacer = strings.NewReplacer(
	"\n", " ",
	"\r", " ",
)

var regexReplaces = strings.NewReplacer(
	"%d{1,1}", "[0-9]+",
	"%d{1,0}", "[0-9]+( [0-9]+)*",
	"%w{1,1}", "[a-zа-я0-9]+",
	"%w{1,0}", "[a-zа-я0-9]+( [a-zа-я0-9]+)*",
)

var withSuffix = regexp.MustCompile("%[dw]{1,([0-9]+)}[^ ]")
var others = regexp.MustCompile("%[dw]{1,([0-9]+)}")

func cleanMessage(message string) string {
	message = html.UnescapeString(message)
	message = spaceReplacer.Replace(message)
	message = nonLatinCyrillicRegex.ReplaceAllString(message, "")

	return strings.Trim(whiteSpaceRegex.ReplaceAllString(message, " "), " ")
}

func cleanTemplate(tpl string) string {
	prefix := fmt.Sprintf("ucelltpl%d", time.Now().UnixMilli())

	tplKey := func(key string) string {
		return prefix + key
	}

	oldnew := []string{}
	newold := []string{}

	addOldNew := func(old, new, pattern string) {
		oldnew = append(oldnew, old)
		oldnew = append(oldnew, new)

		newold = append(newold, new)
		newold = append(newold, pattern)
	}

	for i := 100; i >= 2; i-- {
		addOldNew(fmt.Sprintf("%%w{1,%d}", i), tplKey(fmt.Sprintf("range%dw", i)),
			fmt.Sprintf("%%w{1,%d}", i))
		addOldNew(fmt.Sprintf("%%d{1,%d}", i), tplKey(fmt.Sprintf("range%dd", i)),
			fmt.Sprintf("%%d{1,%d}", i))
	}

	addOldNew("%w+", tplKey("plusw"), "%w{1,0}")
	addOldNew("%d+", tplKey("plusd"), "%d{1,0}")
	addOldNew("%w", tplKey("w"), "%w{1,1}")
	addOldNew("%d", tplKey("d"), "%d{1,1}")

	cleanTpl := cleanMessage(strings.NewReplacer(oldnew...).Replace(tpl))
	return strings.NewReplacer(newold...).Replace(regexp.QuoteMeta(cleanTpl))
}

func regexReplaceFunction(addSpace bool) func(v string) string {
	return func(v string) string {
		n := strings.Split(v, ",")[1]
		suffix := ""
		if addSpace {
			suffix = " " + string(v[len(v)-1])
			n = n[:len(n)-2]
		} else {
			n = n[:len(n)-1]
		}

		nn, _ := strconv.Atoi(n)

		if v[1] == 'w' {
			return fmt.Sprintf("[a-zа-я0-9]+( [a-zа-я0-9]+){1,%d}", nn-1) + suffix
		}

		return fmt.Sprintf("[0-9]+( [0-9]+){1,%d}", nn-1) + suffix
	}
}

func CreateRegexp(tpl string) *regexp.Regexp {
	tpl = regexReplaces.Replace(cleanTemplate(tpl))

	tpl = withSuffix.ReplaceAllStringFunc(tpl, regexReplaceFunction(true))
	tpl = others.ReplaceAllStringFunc(tpl, regexReplaceFunction(false))

	return regexp.MustCompile("(?i)^" + tpl + "$")
}

func IsMatch(r *regexp.Regexp, message string) bool {
	return r.MatchString(cleanMessage(message))
}
