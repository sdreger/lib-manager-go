package publisher

var (
	AllowedSortFields = []string{"id", "name"}
)

type LookupItem struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type lookupEntity struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}
