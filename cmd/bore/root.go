package main

import (
	ud "github.com/apparentlymart/go-userdirs/userdirs"
)

func execute() error {
	panic("not implemented")
}

func defaultConfigPath() string {
	userdirs := ud.ForApp("bore", "trulyao", "dev.trulyao.bore")
	_ = userdirs

	panic("not implemented: defaultConfigPath")
}
