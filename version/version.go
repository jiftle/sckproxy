package version

import "strings"

var (
	Name      string = "skpy service"
	Version   string = "v1.0.1"
	BuildTime string
	Author    string
	Hash      string
)

func GetVer() (out string) {
	strTmp := Version[1:]
	arTmp := strings.Split(strTmp, ".")
	nArTmp := []string{arTmp[0], arTmp[1], arTmp[2]}
	out = strings.Join(nArTmp, ".")
	return
}
