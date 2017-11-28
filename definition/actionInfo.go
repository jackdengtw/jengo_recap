package definition

type ActionInfo struct {
	// Action Info from yml
	// executableAction instance = definition/actionInfo + action/concreteAction
	Name        string
	ClassName   string // Reservation
	PropertyMap map[string]interface{}
}
