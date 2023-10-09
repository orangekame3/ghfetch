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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/TheZoraiz/ascii-image-converter/aic_package"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/nsf/termbox-go"
	"github.com/spf13/cobra"
)

type gitHubUser struct {
	Login                string `json:"login"`
	Name                 string `json:"name"`
	Repos                int    `json:"public_repos"`
	Followers            int    `json:"followers"`
	Following            int    `json:"following"`
	TotalStarsEarned     int    `json:"total_star_earned"`
	TotalCommitsThisYear int    `json:"total_commit_this_year"`
	TotalPRs             int    `json:"total_pr"`
	TotalIssues          int    `json:"total_issues"`
}

var (
	username       string
	highlightColor string
	profileData    map[string]interface{}
	profileString  string
	accessToken    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ghfetch",
	Short: "Fetch GitHub user's profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		if accessToken == "" {
			return errors.New("ACCESS_TOKEN environment variable is not set!")
		}
		paneWidth := getTermWidth() / 2
		defaultWidth := 50
		defaultHeight := 25
		scaleFactor := float64(paneWidth) / float64(defaultWidth)
		if scaleFactor > 1 {
			scaleFactor = 1
		}
		newWidth := int(float64(defaultWidth) * scaleFactor)

		if username == "" {
			return errors.New("Please provide a GitHub username using the --user flag.")
		}
		flags := aic_package.DefaultFlags()
		flags.Dimensions = []int{newWidth, int(float64(defaultHeight) * scaleFactor)}
		flags.Colored = true
		flags.CustomMap = " .-=+#@"
		asciiArt, err := aic_package.Convert(fmt.Sprintf("https://github.com/%s.png", username), flags)
		if err != nil {
			return err
		}
		leftPane := lipgloss.NewStyle().Width(newWidth).Render(asciiArt)

		user, err := fetchUserWithGraphQL(username, accessToken)

		if err != nil {
			return errors.New("Error fetching user information")
		}

		// Create user info pane
		titleColor := colorMap[highlightColor].SprintFunc()
		User := titleColor("User")
		Name := titleColor("Name")
		Repos := titleColor("Repos")
		Followers := titleColor("Followers")
		Following := titleColor("Following")
		TotalStarsEarned := titleColor("Total Stars Earnd")
		TotalCommitsThisYear := titleColor("Total Commit This Year")
		TotalPRs := titleColor("Total PRs")
		TotalIssues := titleColor("Total Issues")

		userInfoPane := []string{

			fmt.Sprintf("  %s: %s", User, username),
			separator(username),
			fmt.Sprintf("  %s: %s", Name, user.Name),
			fmt.Sprintf("  %s: %d", Repos, user.Repos),
			fmt.Sprintf("  %s: %d", Followers, user.Followers),
			fmt.Sprintf("  %s: %d", Following, user.Following),
			fmt.Sprintf("  %s: %d", TotalStarsEarned, user.TotalStarsEarned),
			fmt.Sprintf("  %s: %d", TotalCommitsThisYear, user.TotalCommitsThisYear),
			fmt.Sprintf("  %s: %d", TotalPRs, user.TotalPRs),
			fmt.Sprintf("  %s: %d", TotalIssues, user.TotalIssues),
		}
		if profileString != "" {
			err := json.Unmarshal([]byte(profileString), &profileData)
			if err != nil {
				fmt.Println("error parsing profile json:", err)
				return err
			}

			// Append custom profile data to the userInfoPane
			for key, value := range profileData {
				key := titleColor(key)
				userInfoPane = append(userInfoPane, fmt.Sprintf("  %s: %v", key, value))
			}
		}
		rightPaneContent := strings.Join(userInfoPane, "\n")
		rightPane := lipgloss.NewStyle().Width(paneWidth).Render(rightPaneContent)

		fmt.Print(lipgloss.JoinHorizontal(lipgloss.Center, leftPane, rightPane))
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVarP(&username, "user", "u", "", "GitHub username")
	rootCmd.Flags().StringVarP(&highlightColor, "color", "c", "blue", "Highlight color for text")
	rootCmd.Flags().StringVarP(&profileString, "profile", "p", "", "Additional user profile in JSON format")
	rootCmd.PersistentFlags().StringVar(&accessToken, "access-token", "", "Your GitHub access token")
}

