package apptimize

import "testing"

func TestApptimize(t *testing.T) {
	apptimize := New(&Config{
		APIToken: "<test-api-token>",
	})
	if isSuccess := t.Run("test for code block variants", func(t *testing.T) {
		if v, err := apptimize.Variant("test-user", "test-experiment"); err != nil {
			t.Error(err)
			return
		} else if v == "a" {
			t.Log("run code block a")
		} else if v == "b" {
			t.Log("run code block b")
		}
	}) && t.Run("test for tracking events", func(t *testing.T) {
		if err := apptimize.Track("test-user", "did-something-with-no-attributes"); err != nil {
			t.Error(err)
		} else if err := apptimize.Track("test-user", "did-something-with-map-attributes", map[string]interface{}{
			"property": "value",
		}); err != nil {
			t.Error(err)
		} else if err := apptimize.Track("test-user", "did-something-with-struct-attributes", &struct {
			Property string `json:"property"`
		}{
			Property: "value",
		}); err != nil {
			t.Error(err)
		} else if err := apptimize.Track("test-user", "did-something-with-bad-attributes", map[string]interface{}{
			"property": "value",
		}, &struct {
			Property string `json:"property"`
		}{
			Property: "value",
		}); err != ErrBadAttributes {
			t.Errorf("'%s' shoud be '%s'", err, ErrBadAttributes)
		}
	}); !isSuccess {
		t.Fatal("the apptimize sdk tests failed")
	}
}
