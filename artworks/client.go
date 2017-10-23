package artworks

import (
	"database/sql"
	"fmt"
)

// Artwork is package's main struct, represents an Artwork information: (title,
// author, etc.). It's main objective is preservation and CRUD operations.
//
// Example:
// {
//   ID: 1,
//   Rei: '#elle',
//   CreatedAt: 1489140631,
//   Pro: 'Ayuntamiento de Mahón',
//   Ubi: 'Desconocido',
//   ...
//   Tecnica: '',
// }
type Artwork struct {
	ID        int    `json:"id"`
	Rei       string `json:"rei"`
	CreatedAt int64  `json:"created_at"`
	Ubi       string `json:"ubi"`
	Pro       string `json:"pro"`
	Adq       string `json:"adq"`
	Reg       string `json:"reg"`
	Nom       string `json:"nom"`
	Tit       string `json:"tit"`
	Aut       string `json:"aut"`
	Fec       string `json:"fec"`
	Lug       string `json:"lug"`
	Ico       string `json:"ico"`
	Icc       string `json:"icc"`
	Tip       string `json:"tip"`
	Tec       string `json:"tec"`
	Sop       string `json:"sop"`
	Mat       string `json:"mat"`
	Tin       string `json:"tin"`
	Dim       string `json:"dim"`
	Hue       string `json:"hue"`
	Ins       string `json:"ins"`
	Des       string `json:"des"`
	Est       string `json:"est"`
}

// Client is the Artworks struct that implements the ArtworksController
// interface, it does also has the proper DB configuration to access the
// Artworks data on the database.
type Client struct {
	DB *sql.DB
}

// ArtworksController interface define the required methods to implement
// in order to be able to manage Artworks.
type ArtworksController interface {
	GetArtwork(int) (*Artwork, error)
	GetArtworks() ([]Artwork, error)
	AddUpdateArtwork(string, *Artwork) error
	DeleteArtwork(int) error
}

// GetArtwork returns an Artwork (by it's id) stored in the database.
//
// id - The Artwork id to query on the database.
//
// Returns:
// An Artworks.
// An error otherwise.
func (c *Client) GetArtwork(id int) (*Artwork, error) {
	stmt, err := c.DB.Prepare("SELECT * FROM artworks WHERE id=$1")
	if err != nil {
		return nil, fmt.Errorf("Unable to prepare the artwork's query table. Err: %s", err)
	}

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, fmt.Errorf("Unable to query the artworks table. Err: %s", err)
	}

	defer rows.Close()

	var artwork Artwork

	if rows.Next() {
		err := rows.Scan(
			&artwork.ID,
			&artwork.Rei,
			&artwork.CreatedAt,
			&artwork.Ubi,
			&artwork.Pro,
			&artwork.Adq,
			&artwork.Reg,
			&artwork.Nom,
			&artwork.Tit,
			&artwork.Aut,
			&artwork.Fec,
			&artwork.Lug,
			&artwork.Ico,
			&artwork.Icc,
			&artwork.Tip,
			&artwork.Tec,
			&artwork.Sop,
			&artwork.Mat,
			&artwork.Tin,
			&artwork.Dim,
			&artwork.Hue,
			&artwork.Ins,
			&artwork.Des,
			&artwork.Est,
		)
		if err != nil {
			return nil, fmt.Errorf("Unable to map an Artwork data row. Err: %s", err)
		}
	} else {
		return nil, fmt.Errorf("Unable to find an Artwork with id: %d", id)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Unable to iterate on Artworks data. Err %s", err)
	}

	return &artwork, nil
}

