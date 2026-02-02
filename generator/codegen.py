"""Code generation - Go code for resources, actions, datasources."""

from .schema import get_tf_type, to_field_name, is_complex_object, has_complex_objects, get_array_item_schema, merge_anyof_schema


def gen_schema_attrs(properties, required, has_start=False, create_only=None):
    """Generate schema attributes."""
    create_only = create_only or set()
    lines = []

    if not has_start and not required:  # datasource
        lines.append('\t\t\t"id": schema.StringAttribute{Required: true, Description: "Resource ID"},')
    elif "id" not in properties:
        lines.append('\t\t\t"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},')

    if has_start:
        lines.append('\t\t\t"start_on_create": schema.BoolAttribute{Optional: true, Description: "Start the resource immediately after creation (default: true)"},')

    for name, prop in properties.items():
        if name == "id":
            if not has_start and not required:
                continue
            tf_type = get_tf_type(prop)
            lines.append(f'\t\t\t"id": schema.{tf_type}Attribute{{Computed: true, Description: "Resource ID"}},')
            continue
        if name == "provider":
            continue

        prop = prop[0] if isinstance(prop, list) else prop
        tf_type = get_tf_type(prop)
        is_req = name in required
        desc = prop.get("description", "")[:100].replace('"', '\\"').replace("\n", " ") if isinstance(prop, dict) else ""
        attr_name = name.lower() if name != "CSR" else "csr"

        lines.append(f'\t\t\t"{attr_name}": schema.{tf_type}Attribute{{')

        if not has_start and not required:  # datasource
            lines.append("\t\t\t\tComputed: true,")
        else:
            is_auto = not is_req and isinstance(prop, dict) and "generate" in prop.get("description", "").lower()
            if is_auto:
                lines.append("\t\t\t\tOptional: true,")
                lines.append("\t\t\t\tComputed: true,")
            else:
                lines.append(f"\t\t\t\tRequired: {str(is_req).lower()},")
                lines.append(f"\t\t\t\tOptional: {str(not is_req).lower()},")

        if tf_type == "List":
            lines.append("\t\t\t\tElementType: types.StringType,")
        lines.append(f'\t\t\t\tDescription: "{desc}",')

        if name in create_only and name != "name":
            mod_map = {"String": "stringplanmodifier", "Int64": "int64planmodifier", "Bool": "boolplanmodifier"}
            if tf_type in mod_map:
                lines.append(f"\t\t\t\tPlanModifiers: []planmodifier.{tf_type}{{{mod_map[tf_type]}.RequiresReplace()}},")

        lines.append("\t\t\t},")

    return "\n".join(lines)


def gen_fields(properties, has_start=False):
    """Generate struct fields."""
    lines = []
    is_ds = not has_start and "id" in properties

    if is_ds or "id" not in properties:
        lines.append('\tID types.String `tfsdk:"id"`')
    if has_start:
        lines.append('\tStartOnCreate types.Bool `tfsdk:"start_on_create"`')

    for name, prop in properties.items():
        if name == "provider" or (name == "id" and is_ds):
            continue
        field = to_field_name(name)
        tf_type = get_tf_type(prop)
        lines.append(f'\t{field} types.{tf_type} `tfsdk:"{name.lower()}"`')

    return "\n".join(lines)


def gen_create_params(properties):
    """Generate parameter building code."""
    lines = []
    for name, prop in properties.items():
        if name in ("provider", "id"):
            continue
        field = to_field_name(name)
        tf_type = get_tf_type(prop)

        lines.append(f"\tif !data.{field}.IsNull() {{")

        if tf_type == "Bool":
            lines.append(f'\t\tparams["{name}"] = data.{field}.ValueBool()')
        elif tf_type == "Int64":
            lines.append(f'\t\tparams["{name}"] = data.{field}.ValueInt64()')
        elif tf_type == "Float64":
            lines.append(f'\t\tparams["{name}"] = data.{field}.ValueFloat64()')
        elif tf_type == "List":
            item = get_array_item_schema(prop) if isinstance(prop, dict) else {}
            if isinstance(item, dict) and item.get("type") == "object":
                lines.extend([
                    f"\t\tvar {name}List []string",
                    f"\t\tdata.{field}.ElementsAs(ctx, &{name}List, false)",
                    f"\t\tvar {name}Objs []map[string]interface{{}}",
                    f"\t\tfor _, jsonStr := range {name}List {{",
                    f"\t\t\tvar obj map[string]interface{{}}",
                    f"\t\t\tif err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {{",
                    f'\t\t\t\tresp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse {name} item: %s", err))',
                    f"\t\t\t\treturn",
                    f"\t\t\t}}",
                    f"\t\t\t{name}Objs = append({name}Objs, obj)",
                    f"\t\t}}",
                    f'\t\tparams["{name}"] = {name}Objs',
                ])
            else:
                lines.extend([
                    f"\t\tvar {name}List []string",
                    f"\t\tdata.{field}.ElementsAs(ctx, &{name}List, false)",
                    f'\t\tparams["{name}"] = {name}List',
                ])
        elif is_complex_object(prop):
            lines.extend([
                f"\t\tvar {name}Obj map[string]interface{{}}",
                f"\t\tif err := json.Unmarshal([]byte(data.{field}.ValueString()), &{name}Obj); err != nil {{",
                f'\t\t\tresp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse {name}: %s", err))',
                f"\t\t\treturn",
                f"\t\t}}",
                f'\t\tparams["{name}"] = {name}Obj',
            ])
        else:
            lines.append(f'\t\tparams["{name}"] = data.{field}.ValueString()')

        lines.append("\t}")
    return "\n".join(lines)


