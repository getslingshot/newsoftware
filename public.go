package main

// Facade used to export single Public method
type Facade interface {
	Public() interface{}
}

// Public function for our structs
// Returns Public fields defined to `Public` method associated to a struct
func Public(o interface{}) interface{} {
	if p, ok := o.(Facade); ok {
		return p.Public()
	}

	return o
}
