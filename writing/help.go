package writing

import (
	"fmt"
)

func WriteHelp() {
	var d *Dialog = GlobalDialog
	fmt.Println(d.Help.HelpMenu)
}