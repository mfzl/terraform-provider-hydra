package main

type validStrOptions map[string]bool

func (v validStrOptions) keys() []string {
	keys := []string{}

	for k := range v {
		keys = append(keys, k)
	}

	return keys
}

func toStringSlice(vals []interface{}) []string {
	strSlice := []string{}

	for _, v := range vals {
		strSlice = append(strSlice, v.(string))
	}

	return strSlice
}
