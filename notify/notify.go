// Package notify handles the various was a Tenant might be notified.
// Typically via email and Text.
package notify

import "fmt"

// Notifier sends the message as a notification.
type Notifier interface {
	// "to" is the receiver, "msg" is the content.
	// "to" must be in the correct format as determined by the implementation.
	Notify(to, msg string) error
}

// Email sends notifications via smtp.
type Email struct {
	// smtp.Client // The host email account (the "sender").
}

// Notify via email.
func (email Email) Notify(to, msg string) error {
	return fmt.Errorf("unimplemented")
}

// SMS sends notifications via SMS.
type SMS struct{}

// Notify via sms.
func (sms SMS) Notify(to, msg string) error {
	return fmt.Errorf("unimplemented")
}

// Console sends notifications to the console (for debugging purposes).
type Console struct{}

// Notify via console log.
func (console Console) Notify(to, msg string) error {
	fmt.Printf("[%s]: %s\n", to, msg)
	return nil
}
