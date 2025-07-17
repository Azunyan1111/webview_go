package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Set various cookies
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    "abc123xyz",
			Path:     "/",
			HttpOnly: true,
		})
		
		http.SetCookie(w, &http.Cookie{
			Name:   "user_pref",
			Value:  "dark_mode",
			Path:   "/",
			Secure: false,
		})
		
		http.SetCookie(w, &http.Cookie{
			Name:    "temp_data",
			Value:   "temporary_value",
			Path:    "/",
			Expires: time.Now().Add(24 * time.Hour),
		})
		
		// HTML response
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <title>Cookie Test Server</title>
</head>
<body>
    <h1>Cookie Test Server</h1>
    <p>This page sets the following cookies:</p>
    <ul>
        <li>session_id = abc123xyz (HttpOnly)</li>
        <li>user_pref = dark_mode</li>
        <li>temp_data = temporary_value (expires in 24h)</li>
    </ul>
    <p>Current time: %s</p>
    <script>
        // Also set a client-side cookie
        document.cookie = "client_cookie=from_javascript; path=/";
        console.log("All cookies:", document.cookie);
    </script>
</body>
</html>
`, time.Now().Format(time.RFC3339))
	})
	
	log.Println("Test server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}