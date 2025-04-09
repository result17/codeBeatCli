package vipertools

import (
	"strings"

	"github.com/spf13/viper"
)

func GetString(v *viper.Viper, key string) string {
	return strings.Trim(v.GetString(key), `"'`)
}
