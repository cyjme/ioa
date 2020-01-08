package route

type RouteDefinition struct {
	Id         string   `mapstructure:"id"`
	Uri        string   `mapstructure:"uri"`
	Filters    []string `mapstructure:"filters"`
	Predicates []string `mapstructure:"predicates"`
}

type routeDefinitionReader interface {
	GetRouteDefinitions(config string) ([]RouteDefinition, error)
}

type RouteDefinitionWriter interface {
	AddRouteDefinition(config string, routDefinition RouteDefinition) error
	UpdateRouteDefinition(config string, routDefinition RouteDefinition) error
	DeleteRouteDefinition(config string, routDefinitionIds string) error
}
