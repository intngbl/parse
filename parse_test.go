package parse

import (
	"os"
	"testing"
)

const (
	deviceToken = `0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef`
)

var (
	applicationID string
	clientKey     string
	objectID      string
)

func init() {
	applicationID = os.Getenv("APPLICATION_ID")
	clientKey = os.Getenv("REST_API_KEY")
}

func TestPostInstallationMessage(t *testing.T) {
	app := New(applicationID, clientKey)
	res, err := app.Installation(InstallationMessage{
		DeviceToken: deviceToken,
		DeviceType:  IOS,
		Channels:    []string{""},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.DeviceToken != deviceToken {
		t.Fatal("Expecting a device token.")
	}
	objectID = res.ObjectID
}

func TestPostInstallationMessageBroadcast(t *testing.T) {
	app := New(applicationID, clientKey)
	res, err := app.Installation(InstallationMessage{
		DeviceToken: deviceToken,
		DeviceType:  IOS,
		Channels:    []string{"broadcast"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.DeviceToken != deviceToken {
		t.Fatal("Expecting a device token.")
	}
}

func TestPushMessageToChannel(t *testing.T) {
	app := New(applicationID, clientKey)
	res, err := app.Push(PushMessage{
		Channels: []string{
			"Sports",
		},
		Data: Notification{
			Alert: "The Mets scored! The game is now tied 1-1.",
			Badge: Increment,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Result {
		t.Fatal("Expecting a true result.")
	}
}

func TestPushMessageToUserID(t *testing.T) {
	app := New(applicationID, clientKey)
	res, err := app.Push(PushMessage{
		Where: map[string]interface{}{
			"objectId": objectID,
		},
		Data: Notification{
			Alert: "The Mets scored! The game is now tied 1-1.",
			Badge: Increment,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Result {
		t.Fatal("Expecting a true result.")
	}
}