// GetArtworks returns all the Artworks stored in the database, it may become
// slow as database grow another technique should be applied.
//
// Returns:
// An array of Artworks.
// An error otherwise.
func (c *Client) GetArtworks() ([]Artwork, error) {
	rows, err := c.DB.Query("SELECT * FROM artworks")
	if err != nil {
		return nil, fmt.Errorf("Unable to query the artworks table. Err: %s", err)
	}

	defer rows.Close()

	artworks := make([]Artwork, 0)

	for rows.Next() {
		var artwork Artwork
		err := rows.Scan(
			&artwork.ID,
			&artwork.Rei,
			&artwork.CreatedAt,
			&artwork.Ubi,
			&artwork.Pro,
			&artwork.Adq,
			&artwork.Reg,
			&artwork.Nom,
			&artwork.Tit,
			&artwork.Aut,
			&artwork.Fec,
			&artwork.Lug,
			&artwork.Ico,
			&artwork.Icc,
			&artwork.Tip,
			&artwork.Tec,
			&artwork.Sop,
			&artwork.Mat,
			&artwork.Tin,
			&artwork.Dim,
			&artwork.Hue,
			&artwork.Ins,
			&artwork.Des,
			&artwork.Est,
		)
		if err != nil {
			return nil, fmt.Errorf("Unable to map an Artwork data row. Err: %s", err)
		}

		artworks = append(artworks, artwork)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Unable to iterate on Artworks data. Err %s", err)
	}

	return artworks, nil
}

// AddUpdateArtwork stores an Artwork by performing the action specified on the
// action param.
//
// Available actions:
//
// 'INSERT': For new Artworks.
// 'UPDATE': For existing Artworks.
//
// It will return an error if the given action is not valid.
//
// action: One of the above.
// artwork: The artwork to save.
//
// Returns an error if any.
func (c *Client) AddUpdateArtwork(action string, artwork *Artwork) error {
	sqlStatement := ""
	switch action {
	case "INSERT":
		sqlStatement = fmt.Sprint(
			"INSERT INTO artworks",
			"(rei,ubi,pro,adq,reg,nom,tit,aut,fec,lug,ico,icc,tip,tec,sop,mat,tin,dim,hue,ins,des,est,created_at) ",
			"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? ,?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		break
	case "UPDATE":
		sqlStatement = fmt.Sprint(
			"UPDATE artworks SET ",
			"rei=?,ubi=?,pro=?,",
			"adq=?,reg=?,nom=?,tit=?,aut=?,fec=?,lug=?,ico=?,",
			"icc=?,tip=?,tec=?,sop=?,mat=?,tin=?,dim=?,hue=?,",
			"ins=?,des=?,est=? ",
			"WHERE id=?")
		break
	default:
		return fmt.Errorf("The given action is not valid, it should be either INSERT or UPDATE")
	}

	stmt, err := c.DB.Prepare(sqlStatement)
	if err != nil {
		return fmt.Errorf("Unable to prepare the Artworks INSERT or UPDATE statement. Err: %s", err)
	}

	var values = []interface{}{
		artwork.Rei,
		artwork.Ubi,
		artwork.Pro,
		artwork.Adq,
		artwork.Reg,
		artwork.Nom,
		artwork.Tit,
		artwork.Aut,
		artwork.Fec,
		artwork.Lug,
		artwork.Ico,
		artwork.Icc,
		artwork.Tip,
		artwork.Tec,
		artwork.Sop,
		artwork.Mat,
		artwork.Tin,
		artwork.Dim,
		artwork.Hue,
		artwork.Ins,
		artwork.Des,
		artwork.Est,
	}

	if action == "INSERT" {
		values = append(values, artwork.CreatedAt)
	}

	if action == "UPDATE" {
		values = append(values, artwork.ID)
	}

	res, err := stmt.Exec(values...)
	if err != nil {
		return fmt.Errorf("Unable to execute the Artwork INSERT or UPDATE statement. Err: %s", err)
	}

	if action == "INSERT" {
		// We are being permisive here with the error returned and we will return 0 as ID
		// on the struct if we can't fetch the ID of the inserted item.
		ID, _ := res.LastInsertId()
		artwork.ID = int(ID)
	}

	return nil
}

// DeleteArtwork deletes an Artwork from the database.
//
// ID: The ID of the Artwork to delete.
//
// Returns an error if any.
func (c *Client) DeleteArtwork(ID int) error {
	stmt, err := c.DB.Prepare("DELETE FROM artworks WHERE id=?")
	if err != nil {
		return fmt.Errorf("Unable to prepare the Artwork DELETE statement. Err: %s", err)
	}

	_, err = stmt.Exec(ID)
	if err != nil {
		return fmt.Errorf("Unable to execute the Artwork DELETE statement. Err: %s", err)
	}

	return nil
}
