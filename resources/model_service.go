/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Service struct {
	Key
	Attributes ServiceAttributes `json:"attributes"`
}
type ServiceResponse struct {
	Data     Service  `json:"data"`
	Included Included `json:"included"`
}

type ServiceListResponse struct {
	Data     []Service `json:"data"`
	Included Included  `json:"included"`
	Links    *Links    `json:"links"`
}

// MustService - returns Service from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustService(key Key) *Service {
	var service Service
	if c.tryFindEntry(key, &service) {
		return &service
	}
	return nil
}
