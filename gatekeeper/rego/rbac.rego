package kanthorlabs.gatekeeper

import rego.v1

any := "*"

default allow := false

allow if {
	some privilege in input.privileges

	# make sure the role is exist
	data.permissions[privilege.role]

	some permission in data.permissions[privilege.role]

	# matching
	pass_action(permission.action)
	pass_object(permission.object)
}

pass_action(action) if {
	action == any
}

pass_action(action) if {
	action == input.permission.action
}

pass_object(object) if {
	object == any
}

pass_object(object) if {
	object == input.permission.object
}
