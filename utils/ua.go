package utils

import (
	"bytes"
	"regexp"
	"strings"
)

type UserAgent struct {
	Name      string
	Version   string
	OS        string
	OSVersion string
	Device    string
	Bot       bool
	URL       string
	String    string
}

var ignore = map[string]struct{}{
	"KHTML, like Gecko": struct{}{},
	"U":                 struct{}{},
	"compatible":        struct{}{},
	"Mozilla":           struct{}{},
	"WOW64":             struct{}{},
}

const (
	Windows      = "Windows"
	WindowsPhone = "Windows Phone"
	Android      = "Android"
	MacOS        = "macOS"
	IOS          = "iOS"
	Linux        = "Linux"

	Opera            = "Opera"
	OperaMini        = "Opera Mini"
	OperaTouch       = "Opera Touch"
	Chrome           = "Chrome"
	Firefox          = "Firefox"
	InternetExplorer = "Internet Explorer"
	Safari           = "Safari"
	Edge             = "Edge"
	Vivaldi          = "Vivaldi"

	Googlebot           = "Googlebot"
	Twitterbot          = "Twitterbot"
	FacebookExternalHit = "facebookexternalhit"
)

func ParseUserAgent(userAgent string) UserAgent {
	ua := UserAgent{
		String: userAgent,
	}

	tokens := parseUserAgent(userAgent)

	for k := range tokens {
		if strings.HasPrefix(k, "http://") || strings.HasPrefix(k, "https://") {
			ua.URL = k
			delete(tokens, k)
			break
		}
	}

	switch {
	case tokens.exists("Android"):
		ua.OS = Android
		ua.OSVersion = tokens[Android]
		for s := range tokens {
			if strings.HasSuffix(s, "Build") {
				ua.Device = strings.TrimSpace(s[:len(s)-5])
			}
		}

	case tokens.exists("iPhone"):
		ua.OS = IOS
		ua.OSVersion = tokens.findMacOSVersion()
		ua.Device = "iPhone"

	case tokens.exists("iPad"):
		ua.OS = IOS
		ua.OSVersion = tokens.findMacOSVersion()
		ua.Device = "iPad"

	case tokens.exists("Windows NT"):
		ua.OS = Windows
		ua.OSVersion = tokens["Windows NT"]

	case tokens.exists("Windows Phone OS"):
		ua.OS = WindowsPhone
		ua.OSVersion = tokens["Windows Phone OS"]

	case tokens.exists("Macintosh"):
		ua.OS = MacOS
		ua.OSVersion = tokens.findMacOSVersion()

	case tokens.exists("Linux"):
		ua.OS = Linux
		ua.OSVersion = tokens[Linux]

	}

	switch {

	case tokens.exists("Googlebot"):
		ua.Name = Googlebot
		ua.Version = tokens[Googlebot]
		ua.Bot = true

	case tokens["Opera Mini"] != "":
		ua.Name = OperaMini
		ua.Version = tokens[OperaMini]

	case tokens["OPR"] != "":
		ua.Name = Opera
		ua.Version = tokens["OPR"]

	case tokens["OPT"] != "":
		ua.Name = OperaTouch
		ua.Version = tokens["OPT"]

	case tokens["OPiOS"] != "":
		ua.Name = Opera
		ua.Version = tokens["OPiOS"]

	case tokens["CriOS"] != "":
		ua.Name = Chrome
		ua.Version = tokens["CriOS"]

	case tokens["FxiOS"] != "":
		ua.Name = Firefox
		ua.Version = tokens["FxiOS"]

	case tokens["Firefox"] != "":
		ua.Name = Firefox
		ua.Version = tokens[Firefox]

	case tokens["Vivaldi"] != "":
		ua.Name = Vivaldi
		ua.Version = tokens[Vivaldi]

	case tokens.exists("MSIE"):
		ua.Name = InternetExplorer
		ua.Version = tokens["MSIE"]

	case tokens["EdgiOS"] != "":
		ua.Name = Edge
		ua.Version = tokens["EdgiOS"]

	case tokens["Edge"] != "":
		ua.Name = Edge
		ua.Version = tokens["Edge"]

	case tokens["Edg"] != "":
		ua.Name = Edge
		ua.Version = tokens["Edg"]

	case tokens["EdgA"] != "":
		ua.Name = Edge
		ua.Version = tokens["EdgA"]

	case tokens["bingbot"] != "":
		ua.Name = "Bingbot"
		ua.Version = tokens["bingbot"]

	case tokens["SamsungBrowser"] != "":
		ua.Name = "Samsung Browser"
		ua.Version = tokens["SamsungBrowser"]

	case tokens.exists(Chrome) && tokens.exists(Safari):
		name := tokens.findBestMatch(true)
		if name != "" {
			ua.Name = name
			ua.Version = tokens[name]
			break
		}
		fallthrough

	case tokens.exists("Chrome"):
		ua.Name = Chrome
		ua.Version = tokens["Chrome"]

	case tokens.exists("Safari"):
		ua.Name = Safari
		if v, ok := tokens["Version"]; ok {
			ua.Version = v
		} else {
			ua.Version = tokens["Safari"]
		}

	default:
		if ua.OS == "Android" && tokens["Version"] != "" {
			ua.Name = "Android browser"
			ua.Version = tokens["Version"]
		} else {
			if name := tokens.findBestMatch(false); name != "" {
				ua.Name = name
				ua.Version = tokens[name]
			} else {
				ua.Name = ua.String
			}
			ua.Bot = strings.Contains(strings.ToLower(ua.Name), "bot")
		}
	}

	if !ua.Bot {
		ua.Bot = ua.URL != ""
	}

	if !ua.Bot {
		switch ua.Name {
		case Twitterbot, FacebookExternalHit:
			ua.Bot = true
		}
	}

	return ua
}

