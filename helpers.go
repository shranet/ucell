package ucell_tpl_match

import (
	"fmt"
	"html"
	"regexp"
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
