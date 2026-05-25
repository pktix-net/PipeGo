package config

import (
	"bufio"
	"fmt"
	"os"
	"pipego/internal/route"
	"strings"
)

func LoadRoutes(path string) ([]route.Route, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var routes []route.Route

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}

		// Skip comments
		if strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)

		// Expected:
		// tcp 0.0.0.0:80 -> 192.168.8.8:80
		if len(parts) != 4 {
			return nil, fmt.Errorf("invalid route format: %s", line)
		}

		if parts[2] != "->" {
			return nil, fmt.Errorf("missing -> in route: %s", line)
		}

		r := route.Route{
			Protocol: parts[0],
			Listen:   parts[1],
			Upstream: parts[3],
		}

		routes = append(routes, r)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return routes, nil
}
