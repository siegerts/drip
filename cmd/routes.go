package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var plumberEntryPoint string

func init() {
	rootCmd.AddCommand(routesCmd)
	routesCmd.Flags().StringVarP(&plumberEntryPoint, "entry", "e", "entrypoint.r", "Plumber application entrypoint file")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var routesCmd = &cobra.Command{
	Use:   "routes",
	Short: "Display all routes in your Plumber application",
	Long:  `A quick way to visualize your application's routing structure`,
	Run: func(cmd *cobra.Command, args []string) {
		// gen route structure
		// will need to be recursive to go through file tree
		plumb, _ := regexp.Compile(`(?P<comment>#*).*plumb\("(?P<plumber>[a-zA-Z0-9_]+\.[rR])"\)`)

		routes, _ := regexp.Compile(`#\* @(get|post|put|delete|head)\s/[a-zA-Z0-9\-_\/<>:]+`)
		assets, _ := regexp.Compile(`#\* @assets\s[\.\/a-zA-Z0-9\_]+\s[\.\/a-zA-Z0-9\_]*`)

		// other components
		// programmaticRoutes, _ := regexp.Compile(`#\* @(get|post|put|delete|head)\s/[a-zA-Z0-9\-_\/<>:]+`)
		// mountedRoutes, _ := regexp.Compile(`#\* @(get|post|put|delete|head)\s/[a-zA-Z0-9\-_\/<>:]+`)
		// mountedAssets, _ := regexp.Compile(`#\* @assets\s[\.\/a-zA-Z0-9\_]+\s[\.\/a-zA-Z0-9\_]*`)

		dat, err := ioutil.ReadFile(plumberEntryPoint)
		check(err)

		// if length > 0 then try to read the routes file
		// figure out nests and mounts
		entryMatches := plumb.FindAllStringSubmatch(string(dat), 1)
		// index remains the same if no match
		comment := entryMatches[0][1]
		if len(entryMatches) > 0 && comment == "" {

			dat, err := ioutil.ReadFile(entryMatches[0][2])
			check(err)

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
			// route table
			// refactor into function
			fmt.Println("Routes")
			routeMatches := routes.FindAllStringSubmatch(string(dat), -1)
			for _, match := range routeMatches {
				s := strings.TrimPrefix(match[0], "#*")
				parts := strings.Split(s, " ")
				fmt.Fprintln(w, strings.Join(parts, "\t"))
			}
			w.Flush()

			fmt.Println("Static Assets")
			// static asset table
			assetMatches := assets.FindAllStringSubmatch(string(dat), -1)
			for _, match := range assetMatches {
				s := strings.TrimPrefix(match[0], "#*")
				parts := strings.Split(s, " ")
				fmt.Fprintln(w, strings.Join(parts, "\t"))
			}
			w.Flush()

			// need to deal with mounting and static file routers

			// pr <- plumber$new()
			// pr$handle("GET", "/", function(req, res){
			//   # ...
			// })

			// pr$handle("POST", "/submit", function(req, res){
			//   # ...
			// })

			// 			root <- plumber$new()

			// users <- plumber$new("users.R")
			// root$mount("/users", users)

			// products <- plumber$new("products.R")
			// root$mount("/products", products)

			// 			pr <- plumber$new()

			// stat <- PlumberStatic$new("./myfiles")

			// pr$mount("/assets", stat)

		}
	},
}