def gen_read_mapping(properties, skip_id=False, create_only=None, required=None):
    """Generate code to map API response to state."""
    create_only = create_only or set()
    lines = []

    def should_read_field(name):
        if name in ("provider", "id"):
            return False
        if name in create_only and name not in ("name", "type"):
            return False
        if required is not None and name not in required and name not in ("name", "type"):
            return False
        return True

    fields_to_read = [n for n in properties if should_read_field(n)]
    has_fields = not skip_id or bool(fields_to_read)

    if has_fields:
        lines.extend([
            "\tresultMap, ok := result.(map[string]interface{})",
            "\tif !ok {",
            '\t\tresp.Diagnostics.AddError("Parse Error", "Failed to parse API response")',
            "\t\treturn",
            "\t}",
            "",
        ])
    else:
        lines.extend(["\t_ = result // No fields to read", ""])
        return "\n".join(lines)

    if not skip_id:
        lines.extend([
            '\t\tif v, ok := resultMap["id"]; ok && v != nil {',
            '\t\t\tdata.ID = types.StringValue(fmt.Sprintf("%v", v))',
            "\t\t}",
        ])

    for name in fields_to_read:
        prop = properties[name]
        field = to_field_name(name)
        tf_type = get_tf_type(prop)

        lines.append(f'\t\tif v, ok := resultMap["{name}"]; ok && v != nil {{')

        if tf_type == "Bool":
            lines.append(f"\t\t\tif bv, ok := v.(bool); ok {{ data.{field} = types.BoolValue(bv) }}")
        elif tf_type == "Int64":
            lines.extend([
                f"\t\t\tswitch val := v.(type) {{",
                f"\t\t\tcase float64:",
                f"\t\t\t\tdata.{field} = types.Int64Value(int64(val))",
                f"\t\t\tcase map[string]interface{{}}:",
                f'\t\t\t\tif parsed, ok := val["parsed"]; ok && parsed != nil {{',
                f"\t\t\t\t\tif fv, ok := parsed.(float64); ok {{ data.{field} = types.Int64Value(int64(fv)) }}",
                f"\t\t\t\t}}",
                f"\t\t\t}}",
            ])
        elif tf_type == "Float64":
            lines.append(f"\t\t\tif fv, ok := v.(float64); ok {{ data.{field} = types.Float64Value(fv) }}")
        elif tf_type == "List":
            lines.extend([
                f"\t\t\tif arr, ok := v.([]interface{{}}); ok {{",
                f"\t\t\t\tstrVals := make([]attr.Value, len(arr))",
                f'\t\t\t\tfor i, item := range arr {{ strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }}',
                f"\t\t\t\tdata.{field}, _ = types.ListValue(types.StringType, strVals)",
                f"\t\t\t}}",
            ])
        else:
            lines.extend([
                f"\t\t\tswitch val := v.(type) {{",
                f"\t\t\tcase string:",
                f"\t\t\t\tdata.{field} = types.StringValue(val)",
                f"\t\t\tcase map[string]interface{{}}:",
                f'\t\t\t\tif strVal, ok := val["value"]; ok && strVal != nil {{',
                f'\t\t\t\t\tdata.{field} = types.StringValue(fmt.Sprintf("%v", strVal))',
                f"\t\t\t\t}}",
                f"\t\t\tdefault:",
                f'\t\t\t\tdata.{field} = types.StringValue(fmt.Sprintf("%v", v))',
                f"\t\t\t}}",
            ])

        lines.append("\t\t}")

    return "\n".join(lines)
