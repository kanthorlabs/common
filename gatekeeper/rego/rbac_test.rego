package kanthorlabs.gatekeeper

import rego.v1

test_any_ok if {
	allow with data.permissions as data.data.definitions
		with input.privileges as data.input.administrator.privileges
		with input.permission as {"action": "GET", "object": "/api/application"}
}

test_any_object_only_ok if {
	allow with data.permissions as data.data.definitions
		with input.privileges as data.input.readonly.privileges
		with input.permission as {"action": "GET", "object": "/api/application"}
}

test_any_object_only_ko if {
	not allow with data.permissions as data.data.definitions
		with input.privileges as data.input.readonly.privileges
		with input.permission as {"action": "DELETE", "object": "/api/application"}
}

test_any_action_ok if {
	allow with data.permissions as data.data.definitions
		with input.privileges as data.input.own.privileges
		with input.permission as {"action": "POST", "object": "/api/account/me"}
}

test_any_action_ko if {
	not allow with data.permissions as data.data.definitions
		with input.privileges as data.input.own.privileges
		with input.permission as {"action": "POST", "object": "/api/application"}
}

test_specific_matching_ok if {
	allow with data.permissions as data.data.definitions
		with input.privileges as data.input.multiple.privileges
		with input.permission as {"action": "DELETE", "object": "/api/application/:id"}
}

test_specific_matching_ko if {
	not allow with data.permissions as data.data.definitions
		with input.privileges as data.input.multiple.privileges
		with input.permission as {"action": "PUT", "object": "/api/application/:id"}
}
