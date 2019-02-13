package godatparseURL

import (
	"net/url"
	"regexp"
)

var reSchemeRegex regexp.Regexp = *regexp.MustCompile("(?i)[a-z]+:\\/\\/")
var reVersionRegex = *regexp.MustCompile("(?i)^(dat:\\/\\/)?([^/]+)(\\+[^/]+)(.*)$")

type ParsedDatURL struct {
	Protocol string
	Slashes  bool
	Auth     *URLAuth
	Host     string
	Port     *string
	Hostname string
	Hash     *string
	Search   string
	Query    string
	Pathname string
	Path     string
	Href     string
	Version  string
}

type URLAuth struct {
	User       string
	Pass       string
	PaswordSet bool
}

func parseDatURL(str string) (*ParsedDatURL, error) {
	var err error
	if !reSchemeRegex.MatchString(str) {
		str = "dat://" + str
	}

	var parsed ParsedDatURL
	var version string
	var match [][]string
	var goURL *url.URL
	match = reVersionRegex.FindAllStringSubmatch(str, -1)
	if len(match) != 0 {
		// We are dealing with a dat url with a version.
		version = match[0][3][1:] // TODO: Am I properly slicing the string here?
		var recombinedString string
		recombinedString = match[0][1] + match[0][2] + match[0][4]
		goURL, err = url.Parse(recombinedString)

	} else {
		goURL, err = url.Parse(str)
		version = ""
	}
	if err != nil {
		return nil, err
	}
	var userInfo URLAuth
	if goURL.User != nil {
		var tempPass string
		var tempBool bool
		tempPass, tempBool = goURL.User.Password()
		userInfo = URLAuth{
			User:       goURL.User.Username(),
			Pass:       tempPass,
			PaswordSet: tempBool,
		}
	}

	parsed = ParsedDatURL{
		Protocol: goURL.Scheme,
		Slashes:  true,
		Auth:     &userInfo,
		Host:     goURL.Host,
		Hostname: goURL.Host,
		Query:    goURL.RawQuery,
		Pathname: goURL.Path,
		Href:     str,
		Version:  version,
	}
	return &parsed, nil
}
