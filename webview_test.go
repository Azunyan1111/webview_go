package webview

import (
	"flag"
	"log"
	"os"
	"testing"
)

func Example() {
	w := New(true)
	defer w.Destroy()
	w.SetTitle("Hello")
	w.Bind("noop", func() string {
		log.Println("hello")
		return "hello"
	})
	w.Bind("add", func(a, b int) int {
		return a + b
	})
	w.Bind("quit", func() {
		w.Terminate()
	})
	w.SetHtml(`<!doctype html>
		<html>
			<body>hello</body>
			<script>
				window.onload = function() {
					document.body.innerText = ` + "`hello, ${navigator.userAgent}`" + `;
					noop().then(function(res) {
						console.log('noop res', res);
						add(1, 2).then(function(res) {
							console.log('add res', res);
							quit();
						});
					});
				};
			</script>
		</html>
	)`)
	w.Run()
}

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Verbose() {
		Example()
	}
	os.Exit(m.Run())
}

func TestClearCookies(t *testing.T) {
	// This test needs to run manually with -v flag as it creates a GUI window
	// Example: go test -v -run TestClearCookies
	if !testing.Verbose() {
		t.Skip("Skipping TestClearCookies in non-verbose mode. Run with -v flag to enable.")
		return
	}

	// Test implementation for manual testing
	log.Println("TestClearCookies: This test requires manual verification.")
	log.Println("You should see a webview window that automatically tests cookie clearing.")

	// Create an example that demonstrates ClearCookies functionality
	w := New(false)
	if w == nil {
		t.Fatal("Failed to create webview")
	}
	defer w.Destroy()

	// Enable debug for this test
	Debug = true
	defer func() { Debug = false }()

	w.SetTitle("Cookie Clear Test")

	w.Bind("testComplete", func(message string) {
		log.Println("Test Result:", message)
		w.Terminate()
	})

	w.Bind("clearCookiesFromGo", func() error {
		log.Println("Clearing cookies from Go...")
		err := w.ClearCookies()
		if err != nil {
			log.Println("ClearCookies error:", err)
		} else {
			log.Println("ClearCookies completed successfully")
		}
		return err
	})

	// HTML that demonstrates cookie clearing
	w.SetHtml(`<!doctype html>
		<html>
		<head>
			<style>
				body { font-family: Arial; padding: 20px; }
				button { margin: 10px; padding: 10px; }
				#status { margin: 20px 0; padding: 10px; background: #f0f0f0; }
			</style>
			<script>
			function setCookies() {
				document.cookie = "test1=value1; path=/";
				document.cookie = "test2=value2; path=/";
				document.cookie = "test3=value3; path=/";
				updateStatus("Cookies set: " + document.cookie);
			}
			
			function updateStatus(msg) {
				document.getElementById('status').innerText = msg;
			}
			
			async function clearAndVerify() {
				updateStatus("Clearing cookies...");
				try {
					await window.clearCookiesFromGo();
					// Wait a bit for cookies to be cleared
					await new Promise(resolve => setTimeout(resolve, 1000));
					
					if (document.cookie.length > 0) {
						updateStatus("FAIL: Cookies still present: " + document.cookie);
						testComplete("FAIL: Cookies were not cleared");
					} else {
						updateStatus("SUCCESS: All cookies cleared!");
						testComplete("SUCCESS: Cookies cleared successfully");
					}
				} catch (e) {
					updateStatus("ERROR: " + e.toString());
					testComplete("ERROR: " + e.toString());
				}
			}
			
			window.onload = function() {
				updateStatus("Ready to test cookie clearing");
			};
			</script>
		</head>
		<body>
			<h2>Cookie Clear Test</h2>
			<button onclick="setCookies()">1. Set Cookies</button>
			<button onclick="clearAndVerify()">2. Clear & Verify</button>
			<div id="status">Ready to test cookie clearing</div>
			<p>Click buttons in order to test cookie clearing functionality.</p>
		</body>
		</html>
	`)

	w.Run()
}
