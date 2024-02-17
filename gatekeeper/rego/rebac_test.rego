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
	"application_read": [
		{
			"action": "GET",
			"object": "/api/application",
		},
		{
			"action": "GET",
			"object": "/api/application/:id",
		},
	],
	"application_delete": [{
		"action": "DELETE",
		"object": "/api/application/:id",
	}],
}

test_any_ok if {
	allow with data.permissions as permissions
		with input.privileges as data.administrator.privileges
		with input.permission as {"action": "GET", "object": "/api/application"}
}

test_any_object_only_ok if {
	allow with data.permissions as permissions
		with input.privileges as data.readonly.privileges
		with input.permission as {"action": "GET", "object": "/api/application"}
}

test_any_object_only_ko if {
	not allow with data.permissions as permissions
		with input.privileges as data.readonly.privileges
		with input.permission as {"action": "DELETE", "object": "/api/application"}
}

test_any_action_ok if {
	allow with data.permissions as permissions
		with input.privileges as data.own.privileges
		with input.permission as {"action": "POST", "object": "/api/account/me"}
}

test_any_action_ko if {
	not allow with data.permissions as permissions
		with input.privileges as data.own.privileges
		with input.permission as {"action": "POST", "object": "/api/application"}
}

test_specific_matching_ok if {
	allow with data.permissions as permissions
		with input.privileges as data.multiple.privileges
		with input.permission as {"action": "DELETE", "object": "/api/application/:id"}
}

test_specific_matching_ko if {
	not allow with data.permissions as permissions
		with input.privileges as data.multiple.privileges
		with input.permission as {"action": "PUT", "object": "/api/application/:id"}
}
