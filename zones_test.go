package powerdns_test

import (
	"github.com/joeig/go-powerdns"
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
	"testing"
)

func TestGetZones(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "http://localhost:8080/api/v1/servers/localhost/zones",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				zonesMock := []powerdns.Zone{
					{
						ID:             "example.com.",
						Name:           "example.com.",
						URL:            "/api/v1/servers/localhost/zones/example.com.",
						Kind:           "Native",
						Serial:         1337,
						NotifiedSerial: 1337,
					},
				}
				return httpmock.NewJsonResponse(200, zonesMock)
			} else {
				return httpmock.NewStringResponse(401, "Unauthorized"), nil
			}
		},
	)

	p := powerdns.NewClient("http://localhost:8080/", "localhost", "example.com", "apipw")
	zones, err := p.GetZones()
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(*zones) == 0 {
		t.Error("Received amount of statistics is 0")
	}
}

func TestGetZone(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "http://localhost:8080/api/v1/servers/localhost/zones/example.com",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				zoneMock := powerdns.Zone{
					ID:   "example.com.",
					Name: "example.com.",
					URL:  "/api/v1/servers/localhost/zones/example.com.",
					Kind: "Native",
					RRsets: []powerdns.RRset{
						{
							Name: "example.com.",
							Type: "SOA",
							TTL:  3600,
							Records: []powerdns.Record{
								{
									Content: "a.misconfigured.powerdns.server. hostmaster.example.com. 1337 10800 3600 604800 3600",
								},
							},
						},
					},
					Serial:         1337,
					NotifiedSerial: 1337,
				}
				return httpmock.NewJsonResponse(200, zoneMock)
			} else {
				return httpmock.NewStringResponse(401, "Unauthorized"), nil
			}
		},
	)

	p := powerdns.NewClient("http://localhost:8080/", "localhost", "example.com", "apipw")
	zone, err := p.GetZone()
	if err != nil {
		t.Errorf("%s", err)
	}
	if zone.ID != "example.com." {
		t.Error("Received no zone")
	}
}

func TestNotify(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("PUT", "http://localhost:8080/api/v1/servers/localhost/zones/example.com/notify",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				return httpmock.NewStringResponse(200, "{\"result\":\"Notification queued\"}"), nil
			} else {
				return httpmock.NewStringResponse(401, "Unauthorized"), nil
			}
		},
	)

	p := powerdns.NewClient("http://localhost:8080/", "localhost", "example.com", "apipw")
	notifyResult, err := p.Notify()
	if err != nil {
		t.Errorf("%s", err)
	}
	if notifyResult.Result != "Notification queued" {
		t.Error("Notification was not queued successfully")
	}
}