package config

// Validate the config
func ValidateConfig(conf *Config) (map[string]string, map[string]string, error) {
	return nil, nil, nil
	// // Usermap
	// userMap := map[string]string{}
	// for name, email := range conf.Users {
	// 	if _, ok := userMap[name]; !ok {
	// 		if strings.Contains(name, ".") {
	// 			return nil, nil, fmt.Errorf("ERROR: name in users can not contain a .")
	// 		}
	// 		userMap[name] = user.Email
	// 	} else {
	// 		return nil, nil, fmt.Errorf("ERROR: No duplicate user names")
	// 	}
	// }

	// // Sourcesmap
	// sourcesMap := map[string]string{}
	// for _, source := range config.Sources {
	// 	if _, ok := userMap[source.Name]; !ok {
	// 		if strings.Contains(source.Name, ".") {
	// 			return nil, nil, fmt.Errorf("ERROR: name in sources can not contain a .")
	// 		}
	// 		sourcesMap[source.Name] = source.Url
	// 	} else {
	// 		return nil, nil, fmt.Errorf("ERROR: No duplicate user names")
	// 	}
	// }

	// return userMap, sourcesMap, nil
}
