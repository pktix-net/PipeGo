package config

import (
	"bufio"
	"fmt"
	"os"
	"pipego/internal/route"
	"strconv"
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

		// Split by space but preserve key=value in parts[3]+
		parts := strings.Fields(line)

		// Expected:
		// tcp 0.0.0.0:80 -> 192.168.8.8:80 [allow_asn=4134,4837] [deny_asn=13335]
		if len(parts) < 4 {
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

		// Parse optional ASN filters from remaining parts
		for i := 4; i < len(parts); i++ {
			kv := strings.SplitN(parts[i], "=", 2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("invalid filter format: %s in %s", parts[i], line)
			}
			key := kv[0]
			raw := kv[1]

			asns, err := parseIntList(raw)
			if err != nil {
				return nil, fmt.Errorf("invalid asn list in %s: %w", line, err)
			}

			switch key {
			case "allow_asn":
				r.AllowASN = asns
			case "deny_asn":
				r.DenyASN = asns
			default:
				return nil, fmt.Errorf("unknown filter: %s in %s", key, line)
			}
		}

		routes = append(routes, r)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return routes, nil
}

// parseIntList parses a comma-separated list of integers e.g. "4134,4837,9808"
func parseIntList(s string) ([]int, error) {
	if s == "" {
		return nil, nil
	}
	items := strings.Split(s, ",")
	out := make([]int, 0, len(items))
	for _, item := range items {
		n, err := strconv.Atoi(strings.TrimSpace(item))
		if err != nil {
			return nil, fmt.Errorf("invalid number: %s", item)
		}
		out = append(out, n)
	}
	return out, nil
}
