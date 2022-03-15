package junos

const (
	idSeparator = "_-_"
	defaultWord = "default"

	setWord      = "set"
	setLineStart = setWord + " "
	deleteWord   = "delete"
	deleteLS     = deleteWord + " " // delete line start

	cmdShowConfig = "show configuration "

	routingInstancesW   = "routing-instances "
	setRoutingInstances = setLineStart + routingInstancesW
	delRoutingInstances = deleteLS + routingInstancesW

	emptyWord  = "empty"
	permitWord = "permit"
	discardW   = "discard"
	disableW   = "disable"

	st0Word = "st0"

	evpnWord  = "evpn"
	inetWord  = "inet"
	inet6Word = "inet6"
	mplsWord  = "mpls"
	ospfV2    = "ospf"
	ospfV3    = "ospf3"
)
