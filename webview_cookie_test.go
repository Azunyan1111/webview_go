package webview

import (
	"runtime"
	"testing"
	"time"
)

func TestGetCookies(t *testing.T) {
	// Skip on non-macOS platforms
	if runtime.GOOS != "darwin" {
		t.Skip("GetCookies is only supported on macOS")
	}

	w := New(false)
	if w == nil {
		t.Fatal("Failed to create webview")
	}
	defer w.Destroy()

	// Navigate to a test page that sets cookies
	testHTML := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Cookie Test</title>
	</head>
	<body>
		<h1>Cookie Test Page</h1>
		<script>
			// Set some test cookies
			document.cookie = "test1=value1; path=/";
			document.cookie = "test2=value2; path=/; secure";
			document.cookie = "test3=value3; path=/test";
			
			// Session cookie (no expiration)
			document.cookie = "session=active; path=/";
			
			// Cookie with expiration
			var date = new Date();
			date.setTime(date.getTime() + (24*60*60*1000)); // 1 day
			document.cookie = "persistent=data; expires=" + date.toUTCString() + "; path=/";
		</script>
	</body>
	</html>
	`

	// Set HTML and wait for it to load
	w.SetHtml(testHTML)
	
	// Give time for cookies to be set
	time.Sleep(100 * time.Millisecond)

	// Test in a goroutine since Run blocks
	done := make(chan bool)
	go func() {
		// Wait a bit for the webview to initialize
		time.Sleep(500 * time.Millisecond)

		// Get cookies
		cookies, err := w.GetCookies()
		if err != nil {
			t.Errorf("Failed to get cookies: %v", err)
		} else {
			t.Logf("Retrieved %d cookies", len(cookies))
			
			// Check if we got some cookies
			if len(cookies) == 0 {
				t.Error("Expected to get some cookies, but got none")
			}
			
			// Log cookie details
			for _, cookie := range cookies {
				t.Logf("Cookie: %s=%s (domain=%s, path=%s, secure=%v, httpOnly=%v)",
					cookie.Name, cookie.Value, cookie.Domain, cookie.Path, 
					cookie.Secure, cookie.HTTPOnly)
			}
			
			// Look for our test cookies
			foundTest1 := false
			for _, cookie := range cookies {
				if cookie.Name == "test1" && cookie.Value == "value1" {
					foundTest1 = true
					break
				}
			}
			
			if !foundTest1 {
				t.Error("Expected to find test1 cookie")
			}
		}

		// Terminate the webview
		w.Terminate()
		done <- true
	}()

	// Run the webview (this blocks)
	w.Run()
	
	// Wait for test to complete
	<-done
}

func TestGetCookiesEmpty(t *testing.T) {
	// Skip on non-macOS platforms
	if runtime.GOOS != "darwin" {
		t.Skip("GetCookies is only supported on macOS")
	}

	w := New(false)
	if w == nil {
		t.Fatal("Failed to create webview")
	}
	defer w.Destroy()

	// Navigate to about:blank (no cookies)
	w.Navigate("about:blank")
	
	// Test in a goroutine since Run blocks
	done := make(chan bool)
	go func() {
		// Wait a bit for the webview to initialize
		time.Sleep(500 * time.Millisecond)

		// Get cookies
		cookies, err := w.GetCookies()
		if err != nil {
			t.Errorf("Failed to get cookies: %v", err)
		} else {
			t.Logf("Retrieved %d cookies from blank page", len(cookies))
			// It's okay to have 0 cookies on a blank page
		}

		// Terminate the webview
		w.Terminate()
		done <- true
	}()

	// Run the webview (this blocks)
	w.Run()
	
	// Wait for test to complete
	<-done
}