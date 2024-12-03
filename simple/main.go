package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/Jeffail/tunny"
)

func SendEmail(email string, subject string, body string) {
	fmt.Printf("Sending email to %s\n", email)
	fmt.Printf("Subject %s\n Body: %s\n", subject, body)
	// Simulate sending email
	time.Sleep(2 * time.Second)
}

func main() {
	numCPUs := runtime.NumCPU()
	fmt.Printf("Number of CPUs: %d\n\n", numCPUs)

	pool := tunny.NewFunc(numCPUs, func(payload interface{}) interface{} {
		m, ok := payload.(map[string]interface{})
		if !ok {
			return fmt.Errorf("Unable to extract map")
		}

		// Extract the fields
		email, ok := m["email"].(string)
		if !ok {
			return fmt.Errorf("email field is missing or not a string")
		}

		subject, ok := m["subject"].(string)
		if !ok {
			return fmt.Errorf("subject field is missing or not a string")
		}

		body, ok := m["body"].(string)
		if !ok {
			return fmt.Errorf("body field is missing or not a string")
		}

		SendEmail(email, subject, body)

		return nil
	})
	defer pool.Close()

	for i := 0; i < 100; i++ {
		var data interface{} = map[string]interface{}{
			"email":   "email" + strconv.Itoa(i+1) + "@codeheim.io",
			"subject": "Welcome!",
			"body":    "Thank you for signing up.",
		}

		go func() {
			result := pool.Process(data)
			if result == nil {
				fmt.Println("Mail sent!")
			}
		}()
	}

	for {
		qLen := pool.QueueLength()
		fmt.Printf(" --------------------------------------- Queue length: %d\n", qLen)
		if qLen == 0 {
			break
		}
		time.Sleep(time.Second)
	}

	time.Sleep(3 * time.Second)

}
