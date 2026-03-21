package raindrop

// apiCollectionRef is the JSON representation of a collection reference in Raindrop API responses.
// Used by both bookmark and collection repositories.
type apiCollectionRef struct {
	ID int `json:"$id"`
}
