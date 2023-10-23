package providerfwk

import "fmt"

type resourceJunosNameable interface {
	junosName() string
}

type stringValuable interface {
	ValueString() string
}

func defaultResourceSchemaDescription(
	rsc resourceJunosNameable,
) string {
	return "Provides a " + rsc.junosName() + "."
}

func defaultResourceCouldNotCreateWithEmptyMessage(
	rsc resourceJunosNameable, argName string,
) string {
	return "could not create " + rsc.junosName() + " with empty " + argName
}

func defaultResourceAlreadyExistsMessage(
	rsc resourceJunosNameable, str stringValuable,
) string {
	return fmt.Sprintf(rsc.junosName()+" %q already exists", str.ValueString())
}

func defaultResourceAlreadyExistsInRoutingInstanceMessage(
	rsc resourceJunosNameable, str stringValuable, routingInstance string,
) string {
	return fmt.Sprintf(rsc.junosName()+" %q already exists in routing-instance %q", str.ValueString(), routingInstance)
}

func defaultResourceDoesNotExistsAfterCommitMessage(
	rsc resourceJunosNameable, str stringValuable,
) string {
	return fmt.Sprintf(rsc.junosName()+" %q does not exists after commit "+
		"=> check your config", str.ValueString())
}

func defaultResourceDoesNotExistsInRoutingInstanceAfterCommitMessage(
	rsc resourceJunosNameable, str stringValuable, routingInstance string,
) string {
	return fmt.Sprintf(rsc.junosName()+" %q does not exists in routing-instance %q after commit "+
		"=> check your config", str.ValueString(), routingInstance)
}

func defaultResourceImportDontFindMessage(
	rsc resourceJunosNameable, id string,
) string {
	return fmt.Sprintf("don't find "+rsc.junosName()+" with id %q", id)
}

func defaultResourceImportDontFindIDStrMessage(
	rsc resourceJunosNameable, id, idformat string,
) string {
	return defaultResourceImportDontFindMessage(rsc, id) + fmt.Sprintf(" (id must be <%s>)", idformat)
}
