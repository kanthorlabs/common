package kanthorlabs.gatekeeper

import rego.v1

permissions := {
	"administrator": [{
		"action": "*",
		"object": "*",
	}],
	"readonly": [
		{
			"action": "GET",
			"object": "*",
		},
		{
			"action": "HEAD",
			"object": "*",
		},
	],
	"own": [
		{
			"action": "*",
			"object": "/api/account/me",
		},
		{
			"action": "*",
			"object": "/api/account/password",
		},
	],
}

test_root_ok if {
	allow with data.permissions as permissions
		with input.privileges as data.administrator.privileges
}

test_readonly_ok if {
	allow with data.permissions as permissions
		with input.privileges as data.readonly.privileges
		with input.permission as {"action": "GET", "object": "/api/application"}
}

test_own_ok if {
	allow with data.permissions as permissions
		with input.privileges as data.own.privileges
		with input.permission as {"action": "POST", "object": "/api/account/me"}
}

test_ko if {
	not allow with data.permissions as permissions
		with input.privileges as data.own.privileges
		with input.permission as {"action": "POST", "object": "/api/application"}
}
