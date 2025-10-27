package main

import (
	"common/internal/config"
	"common/pkg/sugar"
	"common/scripts"
	"os"
	"path/filepath"
)

func main() {
	workDir := filepath.Join(sugar.Default(os.Getwd()), "common")
	if err := config.LoadEnv(workDir); err != nil {
		panic(err)
	}

	scripts.CreateDefaultUsers()
	// scripts.FillStocks()
	//
	// scripts.MakeGoodOrder()
	//
	// scripts.MakeBadOrderProductNotEnoughStocks()
	// scripts.MakeBadOrderProductDoesNotExist()
	// scripts.MakeBadOrderProductPaymentError()
}