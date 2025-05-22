package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/go-caching-proxy/internal/cache"
	"github.com/user/go-caching-proxy/internal/proxy"
)

func main() {
	var port int
	var origin string
	var clearCache bool

	rootCmd := &cobra.Command{
		Use:   "caching-proxy",
		Short: "A caching proxy server that forwards requests and caches responses",
		Long: `A caching proxy server that forwards requests to the specified origin server 
		and caches the responses. If the same request is made again, it returns the cached response.`,
		Run: func(cmd *cobra.Command, args []string) {
			if clearCache {
				cache.Clear()
				fmt.Println("Cache cleared successfully")
				return
			}

			if port == 0 || origin == "" {
				cmd.Help()
				return
			}

			fmt.Printf("Starting caching proxy server on port %d, forwarding to %s\n", port, origin)
			proxy.Start(port, origin)
		},
	}

	rootCmd.Flags().IntVar(&port, "port", 0, "Port on which the caching proxy server will run")
	rootCmd.Flags().StringVar(&origin, "origin", "", "URL of the server to which the requests will be forwarded")
	rootCmd.Flags().BoolVar(&clearCache, "clear-cache", false, "Clear the cache")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