var colorMap = map[string]*color.Color{
	"red":     color.New(color.FgRed),
	"green":   color.New(color.FgGreen),
	"yellow":  color.New(color.FgYellow),
	"blue":    color.New(color.FgBlue),
	"magenta": color.New(color.FgMagenta),
	"cyan":    color.New(color.FgCyan),
}

func getTermWidth() int {
	err := termbox.Init()
	if err != nil {
		log.Fatalf("error initializing termbox: %v", err)
	}
	defer termbox.Close()

	width, _ := termbox.Size()
	return width
}

func separator(value string) string {
	titleColor := colorMap[highlightColor].SprintFunc()
	lineLength := 6 + len(value)
	return titleColor("  " + strings.Repeat("-", lineLength))
}

func fetchUserWithGraphQL(username, accessToken string) (*gitHubUser, error) {
	url := "https://api.github.com/graphql"
	token := accessToken // 実際には環境変数などから取得する
	var endCursor string
	var hasNextPage bool = true
	var totalStars int
	var totalRepos int
	var query string

	for hasNextPage {
		if endCursor == "" {
			query = fmt.Sprintf(`{
		user(login: "%s") {
			name
			repositories(first: 100, ownerAffiliations: OWNER, isFork: false, orderBy: {direction: DESC, field: STARGAZERS}) {
				nodes {
					stargazers {
						totalCount
					}
				}
				pageInfo {
					endCursor
					hasNextPage
				}
			}
			followers {
				totalCount
			}
			following {
				totalCount
			}
			starredRepositories {
				totalCount
			}
			contributionsCollection {
				totalCommitContributions
			}
			pullRequests {
				totalCount
			}
			issues {
				totalCount
			}
		}
	}`, username)
		} else {
			query = fmt.Sprintf(`{
		user(login: "%s") {
			name
			repositories(first: 100, ownerAffiliations: OWNER, isFork: false, orderBy: {direction: DESC, field: STARGAZERS}, after: "%s") {
				nodes {
					stargazers {
						totalCount
					}
				}
				pageInfo {
					endCursor
					hasNextPage
				}
			}
			followers {
				totalCount
			}
			following {
				totalCount
			}
			starredRepositories {
				totalCount
			}
			contributionsCollection {
				totalCommitContributions
			}
			pullRequests {
				totalCount
			}
			issues {
				totalCount
			}
		}
	}`, username, endCursor)
		}

		body, err := json.Marshal(map[string]string{
			"query": query,
		})
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return nil, err
		}
		resp.Body.Close()
		for _, repo := range result.Data.User.Repositories.Nodes {
			totalStars += repo.Stargazers.TotalCount
		}
		totalRepos += len(result.Data.User.Repositories.Nodes)
		endCursor = result.Data.User.Repositories.PageInfo.EndCursor
		hasNextPage = result.Data.User.Repositories.PageInfo.HasNextPage
	}
	user := &gitHubUser{
		Name:                 result.Data.User.Name,
		Repos:                totalRepos,
		Followers:            result.Data.User.Followers.TotalCount,
		Following:            result.Data.User.Following.TotalCount,
		TotalStarsEarned:     totalStars,
		TotalCommitsThisYear: result.Data.User.ContributionsCollection.TotalCommitContributions,
		TotalPRs:             result.Data.User.PullRequests.TotalCount,
		TotalIssues:          result.Data.User.Issues.TotalCount,
	}

	return user, nil
}

var result struct {
	Data struct {
		User struct {
			Name      string `json:"name"`
			Followers struct {
				TotalCount int `json:"totalCount"`
			}
			Following struct {
				TotalCount int `json:"totalCount"`
			}
			Repositories struct {
				Nodes []struct {
					Stargazers struct {
						TotalCount int `json:"totalCount"`
					} `json:"stargazers"`
				} `json:"nodes"`
				PageInfo struct {
					EndCursor   string `json:"endCursor"`
					HasNextPage bool   `json:"hasNextPage"`
				} `json:"pageInfo"`
			} `json:"repositories"`
			PullRequests struct {
				TotalCount int `json:"totalCount"`
			}
			Issues struct {
				TotalCount int `json:"totalCount"`
			}
			ContributionsCollection struct {
				TotalCommitContributions int `json:"totalCommitContributions"`
			}
		} `json:"user"`
	} `json:"data"`
}
