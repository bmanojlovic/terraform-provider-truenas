#!/usr/bin/env python3
"""
Generate Terraform provider from TrueNAS native method specifications.
Replaces OpenAPI-based generator with native core.get_methods approach.
"""
import json
import sys
from pathlib import Path

# Load templates
TEMPLATE_DIR = Path(__file__).parent / "templates"
RESOURCE_TEMPLATE = (TEMPLATE_DIR / "resource.go.tmpl").read_text()
RESOURCE_UPDATE_ONLY_TEMPLATE = (
    TEMPLATE_DIR / "resource_update_only.go.tmpl"
).read_text()
RESOURCE_WITH_JSON_TEMPLATE = (TEMPLATE_DIR / "resource_with_json.go.tmpl").read_text()
RESOURCE_VM_DEVICE_TEMPLATE = (TEMPLATE_DIR / "resource_vm_device.go.tmpl").read_text()
RESOURCE_DOC_TEMPLATE = (TEMPLATE_DIR / "resource_doc.md.tmpl").read_text()
DATASOURCE_TEMPLATE = (TEMPLATE_DIR / "datasource.go.tmpl").read_text()
DATASOURCE_DOC_TEMPLATE = (TEMPLATE_DIR / "datasource_doc.md.tmpl").read_text()
DATASOURCE_QUERY_TEMPLATE = (TEMPLATE_DIR / "datasource_query.go.tmpl").read_text()
DATASOURCE_QUERY_DOC_TEMPLATE = (
    TEMPLATE_DIR / "datasource_query_doc.md.tmpl"
).read_text()

TYPE_MAP = {
    "string": "String",
    "integer": "Int64",
    "number": "Float64",
    "boolean": "Bool",
    "array": "List",
    "object": "String",  # object as JSON string
}


def find_latest_spec():
    specs = list(Path(".").glob("truenas-methods-*.json"))
    if not specs:
        print("ERROR: No spec file found. Run: make fetch-spec", file=sys.stderr)
        sys.exit(1)
    return sorted(specs)[-1]


def load_spec():
    spec_file = find_latest_spec()
    print(f"Using: {spec_file}", file=sys.stderr)
    with open(spec_file) as f:
        data = json.load(f)
    return data.get("methods", {}), data.get("_metadata", {})


def get_tf_type(prop_schema):
    """Convert JSON schema to Terraform type."""
    if isinstance(prop_schema, list):
        prop_schema = prop_schema[0] if prop_schema else {}
    if not isinstance(prop_schema, dict):
        return "String"

    # Handle anyOf/oneOf - check for integer type
    if "anyOf" in prop_schema:
        for variant in prop_schema["anyOf"]:
            if variant.get("type") == "integer":
                return "Int64"
        return "String"

    json_type = prop_schema.get("type")
    if "oneOf" in prop_schema or "discriminator" in prop_schema:
        return "String"  # Complex objects as JSON strings
    return TYPE_MAP.get(json_type, "String")


def generate_schema_attrs(
    properties, required, has_start=False, create_only_fields=None
):
    """Generate schema attributes for template."""
    if create_only_fields is None:
        create_only_fields = set()
    attrs = []

    # For datasources, id is always String (user input)
    # For resources, id is Computed and matches schema type
    if not has_start and not required:  # datasource mode
        attrs.append(
            '\t\t\t"id": schema.StringAttribute{Required: true, Description: "Resource ID"},'
        )
    elif "id" not in properties:
        attrs.append(
            '\t\t\t"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},'
        )

    if has_start:
        attrs.append(
            '\t\t\t"start_on_create": schema.BoolAttribute{Optional: true, Description: "Start the resource immediately after creation (default: true)"},'
        )

    for name, prop in properties.items():
        # Skip id if already added, or if datasource mode (id handled separately)
        if name == "id":
            if (
                not has_start and not required
            ):  # datasource - skip, already added as String
                continue
            # For resources, add as Computed with actual type
            tf_type = get_tf_type(prop)
            attrs.append(
                f'\t\t\t"id": schema.{tf_type}Attribute{{Computed: true, Description: "Resource ID"}},'
            )
            continue
        # Skip reserved names
        if name in ["provider"]:
            continue
        if isinstance(prop, list):
            prop = prop[0] if prop else {}

        tf_type = get_tf_type(prop)
        is_req = name in required
        desc = (
            prop.get("description", "")[:100].replace('"', '\\"').replace("\n", " ")
            if isinstance(prop, dict)
            else ""
        )

        # Fix invalid attribute names
        attr_name = name.lower() if name != "CSR" else "csr"

        attrs.append(f'\t\t\t"{attr_name}": schema.{tf_type}Attribute{{')

        # For datasources (no required list), all attributes except id are Computed
        if not has_start and not required:  # datasource mode
            attrs.append(f"\t\t\t\tComputed: true,")
        else:  # resource mode
            attrs.append(f"\t\t\t\tRequired: {str(is_req).lower()},")
            attrs.append(f"\t\t\t\tOptional: {str(not is_req).lower()},")

        if tf_type == "List":
            attrs.append(f"\t\t\t\tElementType: types.StringType,")
        attrs.append(f'\t\t\t\tDescription: "{desc}",')

        # Add RequiresReplace for create-only fields (except 'name' which may support rename)
        if name in create_only_fields and name != "name":
            if tf_type == "String":
                attrs.append(
                    f"\t\t\t\tPlanModifiers: []planmodifier.String{{stringplanmodifier.RequiresReplace()}},"
                )
            elif tf_type == "Int64":
                attrs.append(
                    f"\t\t\t\tPlanModifiers: []planmodifier.Int64{{int64planmodifier.RequiresReplace()}},"
                )
            elif tf_type == "Bool":
                attrs.append(
                    f"\t\t\t\tPlanModifiers: []planmodifier.Bool{{boolplanmodifier.RequiresReplace()}},"
                )

        attrs.append("\t\t\t},")

    return "\n".join(attrs)


