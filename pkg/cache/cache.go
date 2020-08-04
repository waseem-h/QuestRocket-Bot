package cache

type FilterFunc func(interface{}, interface{}) bool

type Cache interface {
    Add(object interface{}) bool
    Contains(object interface{}) bool
    GetAll() []interface{}
    Remove(object interface{}) bool
    IndexOf(object interface{}) int
    Filter(filterFunc FilterFunc, compareTo interface{}) []interface{}
}

type Caching interface {
    PopulateCache() error
    GetCache() Cache
}