func parseUserAgent(userAgent string) (clients properties) {
	clients = make(map[string]string, 0)
	slash := false
	isURL := false
	var buff, val bytes.Buffer
	addToken := func() {
		if buff.Len() != 0 {
			s := strings.TrimSpace(buff.String())
			if _, ign := ignore[s]; !ign {
				if isURL {
					s = strings.TrimPrefix(s, "+")
				}

				if val.Len() == 0 {
					var ver string
					s, ver = checkVer(s)
					clients[s] = ver
				} else {
					clients[s] = strings.TrimSpace(val.String())
				}
			}
		}
		buff.Reset()
		val.Reset()
		slash = false
		isURL = false
	}

	parOpen := false

	bua := []byte(userAgent)
	for i, c := range bua {

		switch {
		case c == 41:
			addToken()
			parOpen = false

		case parOpen && c == 59:
			addToken()

		case c == 40:
			addToken()
			parOpen = true

		case slash && c == 32:
			addToken()

		case slash:
			val.WriteByte(c)

		case c == 47 && !isURL:
			if i != len(bua)-1 && bua[i+1] == 47 && (bytes.HasSuffix(buff.Bytes(), []byte("http:")) || bytes.HasSuffix(buff.Bytes(), []byte("https:"))) {
				buff.WriteByte(c)
				isURL = true
			} else {
				slash = true
			}

		default:
			buff.WriteByte(c)
		}
	}
	addToken()

	return clients
}

func checkVer(s string) (name, v string) {
	i := strings.LastIndex(s, " ")
	if i == -1 {
		return s, ""
	}

	switch s[:i] {
	case "Linux", "Windows NT", "Windows Phone OS", "MSIE", "Android":
		return s[:i], s[i+1:]
	default:
		return s, ""
	}

}

type properties map[string]string

func (p properties) exists(key string) bool {
	_, ok := p[key]
	return ok
}

func (p properties) existsAny(keys ...string) bool {
	for _, k := range keys {
		if _, ok := p[k]; ok {
			return true
		}
	}
	return false
}

func (p properties) findMacOSVersion() string {
	for k, v := range p {
		if strings.Contains(k, "OS") {
			if ver := findVersion(v); ver != "" {
				return ver
			} else if ver = findVersion(k); ver != "" {
				return ver
			}
		}
	}
	return ""
}

func (p properties) findBestMatch(withVerOnly bool) string {
	n := 2
	if withVerOnly {
		n = 1
	}
	for i := 0; i < n; i++ {
		for k, v := range p {
			switch k {
			case Chrome, Firefox, Safari, "Version", "Mobile", "Mobile Safari", "Mozilla", "AppleWebKit", "Windows NT", "Windows Phone OS", Android, "Macintosh", Linux, "GSA":
			default:
				if i == 0 {
					if v != "" {
						return k
					}
				} else {
					return k
				}
			}
		}
	}
	return ""
}

var rxMacOSVer = regexp.MustCompile("[_\\d\\.]+")

func findVersion(s string) string {
	if ver := rxMacOSVer.FindString(s); ver != "" {
		return strings.Replace(ver, "_", ".", -1)
	}
	return ""
}