def generate_fields(properties, has_start=False):
    """Generate struct fields for template."""
    fields = []

    # For datasources, ID is always String (user input)
    # For resources, ID matches schema type or String if not in properties
    is_datasource = not has_start and "id" in properties
    if is_datasource:
        fields.append('\tID types.String `tfsdk:"id"`')
    elif "id" not in properties:
        fields.append('\tID types.String `tfsdk:"id"`')

    if has_start:
        fields.append('\tStartOnCreate types.Bool `tfsdk:"start_on_create"`')

    for name, prop in properties.items():
        # Skip reserved names and fix invalid names
        if name in ["provider"]:
            continue
        field_name = name.title().replace("_", "")
        if name == "CSR":
            field_name = "Csr"
        elif name == "id":
            if is_datasource:
                continue  # Already added as String
            field_name = "ID"
        tf_type = get_tf_type(prop)
        fields.append(f'\t{field_name} types.{tf_type} `tfsdk:"{name.lower()}"`')
    return "\n".join(fields)


def generate_create_params(properties):
    """Generate parameter building code for Create method."""
    lines = []
    for prop_name, prop_schema in properties.items():
        # Skip reserved names and id (id is not sent in create/update params)
        if prop_name in ["provider", "id"]:
            continue
        field_name = prop_name.title().replace("_", "")
        if prop_name == "CSR":
            field_name = "Csr"
        tf_type = get_tf_type(prop_schema)

        lines.append(f"\tif !data.{field_name}.IsNull() {{")

        if tf_type == "Bool":
            lines.append(f'\t\tparams["{prop_name}"] = data.{field_name}.ValueBool()')
        elif tf_type == "Int64":
            lines.append(f'\t\tparams["{prop_name}"] = data.{field_name}.ValueInt64()')
        elif tf_type == "Float64":
            lines.append(
                f'\t\tparams["{prop_name}"] = data.{field_name}.ValueFloat64()'
            )
        elif tf_type == "List":
            lines.append(f"\t\tvar {prop_name}List []string")
            lines.append(
                f"\t\tdata.{field_name}.ElementsAs(ctx, &{prop_name}List, false)"
            )
            lines.append(f'\t\tparams["{prop_name}"] = {prop_name}List')
        else:
            # Check if this is a complex object that needs JSON parsing
            if isinstance(prop_schema, dict) and (
                "oneOf" in prop_schema
                or "discriminator" in prop_schema
                or prop_schema.get("type") == "object"
            ):
                lines.append(f"\t\tvar {prop_name}Obj map[string]interface{{}}")
                lines.append(
                    f"\t\tif err := json.Unmarshal([]byte(data.{field_name}.ValueString()), &{prop_name}Obj); err != nil {{"
                )
                lines.append(
                    f'\t\t\tresp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse {prop_name}: %s", err))'
                )
                lines.append(f"\t\t\treturn")
                lines.append(f"\t\t}}")
                lines.append(f'\t\tparams["{prop_name}"] = {prop_name}Obj')
            else:
                lines.append(
                    f'\t\tparams["{prop_name}"] = data.{field_name}.ValueString()'
                )

        lines.append("\t}")
    return "\n".join(lines)


def generate_read_mapping(properties, skip_id_for_datasource=False):
    """Generate code to map API response to Terraform state."""
    lines = []
    for prop_name, prop_schema in properties.items():
        # Skip reserved names
        if prop_name in ["provider"]:
            continue
        # Skip id for datasources (already set as String from input)
        if prop_name == "id" and skip_id_for_datasource:
            continue
        field_name = prop_name.title().replace("_", "")
        if prop_name == "CSR":
            field_name = "Csr"
        elif prop_name == "id":
            field_name = "ID"
        tf_type = get_tf_type(prop_schema)

        lines.append(f'\t\tif v, ok := resultMap["{prop_name}"]; ok && v != nil {{')

        if tf_type == "Bool":
            lines.append(
                f"\t\t\tif bv, ok := v.(bool); ok {{ data.{field_name} = types.BoolValue(bv) }}"
            )
        elif tf_type == "Int64":
            lines.append(
                f"\t\t\tif fv, ok := v.(float64); ok {{ data.{field_name} = types.Int64Value(int64(fv)) }}"
            )
        elif tf_type == "Float64":
            lines.append(
                f"\t\t\tif fv, ok := v.(float64); ok {{ data.{field_name} = types.Float64Value(fv) }}"
            )
        elif tf_type == "List":
            lines.append(f"\t\t\tif arr, ok := v.([]interface{{}}); ok {{")
            lines.append(f"\t\t\t\tstrVals := make([]attr.Value, len(arr))")
            lines.append(
                f'\t\t\t\tfor i, item := range arr {{ strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }}'
            )
            lines.append(
                f"\t\t\t\tdata.{field_name}, _ = types.ListValue(types.StringType, strVals)"
            )
            lines.append(f"\t\t\t}}")
        else:
            lines.append(
                f'\t\t\tdata.{field_name} = types.StringValue(fmt.Sprintf("%v", v))'
            )

        lines.append("\t\t}")
    return "\n".join(lines)


