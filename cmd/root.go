// Package cmd is a root command.
/*
Copyright © 2023 Takafumi Miyanaga <miya.org.0309@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/TheZoraiz/ascii-image-converter/aic_package"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type gitHubUser struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	Repos     int    `json:"public_repos"`
	Followers int    `json:"followers"`
	Following int    `json:"following"`
}

var username string
var highlightColor string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ghfetch",
	Short: "Fetch GitHub user's profile",
	Run: func(cmd *cobra.Command, args []string) {
		if username == "" {
			fmt.Println("Please provide a GitHub username using the --user flag.")
			return
		}

		flags := aic_package.DefaultFlags()
		flags.Dimensions = []int{50, 25}
		flags.Colored = true
		flags.CustomMap = " .-=+#@"
		asciiArt, err := aic_package.Convert(fmt.Sprintf("https://github.com/%s.png", username), flags)
		if err != nil {
			fmt.Println(err)
			return
		}
		asciiLines := strings.Split(asciiArt, "\n")

		user, err := fetchUser(username)
		if err != nil {
			fmt.Println("Error fetching user information:", err)
			return
		}
		titleColor := colorMap[highlightColor].SprintFunc()
		infoColor := color.New(color.FgWhite).SprintFunc()

		// Displaying username at the top
		fmt.Printf("%s    %s: %s\n", asciiLines[0], titleColor("GitHub User"), titleColor(user.Login))

		userInfo := []string{
			fmt.Sprintf("Name: %s", user.Name),
			fmt.Sprintf("Repos: %d", user.Repos),
			fmt.Sprintf("Followers: %d", user.Followers),
			fmt.Sprintf("Following: %d", user.Following),
		}

		for i := 1; i < len(asciiLines) || i-1 < len(userInfo); i++ {
			left := ""
			if i < len(asciiLines) {
				left = asciiLines[i]
			}
			right := ""
			if i-1 < len(userInfo) {
				splitted := strings.Split(userInfo[i-1], ": ")
				right = titleColor(splitted[0]) + ": " + infoColor(splitted[1])
			}
			fmt.Printf("%-60s    %s\n", left, right)
		}
	},
}

func fetchUser(username string) (*gitHubUser, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s", username)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user gitHubUser
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVarP(&username, "user", "u", "", "GitHub username")
	rootCmd.Flags().StringVarP(&highlightColor, "color", "c", "blue", "Highlight color for text")
}

// Color map for user-specified colors
var colorMap = map[string]*color.Color{
	"red":     color.New(color.FgRed),
	"green":   color.New(color.FgGreen),
	"yellow":  color.New(color.FgYellow),
	"blue":    color.New(color.FgBlue),
	"magenta": color.New(color.FgMagenta),
	"cyan":    color.New(color.FgCyan),
}
