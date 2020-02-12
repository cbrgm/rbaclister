package xhttp

type APIResourceList struct {
	Kind         string `json:"kind"`
	APIVersion   string `json:"apiVersion"`
	GroupVersion string `json:"groupVersion"`
	Resources    []struct {
		Name               string   `json:"name"`
		SingularName       string   `json:"singularName"`
		Namespaced         bool     `json:"namespaced"`
		Kind               string   `json:"kind"`
		Verbs              []string `json:"verbs"`
		StorageVersionHash string   `json:"storageVersionHash,omitempty"`
		ShortNames         []string `json:"shortNames,omitempty"`
		Categories         []string `json:"categories,omitempty"`
		Group              string   `json:"group,omitempty"`
		Version            string   `json:"version,omitempty"`
	} `json:"resources"`
}

func (r *APIResourceList) AsRules() []Rule {
	var rules []Rule
	for _, resource := range r.Resources {
		var rule = Rule{
			APIGroups: []string{
				r.GroupVersion,
			},
			Resources: []string{
				resource.Name,
			},
			Verbs: resource.Verbs,
		}
		rules = append(rules, rule)
	}
	return rules
}
