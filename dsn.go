package turn

import (
	"log"
	"regexp"
)

type DsnParts struct {
	Proto    string
	Username string
	Password string
	Host     string
	Port     int
	Db       string
	Query    map[string]string
}

// TODO https://github.com/go-sql-driver/mysql/blob/dc029498cb5a3efbe44e54dcb5cf080d451450fa/utils.go#L81
// remove regexp
func GetDnsMatches(url string) [][]string {
	comp := `^(?P<proto>.*?):\/\/(?:(?P<user>.*?)(?::(?P<passwd>.*))?@)?(?:(?P<net>[^\(\:]*)\:?(?P<port>[0-9]{0,5}))?\/(?P<dbname>.*?)(?:\?(?P<params>[^\?]*))?$`
	r, err := regexp.Compile(comp)
	matches := r.FindAllStringSubmatch(url, -1)

	if err != nil {
		log.Panic("Could not parse: ", err)
	}

	return matches
}
