package route

type etcdRouteDefinitionReader struct {
}

func (reader *etcdRouteDefinitionReader) GetRouteDefinitions(config string) ([]RouteDefinition, error) {
	routeDefinitions := make([]RouteDefinition, 0)

	return routeDefinitions, nil
}
