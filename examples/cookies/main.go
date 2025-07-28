package main

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/webview/webview_go"
)

func main() {
	if runtime.GOOS != "darwin" {
		log.Fatal("Cookie functionality is only supported on macOS")
	}

	// Create a new webview
	w := webview.New(true) // true enables developer tools
	if w == nil {
		log.Fatal("Failed to create webview")
	}
	defer w.Destroy()

	// Set window properties
	w.SetTitle("Cookie Example")
	w.SetSize(800, 600, webview.HintNone)

	// HTML page that sets and displays cookies
	html := `
<!DOCTYPE html>
<html>
<head>
	<title>Cookie Manager</title>
	<style>
		body {
			font-family: Arial, sans-serif;
			padding: 20px;
			max-width: 800px;
			margin: 0 auto;
		}
		button {
			margin: 5px;
			padding: 10px 15px;
			font-size: 16px;
			cursor: pointer;
		}
		#output {
			margin-top: 20px;
			padding: 15px;
			background: #f0f0f0;
			border-radius: 5px;
			white-space: pre-wrap;
			font-family: monospace;
		}
		.cookie-item {
			margin: 10px 0;
			padding: 10px;
			background: white;
			border-radius: 5px;
		}
	</style>
</head>
<body>
	<h1>Cookie Manager Example</h1>
	<p>This example demonstrates cookie management in WebView.</p>
	
	<div>
		<button onclick="setCookies()">Set Test Cookies</button>
		<button onclick="getCookies()">Get Cookies (via Go)</button>
		<button onclick="clearCookies()">Clear Cookies</button>
	</div>
	
	<div id="output"></div>
	
	<script>
		function setCookies() {
			// Set various types of cookies
			document.cookie = "username=john_doe; path=/";
			document.cookie = "theme=dark; path=/";
			document.cookie = "session_id=abc123; path=/";
			
			// Cookie with expiration (1 hour)
			var date = new Date();
			date.setTime(date.getTime() + (60*60*1000));
			document.cookie = "temp_data=temporary; expires=" + date.toUTCString() + "; path=/";
			
			// Secure cookie (only sent over HTTPS)
			document.cookie = "secure_token=xyz789; path=/; secure";
			
			updateOutput("Cookies have been set!");
			showBrowserCookies();
		}
		
		function getCookies() {
			// This will call the Go function to get cookies
			if (window.getCookiesFromGo) {
				window.getCookiesFromGo();
			} else {
				updateOutput("Go binding not available yet. Try again.");
			}
		}
		
		function clearCookies() {
			// Get all cookies and clear them
			var cookies = document.cookie.split(";");
			
			for (var i = 0; i < cookies.length; i++) {
				var cookie = cookies[i];
				var eqPos = cookie.indexOf("=");
				var name = eqPos > -1 ? cookie.substr(0, eqPos).trim() : cookie.trim();
				// Clear the cookie by setting it to expire in the past
				document.cookie = name + "=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
			}
			
			updateOutput("Cookies have been cleared!");
			showBrowserCookies();
		}
		
		function showBrowserCookies() {
			var output = "\nCurrent browser cookies:\n";
			output += document.cookie || "(no cookies)";
			updateOutput(output);
		}
		
		function updateOutput(text) {
			document.getElementById('output').textContent = text;
		}
		
		function displayCookiesFromGo(cookies) {
			var output = "Cookies retrieved from Go:\n\n";
			output += JSON.stringify(cookies, null, 2);
			updateOutput(output);
		}
		
		// Show current cookies on load
		window.onload = function() {
			showBrowserCookies();
		};
	</script>
</body>
</html>
`

	// Set the HTML
	w.SetHtml(html)

	// Bind a Go function to get cookies
	w.Bind("getCookiesFromGo", func() {
		cookies, err := w.GetCookies()
		if err != nil {
			w.Eval(fmt.Sprintf(`updateOutput("Error getting cookies: %v")`, err))
			return
		}

		// Convert cookies to a format suitable for JavaScript
		w.Eval(fmt.Sprintf(`displayCookiesFromGo(%v)`, toJSON(cookies)))
	})

	// Start a goroutine to get cookies once after initialization
	go func() {
		// Wait for WebView to fully initialize
		time.Sleep(2 * time.Second)

		cookies, err := w.GetCookies()
		if err != nil {
			log.Printf("Error getting cookies: %v", err)
		} else {
			log.Printf("Initial cookies (%d):", len(cookies))
			for _, cookie := range cookies {
				log.Printf("  %s = %s (domain: %s, path: %s, secure: %v)",
					cookie.Name, cookie.Value, cookie.Domain, cookie.Path, cookie.Secure)
			}
		}
	}()

	// Run the webview
	w.Run()
}

func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
