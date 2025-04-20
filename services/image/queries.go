package image

var (
	createImageQuery   = `INSERT INTO images(id, name, size, url) VALUES($1, $2, $3, $4);`
	getImageByUrlQuery = `SELECT id, name, size, url FROM images WHERE url = $1;`
	getImageByIdQuery = `SELECT id, name, size, url FROM images WHERE id = $1;`
	updateImageQuery   = `UPDATE images SET url = $1 WHERE id = $2;`
)
