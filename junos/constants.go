package junos

const (
	idSeparator = "_-_"
	defaultW    = "default"

	setW     = "set"
	setLS    = setW + " " // set line start
	deleteW  = "delete"
	deleteLS = deleteW + " " // delete line start

	cmdShowConfig = "show configuration "

	routingInstancesWS  = "routing-instances " // routing-instances word + space
	setRoutingInstances = setLS + routingInstancesWS
	delRoutingInstances = deleteLS + routingInstancesWS

	emptyW   = "empty"
	permitW  = "permit"
	discardW = "discard"
	disableW = "disable"

	st0Word = "st0"

	evpnW  = "evpn"
	inetW  = "inet"
	inet6W = "inet6"
	mplsW  = "mpls"
	ospfV2 = "ospf"
	ospfV3 = "ospf3"
)