def has_complex_objects(properties):
    """Check if any property needs JSON parsing."""
    for prop_name, prop_schema in properties.items():
        # Skip reserved names
        if prop_name in ["provider"]:
            continue
        if isinstance(prop_schema, dict) and (
            "oneOf" in prop_schema
            or "discriminator" in prop_schema
            or prop_schema.get("type") == "object"
        ):
            return True
    return False


def generate_resource(base_name, methods_dict):
    """Generate resource file from method specs."""
    create_method = f"{base_name}.create"
    update_method = f"{base_name}.update"
    delete_method = f"{base_name}.delete"
    read_method = f"{base_name}.get_instance"

    # Check for lifecycle actions
    has_start = f"{base_name}.start" in methods_dict
    has_stop = f"{base_name}.stop" in methods_dict

    # Check if methods are jobs
    create_spec = methods_dict.get(create_method, {})
    update_spec = methods_dict.get(update_method, {})
    delete_spec = methods_dict.get(delete_method, {})

    create_is_job = create_spec.get("job", False)
    update_is_job = update_spec.get("job", False)
    delete_is_job = delete_spec.get("job", False)

    # Detect ID type from update or delete method (first parameter)
    id_is_string = False
    for spec in [update_spec, delete_spec]:
        if spec and spec.get("accepts"):
            first_param = spec["accepts"][0]
            if first_param.get("type") == "string":
                id_is_string = True
                break

    # Get schema from create or update
    method_spec = create_spec or update_spec
    if not method_spec:
        return None

    accepts = method_spec.get("accepts", [])
    if not accepts:
        return None

    schema = accepts[0] if isinstance(accepts, list) else accepts

    # Handle anyOf schemas by merging all variants
    if "anyOf" in schema:
        merged_props = {}
        for variant in schema["anyOf"]:
            merged_props.update(variant.get("properties", {}))
        properties = merged_props
        # Only mark fields as required if they're required in ALL variants
        all_required = [set(variant.get("required", [])) for variant in schema["anyOf"]]
        required = list(set.intersection(*all_required)) if all_required else []
    else:
        properties = schema.get("properties", {})
        required = schema.get("required", [])

    if not properties:
        return None

    # Get update schema and merge with create (union approach - standard Terraform pattern)
    update_properties = {}
    create_only_fields = set()
    if update_spec and update_spec.get("accepts") and len(update_spec["accepts"]) >= 2:
        # Update has [id, data] format - get data schema (second parameter)
        update_schema = update_spec["accepts"][1]
        if "anyOf" in update_schema:
            for variant in update_schema["anyOf"]:
                update_properties.update(variant.get("properties", {}))
        elif update_schema.get("properties"):
            update_properties = update_schema["properties"]

        # Remove 'id' from update_properties - it's the resource identifier, not a field
        update_properties = {k: v for k, v in update_properties.items() if k != "id"}

        # Identify create-only fields (in create but not in update) - these need ForceNew
        create_only_fields = set(properties.keys()) - set(update_properties.keys())

        # Merge: union of create and update properties
        all_properties = {**properties, **update_properties}
        properties = all_properties
    else:
        # No update method or schema
        create_only_fields = set(properties.keys())
        update_properties = {}

    # Generate code
    resource_name = base_name.replace(".", "_").title().replace("_", "")
    tf_name = base_name.replace(".", "_")
    api_name = base_name

    # Choose Call or CallWithJob based on job flag
    create_call = "CallWithJob" if create_is_job else "Call"
    update_call = "CallWithJob" if update_is_job else "Call"
    delete_call = "CallWithJob" if delete_is_job else "Call"

    # Generate ID handling code based on type
    if id_is_string:
        id_read_code = "\tid = data.ID.ValueString()"
        id_update_code = "\tid = state.ID.ValueString()"
        id_delete_code = "\tid = data.ID.ValueString()"
    else:
        id_read_code = """	id, err = strconv.Atoi(data.ID.ValueString())
	if err != nil {{
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}}"""
        id_update_code = """	id, err = strconv.Atoi(state.ID.ValueString())
	if err != nil {{
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}}"""
        id_delete_code = """	id, err = strconv.Atoi(data.ID.ValueString())
	if err != nil {{
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}}"""

    # Generate lifecycle code if has start action
    lifecycle_code = ""
    if has_start:
        lifecycle_code = f"""
\t// Handle lifecycle action - start on create if requested
\tstartOnCreate := true  // default when not specified
\tif !data.StartOnCreate.IsNull() {{
\t\tstartOnCreate = data.StartOnCreate.ValueBool()
\t}}
\tif startOnCreate {{
\t\tvmID, err := strconv.Atoi(data.ID.ValueString())
\t\tif err != nil {{
\t\t\tresp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
\t\t\treturn
\t\t}}
\t\t_, err = r.client.Call("{api_name}.start", vmID)
\t\tif err != nil {{
\t\t\tresp.Diagnostics.AddWarning("Start Failed", fmt.Sprintf("Resource created but failed to start: %s", err.Error()))
\t\t}}
\t}}"""

    # Determine if strconv import is needed (for int IDs or lifecycle code)
    extra_imports = '\t"strconv"' if (not id_is_string or has_start) else ""

    # Check if any List fields exist
    has_list = any(get_tf_type(p) == "List" for p in properties.values())
    # Check if any complex objects need JSON parsing
    has_json = has_complex_objects(properties)

    # Check which plan modifiers are needed for create-only fields
    needs_string_planmod = False
    needs_int64_planmod = False
    needs_bool_planmod = False
    for field in create_only_fields:
        if field != "name" and field in properties:
            tf_type = get_tf_type(properties[field])
            if tf_type == "String":
                needs_string_planmod = True
            elif tf_type == "Int64":
                needs_int64_planmod = True
            elif tf_type == "Bool":
                needs_bool_planmod = True

    if has_list:
        extra_imports += '\n\t"github.com/hashicorp/terraform-plugin-framework/attr"'
    if has_json:
        extra_imports += '\n\t"encoding/json"'
    if has_stop:
        extra_imports += '\n\t"time"'
    if needs_string_planmod or needs_int64_planmod or needs_bool_planmod:
        extra_imports += '\n\t"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"'
    if needs_string_planmod:
        extra_imports += '\n\t"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"'
    if needs_int64_planmod:
        extra_imports += '\n\t"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"'
    if needs_bool_planmod:
        extra_imports += '\n\t"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"'

    # Get description from method spec
    description = method_spec.get("description", f"TrueNAS {tf_name} resource")
    if description:
        description = description.split("\n")[0][:200].replace('"', '\\"')

    # Generate pre-delete code if has stop action
    predelete_code = ""
    if has_stop:
        predelete_code = f"""
\t// Stop VM before deletion if running
\tvmID, err := strconv.Atoi(data.ID.ValueString())
\tif err != nil {{
\t\tresp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
\t\treturn
\t}}
\t_, _ = r.client.Call("{api_name}.stop", vmID)  // Ignore errors - VM might already be stopped
\ttime.Sleep(2 * time.Second)  // Wait for VM to stop
"""

    # Use special template for vm.device
    if api_name == "vm.device":
        code = RESOURCE_VM_DEVICE_TEMPLATE.format(
            resource_name=resource_name,
            name=tf_name,
            api_name=api_name,
            fields=generate_fields(properties, has_start),
            schema_attrs=generate_schema_attrs(
                properties, required, has_start, create_only_fields
            ),
            create_params=generate_create_params(properties),
            read_mapping=generate_read_mapping(properties),
            lifecycle_code=lifecycle_code,
        )
    else:
        code = RESOURCE_TEMPLATE.format(
            resource_name=resource_name,
            name=tf_name,
            api_name=api_name,
            description=description,
            fields=generate_fields(properties, has_start),
            schema_attrs=generate_schema_attrs(
                properties, required, has_start, create_only_fields
            ),
            create_params=generate_create_params(properties),
            update_params=generate_create_params(
                update_properties if update_properties else properties
            ),
            read_mapping=generate_read_mapping(properties),
            lifecycle_code=lifecycle_code,
            predelete_code=predelete_code,
            id_read_code=id_read_code,
            id_update_code=id_update_code,
            id_delete_code=id_delete_code,
            extra_imports=extra_imports,
            create_call=create_call,
            update_call=update_call,
            delete_call=delete_call,
        )

    return code


