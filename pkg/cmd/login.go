package cmd

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func runLoginCmd(cmd *cobra.Command, args []string) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		url := launcher.New().Headless(false).MustLaunch()
		browser := rod.New().ControlURL(url).MustConnect()
		page := browser.MustPage("https://app.stainlessapi.com/api/auth/signin/github")

		wait := page.EachEvent(func(e *proto.NetworkRequestWillBeSentExtraInfo) {
			cookieJson := e.Headers["cookie"]
			cookies, err := cookieJson.MarshalJSON()
			if err != nil {
				return
			}
			cookieStr := string(cookies)

			if strings.Contains(cookieStr, "__Secure-next-auth.session-token") {
				fmt.Println("✅ Successfully authenticated with Stainless API")
				cookieStr = strings.ReplaceAll(cookieStr, "\"", "")
				viper.Set("credentials", cookieStr)
				err := viper.WriteConfig()
				if err != nil {
					fmt.Fprintf(os.Stderr, "error attempting to write to config file: %s", err)
					page.Close()
					browser.Close()

					wg.Done()
					os.Exit(1)
				}
				fmt.Println("✅ Successfully updated configuration")

				page.Close()
				browser.Close()
				wg.Done()
			}
		})
		wait()
	}()
	wg.Wait()
	fmt.Println("Press Enter to complete login")
	utils.E(fmt.Scanln())
	fmt.Println("✅ Successfully logged in")
	os.Exit(0)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login through Stainless API",
	Run:   runLoginCmd,
}
