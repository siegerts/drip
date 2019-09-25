package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(routesCmd)
	routesCmd.Flags().StringVarP(&entryPoint, "entry", "e", "entrypoint.r", "Plumber application entrypoint file")
}

var routesCmd = &cobra.Command{
	Use:   "routes",
	Short: "Display routes in your Plumber application",
	Long:  `A quick way to visualize your Plumber application's routing structure`,
	Run: func(cmd *cobra.Command, args []string) {
		app := Application{
			entryPoint: entryPoint,
		}
		app.RouteStructure()
	},
}

// RouteStructure outputs the parsed endpoints for a given entrypoint file
// @TODO: need to deal with mounting and static file routers
// gen route structure, maybe write a lexer in the future
func (app *Application) RouteStructure() {

	plumberFile, _ := regexp.Compile(`(?i)(?P<comment>#*).*plumb\("(?P<plumber>[a-zA-Z0-9_]+\.[rR])"\)`)
	routes, _ := regexp.Compile(`(?i)#\*\s*@(get|post|put|delete|head)\s/[a-zA-Z0-9\-_\/<>:]+`)
	assets, _ := regexp.Compile(`(?i)#\*\s*@assets\s*[\.\/a-zA-Z0-9\_]+\s[\.\/a-zA-Z0-9\_]*`)

	// other components
	programmaticRoutes, _ := regexp.Compile(`(?i)\$handle\(\"(get|post|put|delete|head)\",\s*\"\/(?P<route>[a-zA-Z0-9_]+)\"`)

	dat, err := ioutil.ReadFile(app.entryPoint)
	if err != nil {
		fmt.Println("Exiting... Error reading entrypoint file: ", err)
		os.Exit(1)
	}

	entryMatches := plumberFile.FindAllStringSubmatch(string(dat), -1)

	if len(entryMatches) > 0 {
		for _, entry := range entryMatches {

			comment := entry[1]
			if comment != "#" {

				dat, err := ioutil.ReadFile(entry[2])
				if err != nil {
					fmt.Println("Exiting... Error reading plumber file: ", err)
					os.Exit(1)
				}

				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Plumber Verb", "Endpoint", "Handler"})
				data := [][]string{}

				// route table
				// refactor into function
				routeMatches := routes.FindAllStringSubmatch(string(dat), -1)
				for _, match := range routeMatches {
					s := strings.TrimPrefix(match[0], "#*")
					parts := strings.Split(s, " ")

					// route filter
					var printRoute = true
					if routeFilter != "" && !strings.Contains(parts[2], routeFilter) {
						printRoute = false
					}

					// flag for absolute endpoint
					// needs refactored into function
					if printRoute {
						if absoluteHost {
							var endpoint string
							if app.host != "" {

								endpoint = strings.TrimRight(app.host, "/") + ":" + strconv.Itoa(app.port) + parts[2]
							} else {
								endpoint = parts[2]
							}

							data = append(data, []string{parts[1], endpoint, "function"})

						} else {
							data = append(data, []string{parts[1], parts[2], "function"})
						}
					}

				}

				// programmatic routes
				programmaticRouteMatches := programmaticRoutes.FindAllStringSubmatch(string(dat), -1)
				for _, match := range programmaticRouteMatches {
					s := strings.TrimPrefix(match[0], "$handle(")
					parts := strings.Split(strings.Replace(s, "\"", "", -1), ",")
					if absoluteHost {
						var endpoint string
						if app.host != "" {
							endpoint = strings.TrimRight(app.host, "/") + ":" + strconv.Itoa(app.port) + parts[2]
						} else {
							endpoint = parts[2]
						}

						data = append(data, []string{parts[1], endpoint, "function"})

					} else {
						data = append(data, []string{parts[1], parts[2], "function"})
					}

				}

				// static asset table
				assetMatches := assets.FindAllStringSubmatch(string(dat), -1)
				for _, match := range assetMatches {
					s := strings.TrimPrefix(match[0], "#*")
					parts := strings.Split(s, " ")

					if absoluteHost {
						var endpoint string
						if app.host != "" {
							endpoint = strings.TrimRight(app.host, "/") + ":" + strconv.Itoa(app.port) + strings.TrimLeft(parts[2], ".")
						} else {
							endpoint = parts[2]
						}

						data = append(data, []string{parts[1], endpoint, "static assets"})

					} else {
						data = append(data, []string{parts[1], parts[2], "static assets"})
					}
				}
				for _, v := range data {
					table.Append(v)
				}
				fmt.Println()
				table.Render()
				fmt.Println()
			}
		}
	}

}
