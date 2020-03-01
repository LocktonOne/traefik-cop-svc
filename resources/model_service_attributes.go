/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type ServiceAttributes struct {
	Name string `json:"name"`
	Port string `json:"port"`
	Rule string `json:"rule"`
	// priority of the passed rule which will override default rule priority in Traefik. Default priority is len(rule).
	RulePriority *int32 `json:"rule_priority,omitempty"`
	Url          string `json:"url"`
}
