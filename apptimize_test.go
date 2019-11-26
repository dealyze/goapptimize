package apptimize

import (
	"fmt"
	"log"
	"testing"
)

func TestApptimize(t *testing.T) {

	// add your credentials here to run these tests
	apiToken := "<api-token>"
	experiment := "<experiment-with-two-variations>"

	// debug logging
	debug = true
	log.SetFlags(log.Llongfile)

	// unit tests
	apptimize := New(&Config{
		APIToken: apiToken,
	})
	if isSuccess := t.Run("test for code block variants", func(t *testing.T) {
		var baselineCount, variation1Count int
		for i := 0; i < 50; i++ {
			if v, err := apptimize.Variant(fmt.Sprintf("test-user-%d", i), experiment); err != nil {
				t.Error(err)
			} else if v == "baseline" {
				baselineCount++
			} else if v == "variation1" {
				variation1Count++
			}
		}
		t.Logf("%d basline variants and %d variation1", baselineCount, variation1Count)
		if passed := baselineCount > 0 && variation1Count > 0 && baselineCount+variation1Count == 50; !passed {
			t.Fail()
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
