package app

import "net/http"

func Main(configPath string, listenOn string) error {
	config, err := LoadConfigurationFromFile(configPath)
	if err != nil {
		return err
	}
	routes := createAllRoutes(config)

	for path, route := range routes {
		http.HandleFunc(path, route)
	}

	return http.ListenAndServe(listenOn, nil)
}
