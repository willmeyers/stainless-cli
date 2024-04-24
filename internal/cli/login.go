package cli

import (
	"flag"
	"fmt"
	"os"
	"stainless_cli/internal/stainless"
	"strings"
	"sync"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

type LoginCommand struct {
	Refresh bool
	Cookies string
}

func NewLoginCommand(stl *stainless.Stainless) Command {
	return &LoginCommand{}
}

func (cmd *LoginCommand) Parse(fs *flag.FlagSet) error {
	fs.StringVar(&cmd.Cookies, "cookies", "", "Cookie string of current web session")

	return fs.Parse(os.Args[2:])
}

func (cmd *LoginCommand) Exec(stl *stainless.Stainless) error {
	if cmd.Cookies != "" {
		stl.AuthCookies = cmd.Cookies
		return nil
	}

	cachedCookies := os.Getenv("STAINLESS_NEXTAUTH_SESSION_COOKIES")
	if cachedCookies != "" {
		stl.AuthCookies = cachedCookies
		return nil
	}

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
				fmt.Println("🔒 successfully logged in")
				fmt.Println("❗️ copy/paste and export environment variable below to start using CLI")
				fmt.Printf("STAINLESS_NEXTAUTH_SESSION_COOKIES=%s\n", cookieStr)
				page.Close()
				browser.Close()
				wg.Done()
			}
		})
		wait()
	}()
	wg.Wait()
	return nil
}
