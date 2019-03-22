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
	Port     string
	Db       string
	Query    map[string]string
}

// TODO https://github.com/go-sql-driver/mysql/blob/dc029498cb5a3efbe44e54dcb5cf080d451450fa/utils.go#L81
// remove regexp
func GetDsnMatches(url string) [][]string {
	comp := `^(?P<proto>.*?):\/\/(?:(?P<user>.*?)(?::(?P<passwd>.*))?@)?(?:(?P<net>[^\(\:]*)\:?(?P<port>[0-9]{0,5}))?\/(?P<dbname>.*?)(?:\?(?P<params>[^\?]*))?$`
	r, err := regexp.Compile(comp)
	matches := r.FindAllStringSubmatch(url, -1)

	if err != nil {
		log.Panic("Could not parse: ", err)
	}

	return matches
}

func MakeDsnParts(url string) DsnParts {
	matches := GetDsnMatches(url)

	proto := matches[0][1]
	username := matches[0][2]
	password := matches[0][3]
	host := matches[0][4]
	port := matches[0][5]
	db := matches[0][6]

	parts := DsnParts{
		Proto:    proto,
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
		Db:       db,
	}

	return parts
}
