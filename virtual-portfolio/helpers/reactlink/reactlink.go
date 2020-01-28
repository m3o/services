package reactlink

import "fmt"

// Create returns an encoded link
func Create(resourceType, resourceUUID, resourceName string) string {
	return fmt.Sprintf("<&%v:%v>%v<&/%v>", resourceType, resourceUUID, resourceName, resourceType)
}
