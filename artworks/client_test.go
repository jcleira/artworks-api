package artworks

import (
	"fmt"
	"reflect"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestCheckValidDate(t *testing.T) {
	tests := []struct {
		date          int64
		expectedError error
	}{
		{
			date:          1489140630, // A valid date
			expectedError: nil,
		},
		{
			date:          1480000000, // A non valid date
			expectedError: fmt.Errorf("Unable to query Snapshots from over one month"),
		},
		{
			date:          -1, // A non valid date
			expectedError: fmt.Errorf("Unable to query Snapshots from over one month"),
		},
	}

	for _, test := range tests {
		err := ValidateDate(test.date)
		if test.expectedError != nil {
			if err.Error() != test.expectedError.Error() {
				t.Errorf("The returned error or no error from GetSnapshots don't math the test case. Got: %s Expected: %s", err, test.expectedError)
				return
			}
		}
	}
}

func TestGetSnapshots(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Unable to open a stub database connection. Err %s", err)
	}

	defer db.Close()

	snapshotsClient := Client{
		DB: db,
	}

	tests := []struct {
		date              int64
		expectedSnapshots []Snapshot
		expectedError     error
	}{
		{
			date: 1489140630, // A valid date
			expectedSnapshots: []Snapshot{
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
			},
			expectedError: nil,
		},
		{
			date:              1480000000, // A non valid date
			expectedSnapshots: nil,
			expectedError:     fmt.Errorf("Unable to query Snapshots from over one month"),
		},
	}

	for _, test := range tests {
		columns := []string{
			"id", "code", "created_at", "cups", "wins",
			"card_1", "card_2", "card_3", "card_4",
			"card_5", "card_6", "card_7", "card_8"}

		mock.ExpectQuery("SELECT (.+) FROM player_snapshots where created_at >= ?").
			WithArgs(test.date).
			WillReturnRows(sqlmock.NewRows(columns).
				AddRow(1, "#EU82REE", 1489140631, 2839, 232,
					"giant_skeleton", "goblins_barrel", "pekka", "princess",
					"ice_wizard", "rocket", "mortar", "prince").
				AddRow(2, "#F423432", 1489140633, 3839, 432,
					"royal_giant", "princess", "poison", "miner",
					"furnace", "rage", "zap", "prince"))

		snapshots, err := snapshotsClient.GetSnapshots(test.date)
		if test.expectedError != nil {
			if err.Error() != test.expectedError.Error() {
				t.Errorf("Returned error from GetSnapshots don't math the test case. Got: %s Expected: %s", err, test.expectedError)
				return
			}
		}

		if !reflect.DeepEqual(snapshots, test.expectedSnapshots) {
			t.Errorf("The returned Snapshots from GetSnapshots() don't match the expected. Got: %v Expected: %v", snapshots, test.expectedSnapshots)
		}

		if test.expectedError == nil {
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expections: %s", err)
				return
			}
		}
	}
}

func TestAddUpdateSnapshot(t *testing.T) {
	tests := []struct {
		action        string
		snapshot      *Snapshot
		expectedError error
	}{
		{
			action: "INSERT",
			snapshot: &Snapshot{
				Code: "#EU82REE", CreatedAt: 1489140631, Cups: 2839, Wins: 232,
				Card1: "giant_skeleton", Card2: "goblins_barrel", Card3: "pekka", Card4: "princess",
				Card5: "ice_wizard", Card6: "rocket", Card7: "mortar", Card8: "prince",
			},
			expectedError: nil,
		},
		{
			action: "UPDATE",
			snapshot: &Snapshot{
				ID: 2, Code: "#F423432", CreatedAt: 1489140633, Cups: 3839, Wins: 432,
				Card1: "royal_giant", Card2: "princess", Card3: "poison", Card4: "miner",
				Card5: "furnace", Card6: "rage", Card7: "zap", Card8: "prince",
			},
			expectedError: nil,
		},
		{
			action:        "FOO",
			expectedError: fmt.Errorf("The given action is not valid, it should be either INSERT or UPDATE"),
		},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Unable to open a stub database connection. Err %s", err)
	}
	defer db.Close()

	snapshotsClient := Client{
		DB: db,
	}

	for _, test := range tests {
		if test.expectedError == nil {
			switch test.action {
			case "INSERT":
				mock.ExpectPrepare("INSERT INTO player_snapshots").ExpectExec().WithArgs(
					test.snapshot.Code, test.snapshot.Cups, test.snapshot.Wins,
					test.snapshot.Card1, test.snapshot.Card2, test.snapshot.Card3,
					test.snapshot.Card4, test.snapshot.Card5, test.snapshot.Card6,
					test.snapshot.Card7, test.snapshot.Card8, test.snapshot.CreatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1))
				break
			case "UPDATE:":
				mock.ExpectPrepare("UPDATE player_snapshots SET").ExpectExec().WithArgs(
					test.snapshot.Code, test.snapshot.Cups, test.snapshot.Wins,
					test.snapshot.Card1, test.snapshot.Card2, test.snapshot.Card3,
					test.snapshot.Card4, test.snapshot.Card5, test.snapshot.Card6,
					test.snapshot.Card7, test.snapshot.Card8).
					WillReturnResult(sqlmock.NewResult(0, 1))
			}
		}

		err = snapshotsClient.AddUpdateSnapshot(test.action, test.snapshot)
		if test.expectedError != nil {
			if err.Error() != test.expectedError.Error() {
				t.Errorf("The returned error or no error from AddUpdateSnapshot don't math the test case. Got: %s Expected: %s", err, test.expectedError)
				return
			}
		}

		if test.expectedError == nil {
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expections: %s", err)
				return
			}
		}
	}
}

func TestDeleteSnapshot(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Unable to open a stub database connection. Err %s", err)
	}
	defer db.Close()

	snapshotsClient := Client{
		DB: db,
	}

	mock.ExpectPrepare("DELETE FROM player_snapshots").ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

	err = snapshotsClient.DeleteSnapshot(1)
	if err != nil {
		t.Errorf("DeleteSnapshot returned a non expected error. Err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expections: %s", err)
		return
	}
}
