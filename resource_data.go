package main

import "github.com/hashicorp/terraform/helper/schema"

// GetString pulls a string out of ResourceData.
func GetString(d *schema.ResourceData, key string) string {
	return d.Get(key).(string)
}

// GetInt pulls an int out of ResourceData.
func GetInt(d *schema.ResourceData, key string) int {
	return d.Get(key).(int)
}

// GetUint pulls a uint out of ResourceData.
func GetUint(d *schema.ResourceData, key string) uint {
	return uint(d.Get(key).(int))
}

// GetBool pulls a bool out of ResourceData.
func GetBool(d *schema.ResourceData, key string) bool {
	return d.Get(key).(bool)
}

// GetSliceString pulls a slice of strings out of ResourceData.
func GetSliceString(d *schema.ResourceData, key string) []string {
	slice := d.Get(key).([]interface{})

	stringSlice := make([]string, len(slice))
	for i, str := range slice {
		stringSlice[i] = str.(string)
	}

	return stringSlice
}
