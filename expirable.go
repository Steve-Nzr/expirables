package expirables

import (
	"time"
)

// ExpirableRefresher is the function that is called each refresh of the Expirable Variable
type ExpirableRefresher func() interface{}

// Expirable in-memory variable
// Its value is an interface{} that is refreshed every {ttl}
// with the given refresher
type Expirable struct {
	refresher  ExpirableRefresher
	ttl        time.Duration
	value      interface{}
	expiration time.Time
}

func (v *Expirable) init() *Expirable {
	v.set(v.refresher())
	return v
}

func (v *Expirable) set(val interface{}) *Expirable {
	v.value = val
	v.expiration = time.Now().Add(v.ttl)
	return v
}

// Get the value of the stored variable.
// Calling this function could trigger a refresh on the value
// and potentially slow the function execution
func (v *Expirable) Get() interface{} {
	if time.Since(v.expiration) > 0 {
		return v.set(v.refresher()).value
	}

	return v.value
}

// NewExpirable creates a new Expirable variable with the given Refresher & TTL
func NewExpirable(refresher ExpirableRefresher, TTL time.Duration) *Expirable {
	exp := new(Expirable)
	exp.refresher = refresher
	exp.ttl = TTL

	return exp.init()
}
