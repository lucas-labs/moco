package console

// draws ascii header

import (
	"fmt"
	"moco/config"
	c "moco/utils/console/color"
)

func Header(config *config.Config) {
	var lines []string

	if config != nil {
		lines = []string{
			c.BrightGreenBold + " _ _  ___  ___  ___    " + c.Reset,
			c.BrightGreenBold + "| | )|   )|    |   )   " + c.Reset + "Serving: " + config.Name + c.Blue + "@" + c.Reset + config.Version,
			c.BrightGreenBold + "|  / |__/ |__  |__/    " + c.Reset + config.Description + c.Blue + " | " + c.Reset + fmt.Sprintf("http://127.0.0.1:%d", config.Port),
			"",
		}
	} else {
		lines = []string{
			"\033[1m\033[92m  _ _  ___  ___  ___ \033[0m",
			"\033[1m\033[92m | | )|   )|    |   )\033[0m",
			"\033[1m\033[92m |  / |__/ |__  |__/ \033[0m\n",
		}
	}

	for _, line := range lines {
		fmt.Println(line)
	}
}
