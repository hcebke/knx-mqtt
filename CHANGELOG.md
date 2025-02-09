# Version 1.2
Replaces / with _ in GroupAddress Names to avoid unwanted MQTT topic separations.

# Version 1.1
Added configuration parameter emitValueAsString. When set to `false`, values will preserve its type when sent as `json`. If set to `true`, behaviour is compatible with previous version.