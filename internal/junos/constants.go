package junos

const (
	IDSeparator = "_-_"
	DefaultW    = "default"

	SetW     = "set"
	SetLS    = SetW + " " // set line start
	DeleteW  = "delete"
	DeleteLS = DeleteW + " " // delete line start

	CmdShowConfig          = "show configuration "
	PipeDisplaySet         = " | display set"
	PipeDisplaySetRelative = PipeDisplaySet + " relative"

	RoutingInstancesWS = "routing-instances " // routing-instances word + space

	RoutingOptionsWS = "routing-options "
	RibInet60WS      = "rib inet6.0 "

	EmptyW   = "empty"
	PermitW  = "permit"
	DiscardW = "discard"
	DisableW = "disable"

	St0Word = "st0"

	InetW  = "inet"
	Inet6W = "inet6"

	OspfV2 = "ospf"
	OspfV3 = "ospf3"

	CantReadValuesNotEnoughFields = "can't read values for %s in '%s': not enough fields"

	EnvHost                       = "JUNOS_HOST"
	EnvPort                       = "JUNOS_PORT"
	EnvUsername                   = "JUNOS_USERNAME"
	EnvPassword                   = "JUNOS_PASSWORD"
	EnvKeyPem                     = "JUNOS_KEYPEM"
	EnvKeyFile                    = "JUNOS_KEYFILE"
	EnvKeyPass                    = "JUNOS_KEYPASS"
	EnvGroupInterfaceDelete       = "JUNOS_GROUP_INTERFACE_DELETE"
	EnvSleepShort                 = "JUNOS_SLEEP_SHORT"
	EnvSleepLock                  = "JUNOS_SLEEP_LOCK"
	EnvCommitConfirmed            = "JUNOS_COMMIT_CONFIRMED"
	EnvCommitConfirmedWaitPercent = "JUNOS_COMMIT_CONFIRMED_WAIT_PERCENT"
	EnvSleepSSHClosed             = "JUNOS_SLEEP_SSH_CLOSED"
	EnvSSHTimeoutToEstablish      = "JUNOS_SSH_TIMEOUT_TO_ESTABLISH"
	EnvSSHRetryToEstablish        = "JUNOS_SSH_RETRY_TO_ESTABLISH"
	EnvFilePermission             = "JUNOS_FILE_PERMISSION"
	EnvLogPath                    = "JUNOS_LOG_PATH"
	EnvFakecreateSetfile          = "JUNOS_FAKECREATE_SETFILE"
	EnvFakeupdateAlso             = "JUNOS_FAKEUPDATE_ALSO"
	EnvFakedeleteAlso             = "JUNOS_FAKEDELETE_ALSO"

	DefaultInterfaceTestAcc        = "ge-0/0/3"
	DefaultInterfaceTestAcc2       = "ge-0/0/4"
	DefaultInterfaceSwitchTestAcc  = "xe-0/0/3"
	DefaultInterfaceSwitchTestAcc2 = "xe-0/0/4"
)