def generate_attr_types(properties):
    """Generate AttrTypes map for ListValueFrom"""
    lines = []
    for name, prop in sorted(properties.items()):
        if name in ["provider"]:
            continue
        tf_type = get_tf_type(prop)
        # Skip List types - too complex for query datasources
        if tf_type == "List":
            continue
        # ID is always String in datasources
        if name == "id":
            tf_type = "String"
        attr_name = name.lower() if name != "CSR" else "csr"

        if tf_type == "String":
            lines.append(f'\t\t\t"{attr_name}": types.StringType,')
        elif tf_type == "Int64":
            lines.append(f'\t\t\t"{attr_name}": types.Int64Type,')
        elif tf_type == "Bool":
            lines.append(f'\t\t\t"{attr_name}": types.BoolType,')
        elif tf_type == "Float64":
            lines.append(f'\t\t\t"{attr_name}": types.Float64Type,')
    return "\n".join(lines)


def generate_query_datasource(base_name, methods_dict):
    """Generate query data source for listing multiple resources."""
    query_method = f"{base_name}.query"

    query_spec = methods_dict.get(query_method, {})
    if not query_spec:
        return None

    # Get return schema
    returns = query_spec.get("returns", [])
    if not returns:
        return None

    # Query returns array of items - handle anyOf wrapper
    schema = returns[0] if isinstance(returns, list) else returns

    # Handle anyOf - look for array type
    if "anyOf" in schema:
        for variant in schema["anyOf"]:
            if isinstance(variant, dict) and variant.get("type") == "array":
                schema = variant
                break

    if schema.get("type") != "array":
        return None

    items_schema = schema.get("items", {})
    if isinstance(items_schema, list):
        items_schema = items_schema[0] if items_schema else {}

    properties = items_schema.get("properties", {})
    if not properties:
        return None

    # Filter out List types - too complex for query datasources
    filtered_properties = {
        k: v for k, v in properties.items() if get_tf_type(v) != "List"
    }
    if not filtered_properties:
        return None

    # Generate code
    resource_name = base_name.replace(".", "_").title().replace("_", "") + "s"  # Plural
    tf_name = base_name.replace(".", "_") + "s"  # Plural
    api_name = base_name
    description = query_spec.get("description") or f"Query {tf_name} resources"
    description = description.split("\n")[0][:200].replace('"', '\\"')

    # Generate read mapping for items
    read_mapping_lines = []
    for prop_name, prop_schema in filtered_properties.items():
        if prop_name in ["provider"]:
            continue
        field_name = prop_name.title().replace("_", "")
        if prop_name == "CSR":
            field_name = "Csr"
        elif prop_name == "id":
            field_name = "ID"
        tf_type = get_tf_type(prop_schema)

        read_mapping_lines.append(
            f'\t\tif v, ok := resultMap["{prop_name}"]; ok && v != nil {{'
        )
        if tf_type == "Bool":
            read_mapping_lines.append(
                f"\t\t\tif bv, ok := v.(bool); ok {{ itemModel.{field_name} = types.BoolValue(bv) }}"
            )
        elif (
            tf_type == "Int64" and prop_name != "id"
        ):  # ID is always String in datasources
            read_mapping_lines.append(
                f"\t\t\tif fv, ok := v.(float64); ok {{ itemModel.{field_name} = types.Int64Value(int64(fv)) }}"
            )
        elif tf_type == "Float64":
            read_mapping_lines.append(
                f"\t\t\tif fv, ok := v.(float64); ok {{ itemModel.{field_name} = types.Float64Value(fv) }}"
            )
        elif tf_type == "List":
            read_mapping_lines.append(f"\t\t\t// Skip complex list types for now")
        else:
            read_mapping_lines.append(
                f'\t\t\titemModel.{field_name} = types.StringValue(fmt.Sprintf("%v", v))'
            )
        read_mapping_lines.append("\t\t}")

    code = DATASOURCE_QUERY_TEMPLATE.format(
        resource_name=resource_name,
        name=tf_name,
        api_name=api_name,
        description=description,
        fields=generate_fields(filtered_properties, False),
        schema_attrs=generate_schema_attrs(filtered_properties, [], False),
        read_mapping="\n".join(read_mapping_lines),
        attr_types=generate_attr_types(filtered_properties),
    )

    return code


