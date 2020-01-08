package path

import (
	"ioa/context"
	"regexp"
	"strings"
)

type predicate struct {
	config config
}

type route struct {
	pattern             string
	patternRegexpString string
	patternRegexp       *regexp.Regexp
	params              map[int]string
	paramValues         []string
}

type config struct {
	route route
}

func New(arg string) (predicate, error) {
	p := predicate{}
	route, err := buildRoute(arg)
	if err != nil {
		return p, err
	}
	p.config.route = route

	return p, nil
}

func (p *predicate) Name() string {
	return "path"
}

func (p *predicate) Apply(ctx *context.Context) bool {

	matched := p.matchRoute(ctx.Request.URL.Path)
	if !matched {
		return false
	}

	ctx.PathParam = make(map[string]string)

	for i, param := range p.config.route.params {
		ctx.PathParam[param] = p.config.route.paramValues[i]
	}

	return matched
}

func buildRoute(pattern string) (route, error) {
	route := route{}
	route.pattern = pattern
	parts := strings.Split(pattern, "/")

	j := 0

	params := make(map[int]string)
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			expr := "([^/]+)"

			//a user may choose to override the defult expression
			// similar to expressjs: ‘/user/:id([0-9]+)’
			if index := strings.Index(part, "("); index != -1 {
				expr = part[index:]
				part = part[:index]
			}
			params[j] = part
			parts[i] = expr
			j++
		}
	}

	//recreate the url pattern, with parameters replaced
	//by regular expressions. then compile the regex
	pattern = strings.Join(parts, "/")

	route.patternRegexpString = pattern
	regex, regexErr := regexp.Compile(pattern)

	if regexErr != nil {
		return route, regexErr
	}

	//now create the Route
	route.patternRegexp = regex
	route.params = params

	return route, nil
}

func (p *predicate) matchRoute(path string) bool {
	if !p.config.route.patternRegexp.MatchString(path) {
		p.config.route.paramValues = []string{}
		return false
	}

	matches := p.config.route.patternRegexp.FindStringSubmatch(path)
	p.config.route.paramValues = matches[1:]

	return len(matches[0]) == len(path)
}
