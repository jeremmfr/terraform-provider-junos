package junos

func SyslogFacilities() []string {
	return []string{
		"authorization",
		"daemon",
		"ftp",
		"kernel",
		"user",
		"local0",
		"local1",
		"local2",
		"local3",
		"local4",
		"local5",
		"local6",
		"local7",
	}
}

func SyslogSeverity() []string {
	return []string{
		"alert",
		"any",
		"critical",
		"emergency",
		"error",
		"info",
		"none",
		"notice",
		"warning",
	}
}
