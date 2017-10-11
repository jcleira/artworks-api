package artworks

import "fmt"

// FakeClient implements the SnapshotController interface, as the 'real'
// players.Client struct. It has been created for testing purposes.
type FakeClient struct{}

// GetArtworks return an array of mocked Snapshots if a valid date has been
// given.
func (tc *FakeClient) GetArtworks(date int64) ([]Artwork, error) {
	return []Artwork{
		{
			ID: 1, Rei: "#EU82REE", CreatedAt: 1489140631,
		},
		{
			ID: 2, Rei: "#F423432", CreatedAt: 1489140633,
		},
	}, nil
}

// AddUpdateArtwork return nil if the proper action was sent, error otherwise.
func (tc *FakeClient) AddUpdateArtwork(action string, artwork *Artwork) error {
	switch action {
	case "INSERT":
	case "UPDATE":
		return nil
	default:
		return fmt.Errorf("The given action is not valid, it should be either INSERT or UPDATE")
	}

	return nil
}

// DeleteArtwork return always nil.
func (tc *FakeClient) DeleteArtwork(ID int) error {
	return nil
}