def generate_datasource(base_name, methods_dict):
    """Generate data source file from method specs."""
    get_method = f"{base_name}.get_instance"

    get_spec = methods_dict.get(get_method, {})
    if not get_spec:
        return None

    # Get return schema
    returns = get_spec.get("returns", [])
    if not returns:
        return None

    schema = returns[0] if isinstance(returns, list) else returns
    properties = schema.get("properties", {})

    if not properties:
        return None

    # Check if any List fields exist
    has_list = any(get_tf_type(p) == "List" for p in properties.values())

    extra_imports = ""
    if has_list:
        extra_imports = '\n\t"github.com/hashicorp/terraform-plugin-framework/attr"'

    # Determine ID type and parameter format from schema
    id_type = get_tf_type(properties.get("id", {"type": "string"}))
    if id_type == "Int64":
        # ID in schema is Int64, but datasource input is String - convert
        id_param = (
            "func() int { id, _ := strconv.Atoi(data.ID.ValueString()); return id }()"
        )
        extra_imports += '\n\t"strconv"'
    else:
        id_param = "data.ID.ValueString()"

    # Generate code
    resource_name = base_name.replace(".", "_").title().replace("_", "")
    tf_name = base_name.replace(".", "_")
    api_name = base_name
    description = get_spec.get("description") or f"Retrieves TrueNAS {tf_name} data"
    description = description.split("\n")[0][:200].replace('"', '\\"')

    code = DATASOURCE_TEMPLATE.format(
        resource_name=resource_name,
        name=tf_name,
        api_name=api_name,
        description=description,
        fields=generate_fields(properties, False),
        schema_attrs=generate_schema_attrs(properties, [], False),
        read_mapping=generate_read_mapping(properties, skip_id_for_datasource=True),
        extra_imports=extra_imports,
        id_param=id_param,
    )

    return code


