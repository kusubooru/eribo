// +build !prod

package mysql

func legacyDefault() string {
	return "DEFAULT CURRENT_TIMESTAMP"
}
