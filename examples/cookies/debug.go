package main

import (
	"log"
	"runtime"
	"time"

	"github.com/webview/webview_go"
)

func main() {
	if runtime.GOOS != "darwin" {
		log.Fatal("Cookie functionality is only supported on macOS")
	}

	// Enable debug logging
	webview.Debug = true

	// Create a new webview
	w := webview.New(true) // true enables developer tools
	if w == nil {
		log.Fatal("Failed to create webview")
	}
	defer w.Destroy()

	// Set window properties
	w.SetTitle("Cookie Debug")
	w.SetSize(800, 600, webview.HintNone)

	// Navigate to a cookie test website
	w.Navigate("https://www.tohoho-web.com/cgi/wwwcook.cgi")

	// Inject JavaScript to check cookies
	w.Init(`
		window.checkCookies = function() {
			console.log('Document cookies:', document.cookie);
			return document.cookie;
		};
	`)

	// Try to get cookies after page loads
	go func() {
		// Wait for page to load and cookies to be set
		time.Sleep(3 * time.Second)

		// Check JavaScript cookies
		w.Eval(`console.log('Checking cookies from JavaScript...');`)
		w.Eval(`console.log('document.cookie =', document.cookie);`)

		log.Println("Attempting to get cookies via GetCookies...")
		cookies, err := w.GetCookies()
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			log.Printf("Success: Got %d cookies", len(cookies))
			for _, cookie := range cookies {
				log.Printf("  %s = %s (domain: %s)", cookie.Name, cookie.Value, cookie.Domain)
			}
		}

		// Try again after a delay
		time.Sleep(3 * time.Second)
		log.Println("Second attempt to get cookies...")
		cookies2, err2 := w.GetCookies()
		if err2 != nil {
			log.Printf("Error on second attempt: %v", err2)
		} else {
			log.Printf("Success on second attempt: Got %d cookies", len(cookies2))
			for _, cookie := range cookies2 {
				log.Printf("  %s = %s (domain: %s)", cookie.Name, cookie.Value, cookie.Domain)
			}
		}
	}()

	// Run the webview
	log.Println("Starting webview...")
	w.Run()
}
