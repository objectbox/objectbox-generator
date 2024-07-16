function (add_schema_file TARGET SCHEMA_FILE)
	set(OBX_SCHEMA_FILES ${OBX_SCHEMA_FILES} ${SCHEMA_FILE} PARENT_SCOPE)
endfunction()

function (generate_schema_bindings TARGET OUTPUT_DIR)
	set(OBX_GENERATED_CPP_FILES "")
	set(OBX_GENERATOR_COMMANDS "")

	foreach (SCHEMA_FILE ${OBX_SCHEMA_FILES})
		get_filename_component(SCHEMA_NAME ${SCHEMA_FILE} NAME_WE)
		get_filename_component(SCHEMA_DIR ${SCHEMA_FILE} DIRECTORY)

		set(OBX_GENERATED_CPP_FILE "${CMAKE_CURRENT_SOURCE_DIR}/${SCHEMA_DIR}/${SCHEMA_NAME}.obx.cpp")
		set(OBX_GENERATED_CPP_FILES ${OBX_GENERATED_CPP_FILES} ${OBX_GENERATED_CPP_FILE})

		set(OBX_GENERATOR_COMMAND COMMAND ${ObjectBoxGenerator_EXECUTABLE} -out ${OUTPUT_DIR} -cpp ${SCHEMA_FILE})
		set(OBX_GENERATOR_COMMANDS ${OBX_GENERATOR_COMMANDS} ${OBX_GENERATOR_COMMAND})
	endforeach()

	#message(STATUS OBX_GENERATOR_COMMANDS ${OBX_GENERATOR_COMMANDS})

	add_custom_command(
		OUTPUT
			__OBX_generate_schema_bindings_force_run
			${OBX_GENERATED_CPP_FILES}
		WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
		${OBX_GENERATOR_COMMANDS}
	)

	target_sources(${TARGET} PRIVATE
		__OBX_generate_schema_bindings_force_run
		${OBX_GENERATED_CPP_FILES}) # TODO Choose visibility externally ?
endfunction()
