package apigee

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func convertSetToArray(set *schema.Set) []string {
	setList := set.List()
	retVal := []string{}
	for _, s := range setList {
		retVal = append(retVal, s.(string))
	}
	return retVal
}