def generate_datasource_docs(base_name, properties, description):
    """Generate data source documentation"""
    tf_name = base_name.replace(".", "_")

    # Build attributes list
    attrs = []
    for name, prop in sorted(properties.items()):
        if name in ["id"]:
            continue
        tf_type = get_tf_type(prop)
        desc = prop.get("description", "") if isinstance(prop, dict) else ""
        desc = desc.replace("\n", " ").strip()[:200]
        attrs.append(f"- `{name}` ({tf_type}) - {desc}")

    doc = DATASOURCE_DOC_TEMPLATE.format(
        resource_type=tf_name,
        description=description,
        name=tf_name,
        attrs=chr(10).join(attrs) if attrs else "- None",
    )

    docs_dir = Path("docs/data-sources")
    docs_dir.mkdir(parents=True, exist_ok=True)
    (docs_dir / f"{tf_name}.md").write_text(doc)


def generate_query_datasource_docs(base_name, properties, description):
    """Generate query data source documentation"""
    tf_name = base_name.replace(".", "_") + "s"  # Plural

    # Build attributes list
    attrs = []
    for name, prop in sorted(properties.items()):
        if name in ["id"]:
            continue
        tf_type = get_tf_type(prop)
        desc = prop.get("description", "") if isinstance(prop, dict) else ""
        desc = desc.replace("\n", " ").strip()[:200]
        attrs.append(f"- `{name}` ({tf_type}) - {desc}")

    doc = DATASOURCE_QUERY_DOC_TEMPLATE.format(
        resource_type=tf_name,
        description=description,
        name=base_name.replace(".", "_"),
        attrs=chr(10).join(attrs) if attrs else "- None",
    )

    docs_dir = Path("docs/data-sources")
    docs_dir.mkdir(parents=True, exist_ok=True)
    (docs_dir / f"{tf_name}.md").write_text(doc)


