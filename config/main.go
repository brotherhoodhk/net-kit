package config

import "fmt"

func init() {
	fmt.Println("config init")
	Parse["xml"] = xmlparse
	Parse["json"] = jsonparse
	Parse["yaml"] = yamlparse
	Format["xml"] = xmlformat
	Format["json"] = jsonformat
	Format["yaml"] = yamlformat
}
