package common

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var promptColor = color.New(color.FgYellow, color.Bold)

// SetConfigVal sets value to key if value exists
// If value doesn't exist, then the value is read from stdin
func SetConfigVal(key, value, usage, prompt string) error {
	// check if value exists
	if value != "" {
		viper.Set(key, value)
		return nil
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println(usage)
	fmt.Println()
	promptColor.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	fmt.Println()

	// set config
	viper.Set(key, strings.TrimSpace(input))
	return nil
}