def generate_resource_docs(
    base_name, properties, required, description, methods_dict, anyof_variants=None
):
    """Generate Terraform documentation markdown"""
    tf_name = base_name.replace(".", "_")
    has_start = f"{base_name}.start" in methods_dict

    # Detect discriminator field for anyOf schemas
    discriminator_field = None
    variant_info = {}
    discriminator_all_values = []

    if anyof_variants:
        # Find discriminator field (usually 'type')
        for field_name in ["type", "kind", "variant"]:
            if field_name in properties:
                prop = properties[field_name]
                if isinstance(prop, dict) and "enum" in prop:
                    discriminator_field = field_name
                    break

        # Build variant information and collect all discriminator values
        for variant in anyof_variants:
            variant_props = variant.get("properties", {})
            variant_required = set(variant.get("required", []))

            # Identify variant by discriminator value
            variant_name = None
            if discriminator_field and discriminator_field in variant_props:
                disc_prop = variant_props[discriminator_field]
                if isinstance(disc_prop, dict):
                    if "enum" in disc_prop and disc_prop["enum"]:
                        variant_name = disc_prop["enum"][0]
                        discriminator_all_values.extend(disc_prop["enum"])
                    elif "default" in disc_prop:
                        variant_name = disc_prop["default"]
                        discriminator_all_values.append(disc_prop["default"])

            if variant_name:
                variant_info[variant_name] = {
                    "properties": set(variant_props.keys()),
                    "required": variant_required,
                }

    # Build example with required fields
    example_lines = []
    for name in sorted(required):
        if name in properties and name not in ["uuid", "id"]:
            prop = properties[name]
            tf_type = get_tf_type(prop)
            if tf_type == "String":
                example_lines.append(f'  {name} = "example-value"')
            elif tf_type == "Int64":
                example_lines.append(f"  {name} = 1")
            elif tf_type == "Bool":
                example_lines.append(f"  {name} = true")
            elif tf_type == "Float64":
                example_lines.append(f"  {name} = 1.0")
            elif tf_type == "List":
                example_lines.append(f'  {name} = ["item1"]')

    if has_start and len(example_lines) < 8:
        example_lines.append("  start_on_create = true")

    # Build schema documentation
    required_args = []
    optional_args = []

    if has_start:
        optional_args.append(
            "- `start_on_create` (Bool) - Start the resource immediately after creation. Default: `true`"
        )

    for name, prop in sorted(properties.items()):
        if name in ["provider", "uuid", "id"]:
            continue
        tf_type = get_tf_type(prop)
        desc = prop.get("description", "") if isinstance(prop, dict) else ""
        desc = desc.replace("\n", " ").strip()[:200]

        # Add default value if present
        if isinstance(prop, dict) and "default" in prop:
            desc += f" Default: `{prop['default']}`"

        # Add enum values if present
        if isinstance(prop, dict) and "enum" in prop:
            # For discriminator field in anyOf schemas, show all variant values
            if name == discriminator_field and discriminator_all_values:
                enum_vals = ", ".join(
                    f"`{v}`" for v in sorted(set(discriminator_all_values))[:10]
                )
            else:
                enum_vals = ", ".join(f"`{v}`" for v in prop["enum"][:10])
            desc += f" Valid values: {enum_vals}"

        # Add variant applicability if this is an anyOf schema
        if variant_info and name != discriminator_field:
            applicable_variants = [
                v for v, info in variant_info.items() if name in info["properties"]
            ]
            if applicable_variants and len(applicable_variants) < len(variant_info):
                variant_list = ", ".join(f"`{v}`" for v in applicable_variants)
                desc += f" **Applies to:** {variant_list}"

        arg_line = f"- `{name}` ({tf_type}) - {desc}"

        if name in required:
            required_args.append(arg_line)
        else:
            optional_args.append(arg_line)

    # Build variant examples if anyOf present
    variant_examples = ""
    if variant_info and discriminator_field:
        variant_examples = "\n\n## Variants\n\n"
        variant_examples += f"This resource has **{len(variant_info)} variants** controlled by the `{discriminator_field}` field. "
        variant_examples += "Choose the appropriate variant for your use case:\n\n"

        for variant_name, info in sorted(variant_info.items()):
            variant_examples += f"### {variant_name}\n\n"

            # Build example for this variant
            variant_example_lines = []
            variant_required = sorted(info["required"] - {"provider", "uuid", "id"})

            # Add discriminator field first
            variant_example_lines.append(f'  {discriminator_field} = "{variant_name}"')

            # Add other required fields
            for req_field in variant_required:
                if req_field == discriminator_field:
                    continue
                if req_field in properties:
                    prop = properties[req_field]
                    tf_type = get_tf_type(prop)
                    if tf_type == "String":
                        variant_example_lines.append(f'  {req_field} = "value"')
                    elif tf_type == "Int64":
                        variant_example_lines.append(f"  {req_field} = 1024")
                    elif tf_type == "Bool":
                        variant_example_lines.append(f"  {req_field} = true")

            # Show example
            variant_examples += "```terraform\n"
            variant_examples += f'resource "truenas_{tf_name}" "example" {{\n'
            variant_examples += "\n".join(variant_example_lines)
            variant_examples += "\n}\n```\n\n"

            # Show required fields
            if variant_required:
                variant_examples += (
                    "**Required fields:** "
                    + ", ".join(f"`{r}`" for r in variant_required)
                    + "\n\n"
                )

            # Show variant-specific optional fields (not in all variants)
            variant_props = info["properties"] - {"provider", "uuid", "id"}
            all_props = set()
            for v_info in variant_info.values():
                all_props.update(v_info["properties"])
            variant_specific = sorted(
                (
                    variant_props
                    - all_props.intersection(
                        *[v["properties"] for v in variant_info.values()]
                    )
                )
                - info["required"]
            )

            if variant_specific:
                variant_examples += "**Key optional fields:** " + ", ".join(
                    f"`{p}`" for p in variant_specific[:8]
                )
                if len(variant_specific) > 8:
                    variant_examples += f" (and {len(variant_specific) - 8} more)"
                variant_examples += "\n\n"

    # Build generic example section (skip for anyOf resources)
    generic_example = ""
    if not variant_info:
        generic_example = "\n## Example Usage\n\n```terraform\n"
        generic_example += f'resource "truenas_{tf_name}" "example" {{\n'
        generic_example += (
            chr(10).join(example_lines)
            if example_lines
            else "  # Configure required attributes"
        )
        generic_example += "\n}\n```\n"

    doc = RESOURCE_DOC_TEMPLATE.format(
        resource_type=tf_name,
        description=description,
        required_args=chr(10).join(required_args) if required_args else "- None",
        optional_args=chr(10).join(optional_args) if optional_args else "- None",
        variant_examples=variant_examples,
        generic_example=generic_example,
    )

    docs_dir = Path("docs/resources")
    docs_dir.mkdir(parents=True, exist_ok=True)
    (docs_dir / f"{tf_name}.md").write_text(doc)


def generate_provider(resources, datasources):
    """Generate provider.go from template"""
    with open("templates/provider.go.tmpl", "r") as f:
        template = f.read()

    # Build resource list
    resource_funcs = [
        f"New{r.replace('.', '_').title().replace('_', '')}Resource" for r in resources
    ]
    resource_list = ",\n\t\t".join(resource_funcs)

    # Build datasource list
    datasource_funcs = [
        f"New{d.replace('.', '_').title().replace('_', '')}DataSource"
        for d in datasources
    ]
    datasource_list = ",\n\t\t".join(datasource_funcs)
    if datasource_list:
        datasource_list += ","  # Add trailing comma

    # Replace template variables
    code = template.replace("{{resource_list}}", resource_list)
    code = code.replace("{{datasource_list}}", datasource_list)

    with open("internal/provider/provider.go", "w") as f:
        f.write(code)

    print("✅ Generated provider.go", file=sys.stderr)


