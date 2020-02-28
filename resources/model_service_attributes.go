/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type ServiceAttributes struct {
	// defines how often Traefik should check service's health
	HealthcheckInterval *string `json:"healthcheck_interval,omitempty"`
	// defines how much time Traefik will wait for the [200;300) answer from the server before deleting it from routing tables
	HealthcheckTimout *string `json:"healthcheck_timout,omitempty"`
	Name              string  `json:"name"`
	Port              string  `json:"port"`
	Rule              string  `json:"rule"`
	// priority of the passed rule which will override default rule priority in Traefik. Default priority is len(rule).
	RulePriority *int32 `json:"rule_priority,omitempty"`
	Url          string `json:"url"`
}
