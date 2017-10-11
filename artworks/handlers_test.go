package artworks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	"github.com/xeipuuv/gojsonschema"
)

func ConfigureSnapshotSchema() (*gojsonschema.Schema, error) {
	jsonSchemaData, err := ioutil.ReadFile("./snapshot_schema.json")
	if err != nil {
		return nil, fmt.Errorf("Unable to read the Snapshot JSON schema file. Err: %s", err)
	}

	schemaLoader := gojsonschema.NewStringLoader(string(jsonSchemaData))
	snapshotSchema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return nil, fmt.Errorf("Unable to load the Snapshot JSON schema. Err: %s", err)
	}

	return snapshotSchema, nil
}

func TestGetSnapshotsHandler(t *testing.T) {
	server := httptest.NewServer(GetSnapshotsHandler(&FakeClient{}))
	defer server.Close()

	tests := []struct {
		params     string
		statusCode int
	}{
		{
			params:     "?date=1489140631",
			statusCode: http.StatusOK,
		},
		{
			params:     "?date=foo",
			statusCode: http.StatusBadRequest,
		},
		{
			params:     "?date=",
			statusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		r, err := http.Get(fmt.Sprint(server.URL, test.params))
		if err != nil {
			t.Error("Unable to perform GetSnapshots request. Err: %s", err)
			return
		}

		defer r.Body.Close()

		if r.StatusCode != test.statusCode {
			t.Errorf("The response Status Code don't match Got: %d Expected: %d", r.StatusCode, test.statusCode)
			return
		}

		if test.statusCode == http.StatusOK {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				panic(err)
			}

			var snapshots []Snapshot
			err = json.Unmarshal(body, &snapshots)
			if err != nil {
				t.Errorf("Unable to Unmarshal the Body content of the GetSnapshots request: %s", err)
				return
			}

			if len(snapshots) != 2 {
				t.Errorf("The returned number of Snapshots don't match the expected amount. Got: %d Expected: 3", len(snapshots))
				return
			}

			expectedSnapshots := []Snapshot{
				{
					ID: 1, Code: "#EU82REE", CreatedAt: 1489140631, Cups: 2839, Wins: 232,
					Card1: "giant_skeleton", Card2: "goblins_barrel", Card3: "pekka", Card4: "princess",
					Card5: "ice_wizard", Card6: "rocket", Card7: "mortar", Card8: "prince",
				},
				{
					ID: 2, Code: "#F423432", CreatedAt: 1489140633, Cups: 3839, Wins: 432,
					Card1: "royal_giant", Card2: "princess", Card3: "poison", Card4: "miner",
					Card5: "furnace", Card6: "rage", Card7: "zap", Card8: "prince",
				},
			}

			if !reflect.DeepEqual(snapshots, expectedSnapshots) {
				t.Errorf("The returned Snapshots don't match the expected. Got: %v Expected: %v", snapshots, expectedSnapshots)
				return
			}
		}
	}
}

func TestAddSnapshotHandler(t *testing.T) {
	snapshotJSON := fmt.Sprint(
		"{ \"code\": \"#EU82REEE\", \"created_at\": 1489140631, \"cups\": 2839, \"wins\": 232,",
		"\"card_1\": \"giant_skeleton\", \"card_2\": \"goblins_barrel\", \"card_3\": \"pekka\",",
		"\"card_4\": \"princess\", \"card_5\": \"ice_wizard\", \"card_6\": \"rocket\",",
		"\"card_7\": \"mortar\", \"card_8\": \"prince\" }")

	snapshotSchema, err := ConfigureSnapshotSchema()
	if err != nil {
		t.Error("Unable to load the Snapshot Schema. Err: %s", err)
		return
	}

	server := httptest.NewServer(AddSnapshotHandler(&FakeClient{}, snapshotSchema))
	defer server.Close()

	tests := []struct {
		snapshotJSON []byte
		statusCode   int
	}{
		{
			snapshotJSON: []byte(snapshotJSON),
			statusCode:   201,
		},
		{
			snapshotJSON: []byte("{}"),
			statusCode:   400,
		},
		{
			snapshotJSON: []byte("{ \"code\": \"#EU82REEE\", \"created_at\": 1489140631, \"cups\": 2839, \"wins\": 232 }"),
			statusCode:   400,
		},
	}

	for _, test := range tests {
		req, err := http.NewRequest(http.MethodPut, server.URL, bytes.NewBuffer(test.snapshotJSON))
		if err != nil {
			t.Error("Unable to perform AddSnapshot request. Err: %s", err)
			return
		}

		resp, err := http.DefaultClient.Do(req)

		if resp.StatusCode != test.statusCode {
			t.Errorf("The response status code don't match the expected. Got: %d Expected: %d", resp.StatusCode, test.statusCode)
			return
		}
	}
}

