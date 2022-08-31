package app

import "net/http"

func Main(configPath string) error {
	config, err := LoadConfigurationFromFile(configPath)
	if err != nil {
		return err
	}
	routes := createAllRoutes(config)

	for path, route := range routes {
		http.HandleFunc(path, route)
	}

	return http.ListenAndServe(":8090", nil)
}
