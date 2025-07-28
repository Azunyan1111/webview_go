package main

import (
	"log"
	"net/http"
	"time"
	"github.com/webview/webview_go"
)

func main() {
	// Start simple HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!doctype html>
<html>
<head><title>Cookie Test</title></head>
<body>
	<h1>Cookie Test Page</h1>
	<p>Check console for cookie operations</p>
</body>
</html>`))
	})

	go func() {
		log.Println("Starting HTTP server on http://localhost:8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	// Create webview
	log.Println("Creating webview...")
	w := webview.New(false)
	if w == nil {
		log.Fatal("Failed to create webview")
	}
	defer w.Destroy()

	webview.Debug = true

	w.SetTitle("Simple Cookie Test")
	w.SetSize(600, 400, webview.HintNone)

	// Navigate to a real HTTPS site
	w.Navigate("https://example.com")

	// Test cookies from goroutine
	go func() {
		log.Println("Waiting 3 seconds for page to load...")
		time.Sleep(3 * time.Second)

		// Test GetCookies
		log.Println("Testing GetCookies...")
		cookies, err := w.GetCookies()
		if err != nil {
			log.Printf("GetCookies error: %v", err)
		} else {
			log.Printf("GetCookies success: %d cookies", len(cookies))
		}

		// Test SetCookie
		log.Println("Testing SetCookie...")
		testCookie := webview.Cookie{
			Name:   "test_cookie",
			Value:  "test_value_123",
			Domain: ".example.com",
			Path:   "/",
		}
		
		err = w.SetCookie(testCookie)
		if err != nil {
			log.Printf("SetCookie error: %v", err)
		} else {
			log.Println("SetCookie success")
		}

		// Wait and get cookies again
		time.Sleep(1 * time.Second)
		
		log.Println("Getting cookies after set...")
		cookies, err = w.GetCookies()
		if err != nil {
			log.Printf("GetCookies error: %v", err)
		} else {
			log.Printf("GetCookies success: %d cookies", len(cookies))
			// Check if our cookie is there
			for _, c := range cookies {
				if c.Name == "test_cookie" {
					log.Printf("✓ Found our test cookie: %s=%s", c.Name, c.Value)
				}
			}
		}

		// Test ClearCookies
		log.Println("Testing ClearCookies...")
		err = w.ClearCookies()
		if err != nil {
			log.Printf("ClearCookies error: %v", err)
		} else {
			log.Println("ClearCookies success")
		}

		// Final check
		time.Sleep(1 * time.Second)
		cookies, err = w.GetCookies()
		if err != nil {
			log.Printf("Final GetCookies error: %v", err)
		} else {
			log.Printf("Final GetCookies: %d cookies (should be 0)", len(cookies))
		}

		log.Println("Test completed!")
		
		// Wait a bit and then terminate
		time.Sleep(2 * time.Second)
		log.Println("Terminating webview...")
		w.Terminate()
	}()

	w.Run()
}