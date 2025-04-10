package installer

// Installer is used for fetching all kinds of external files - clients, configs etc.
// It's primarly used to encapsulate fetching logic for easier installations and mocking.
type Installer interface {
	// Fetch uses HTTP GET to fetch data from the url and returns the response body as bytes.
	Fetch()

	// InstallFile uses HTTP GET to fetch data from the url and writes it to dest.
	InstallFile(url string, dest string) (err error)
}
