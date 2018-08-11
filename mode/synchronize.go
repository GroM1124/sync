package mode

import (
	"fmt"
	"sort"
	"syncdata/engine"
)

func Sync(c Config) error {
	var ext []string
	if c.Mask.On {
		ext = c.Mask.Ext
	}
	ex1, in1, err := engine.ReadDir(c.Path1, ext)
	if err != nil {
		return err
	}
	ex2, in2, err := engine.ReadDir(c.Path2, ext)
	if err != nil {
		return err
	}
	if c.Mask.Include {
		ex1, in1 = in1, ex1
		ex2, in2 = in2, ex2
	}
	if c.Mask.Verbose {
		excluded := append(ex1, ex2...)
		sort.Slice(excluded, func(i, j int) bool {
			return excluded[i].Abs < excluded[j].Abs
		})
		for _, fi := range excluded {
			fmt.Printf("%q\t%v\t%q\t%q\n", fi.Abs, fi.Size, fi.Time, "the file is excluded by a mask")
		}
	}
	match, dfr := engine.Compare(in1, in2, c.Path1, c.Path2)
	for _, action := range match {
		fmt.Println(action.Description())
	}
	if len(dfr) == 0 {
		fmt.Println("No files for synchronization\n")
		return nil
	}
	for _, action := range dfr {
		fmt.Println(action.Description())
	}
	fmt.Println("Pleace enter \"Y\" for synchronization " +
		"or enter any other character to cancel synchronization\n")
	var ask string
	fmt.Scanln(&ask)
	if ask == "y" || ask == "Y" {
		for _, action := range dfr {
			err := action.Apply()
			if err != nil {
				return err
			}
		}
	} else {
		fmt.Println("Synchronize canceled by user\n")
	}
	return nil
}