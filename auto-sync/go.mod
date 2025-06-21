module auto-sync

go 1.20

require (
	github.com/kardianos/service v1.2.1 // windows service helper
	github.com/zalando/go-keyring v0.2.1 // credential manager wrapper
	gopkg.in/yaml.v3 v3.0.1 // config file
)

require gopkg.in/natefinch/lumberjack.v2 v2.2.1

require (
	github.com/alessio/shellescape v1.4.1 // indirect
	github.com/danieljoos/wincred v1.1.0 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	golang.org/x/sys v0.0.0-20220429233432-b5fbb4746d32 // indirect
)