def main():
    print("=" * 60, file=sys.stderr)
    print("TrueNAS Provider Generator (Native Spec)", file=sys.stderr)
    print("=" * 60, file=sys.stderr)

    methods, metadata = load_spec()
    print(f"Version: {metadata.get('truenas_version')}", file=sys.stderr)
    print(f"Methods: {len(methods)}", file=sys.stderr)

    # Find resources
    resources = {}
    for method_name in methods.keys():
        if method_name.endswith(".create"):
            base = method_name[:-7]
            resources[base] = methods

    print(f"Resources: {len(resources)}", file=sys.stderr)

    # Generate
    output_dir = Path("internal/provider")
    count = 0
    generated_resources = []

    # Skip resources with complex array handling for now
    skip_resources = {
        "nvmet.port",  # No properties in create schema
    }

    for base_name in resources.keys():
        if base_name in skip_resources:
            continue
        code = generate_resource(base_name, methods)
        if code:
            filename = f"resource_{base_name.replace('.', '_')}_generated.go"
            (output_dir / filename).write_text(code)
            generated_resources.append(base_name)
            count += 1

            # Generate documentation
            create_spec = methods.get(f"{base_name}.create", {})
            accepts = create_spec.get("accepts", [])
            if accepts:
                schema = accepts[0] if isinstance(accepts, list) else accepts

                # Handle anyOf schemas by merging all variants (same as resource generation)
                anyof_variants = None
                if "anyOf" in schema:
                    anyof_variants = schema["anyOf"]
                    merged_props = {}
                    for variant in schema["anyOf"]:
                        merged_props.update(variant.get("properties", {}))
                    properties = merged_props
                    # Only mark fields as required if they're required in ALL variants
                    all_required = [
                        set(variant.get("required", [])) for variant in schema["anyOf"]
                    ]
                    required = (
                        list(set.intersection(*all_required)) if all_required else []
                    )
                else:
                    properties = schema.get("properties", {})
                    required = schema.get("required", [])

                description = (
                    create_spec.get("description")
                    or f"Manages TrueNAS {base_name} resources"
                )
                description = description.split("\n")[0][:200]
                generate_resource_docs(
                    base_name,
                    properties,
                    required,
                    description,
                    methods,
                    anyof_variants,
                )

    print(f"\n✅ Generated {count} resources", file=sys.stderr)

    # Generate data sources
    datasource_candidates = [
        "vm",
        "pool",
        "pool.dataset",
        "disk",
        "user",
        "group",
        "interface",
        "service",
    ]
    datasource_dir = Path("internal/provider")
    ds_count = 0
    generated_datasources = []

    for base_name in datasource_candidates:
        if f"{base_name}.get_instance" not in methods:
            continue
        code = generate_datasource(base_name, methods)
        if code:
            filename = f"datasource_{base_name.replace('.', '_')}_generated.go"
            (datasource_dir / filename).write_text(code)
            generated_datasources.append(base_name)
            ds_count += 1

            # Generate documentation
            get_spec = methods.get(f"{base_name}.get_instance", {})
            returns = get_spec.get("returns", [])
            if returns:
                schema = returns[0] if isinstance(returns, list) else returns
                properties = schema.get("properties", {})
                description = (
                    get_spec.get("description") or f"Retrieves TrueNAS {base_name} data"
                )
                description = description.split("\n")[0][:200]
                generate_datasource_docs(base_name, properties, description)

    print(f"✅ Generated {ds_count} data sources", file=sys.stderr)

    # Generate query data sources
    query_candidates = [
        "vm",
        "pool",
        "pool.dataset",
        "disk",
        "user",
        "group",
        "interface",
        "service",
    ]
    query_dir = Path("internal/provider")
    query_count = 0
    generated_query_datasources = []

    for base_name in query_candidates:
        if f"{base_name}.query" not in methods:
            continue
        code = generate_query_datasource(base_name, methods)
        if code:
            filename = f"datasource_{base_name.replace('.', '_')}s_generated.go"
            (query_dir / filename).write_text(code)
            generated_query_datasources.append(base_name + "s")
            query_count += 1

            # Generate documentation
            query_spec = methods.get(f"{base_name}.query", {})
            returns = query_spec.get("returns", [])
            if returns:
                schema = returns[0] if isinstance(returns, list) else returns
                items_schema = schema.get("items", {})
                if isinstance(items_schema, list):
                    items_schema = items_schema[0] if items_schema else {}
                properties = items_schema.get("properties", {})
                # Filter out List types
                filtered_properties = {
                    k: v for k, v in properties.items() if get_tf_type(v) != "List"
                }
                description = (
                    query_spec.get("description") or f"Query {base_name} resources"
                )
                description = description.split("\n")[0][:200]
                generate_query_datasource_docs(
                    base_name, filtered_properties, description
                )

    print(f"✅ Generated {query_count} query data sources", file=sys.stderr)

    # Generate provider.go with only successfully generated resources
    generate_provider(
        generated_resources, generated_datasources + generated_query_datasources
    )


if __name__ == "__main__":
    main()
