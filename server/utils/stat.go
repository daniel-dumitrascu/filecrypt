package utils

import "os"

func IsDir(path string) (bool, error) {
	log := GetLogger()
	info, err := os.Stat(path)
	if err != nil {
		log.Error("Error getting the info stats for (", path, "): ", err)
		return false, err
	}

	return info.IsDir(), nil
}