func TestUpdateSnapshotHandler(t *testing.T) {
	snapshotJSON := fmt.Sprint(
		"{ \"id\": 1, \"code\": \"#EU82REE\", \"date\": 1489140631, \"cups\": 2839, \"wins\": 232,",
		"\"card_1\": \"giant_skeleton\", \"card_2\": \"goblins_barrel\", \"card_3\": \"pekka\",",
		"\"card_4\": \"princess\", \"card_5\": \"ice_wizard\", \"card_6\": \"rocket\",",
		"\"card_7\": \"mortar\", \"card_8\": \"prince\" }")

	snapshotSchema, err := ConfigureSnapshotSchema()
	if err != nil {
		t.Error("Unable to load the Snapshot Schema. Err: %s", err)
		return
	}

	r := mux.NewRouter()
	r.Handle("/players/snapshots/{id:[0-9]+}", UpdateSnapshotHandler(&FakeClient{}, snapshotSchema))

	server := httptest.NewServer(r)
	defer server.Close()

	tests := []struct {
		url          string
		snapshotJSON []byte
		statusCode   int
	}{
		{
			url:          "/players/snapshots/1", // ok
			snapshotJSON: []byte(snapshotJSON),
			statusCode:   http.StatusNoContent,
		},
		{
			url:          "/players/snapshots/2", // the snapshot id don't match the JSON one
			snapshotJSON: []byte(snapshotJSON),
			statusCode:   http.StatusBadRequest,
		},
		{
			url:        "/players/snapshots/foo", // non numeric id on url
			statusCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		req, err := http.NewRequest(http.MethodPut, fmt.Sprint(server.URL, test.url), bytes.NewBuffer(test.snapshotJSON))
		if err != nil {
			t.Error("Unable to perform UpdateSnapshot request. Err: %s", err)
			return
		}

		resp, err := http.DefaultClient.Do(req)

		if resp.StatusCode != test.statusCode {
			t.Errorf("The response status code don't match the expected. Got: %d Expected: %d", resp.StatusCode, test.statusCode)
			return
		}
	}
}

func TestDeleteSnapshotHandler(t *testing.T) {
	tests := []struct {
		url          string
		snapshotJSON []byte
		statusCode   int
	}{
		{
			url:        "/players/snapshots/1", // ok
			statusCode: http.StatusNoContent,
		},
		{
			url:        "/players/snapshots/foo", // non numeric id on url
			statusCode: http.StatusNotFound,
		},
	}

	r := mux.NewRouter()
	r.Handle("/players/snapshots/{id:[0-9]+}", DeleteSnapshotHandler(&FakeClient{}))

	server := httptest.NewServer(r)
	defer server.Close()

	for _, test := range tests {
		req, err := http.NewRequest(http.MethodPut, fmt.Sprint(server.URL, test.url), bytes.NewBuffer(test.snapshotJSON))
		if err != nil {
			t.Error("Unable to perform DeleteSnapshot request. Err: %s", err)
			return
		}

		resp, err := http.DefaultClient.Do(req)

		if resp.StatusCode != test.statusCode {
			t.Errorf("The response status code don't match the expected. Got: %d Expected: %d", resp.StatusCode, test.statusCode)
			return
		}
	}
}
