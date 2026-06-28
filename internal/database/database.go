package database

import (
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"

	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
)

func Init(dbPath string) error {
	link := fmt.Sprintf("sqlite::@file(%s)", dbPath)
	if err := gdb.SetConfig(gdb.Config{
		"default": gdb.ConfigGroup{
			{
				Type:   "sqlite",
				Link:   link,
				Prefix: "ydsterm_",
				Debug:  true,
			},
		},
	}); err != nil {
		return fmt.Errorf("set config: %w", err)
	}
	_, err := gdb.Instance("default")
	return err
}
