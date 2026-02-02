#!/usr/bin/env python3
"""TrueNAS Terraform Provider Generator - Config-driven entry point."""
import json
import sys
from pathlib import Path

import yaml

from generator.schema import get_tf_type, merge_anyof_schema, has_complex_objects, get_array_item_schema, is_complex_object, to_field_name
from generator.codegen import gen_schema_attrs, gen_fields, gen_create_params, gen_read_mapping
from generator.docs import gen_resource_docs, gen_datasource_docs, gen_action_docs

TEMPLATE_DIR = Path(__file__).parent / "templates"
CONFIG_FILE = Path(__file__).parent / "config.yaml"


def load_templates():
    return {k: (TEMPLATE_DIR / f"{k}.tmpl").read_text() for k in [
        "resource.go", "resource_update_only.go", "resource_with_json.go",
        "resource_vm_device.go", "resource_uploadable.go", "action_resource.go",
        "action_uploadable.go", "resource_doc.md", "datasource.go",
        "datasource_doc.md", "datasource_query.go", "datasource_query_doc.md",
    ]}


def load_config():
    if CONFIG_FILE.exists():
        with open(CONFIG_FILE) as f:
            return yaml.safe_load(f)
    return {}


def find_latest_spec():
    specs = list(Path(".").glob("truenas-methods-*.json"))
    if not specs:
        sys.exit("ERROR: No spec file found. Run: make fetch-spec")
    return sorted(specs)[-1]


def load_spec():
    spec_file = find_latest_spec()
    print(f"Using: {spec_file}", file=sys.stderr)
    with open(spec_file) as f:
        data = json.load(f)
    return data.get("methods", {}), data.get("_metadata", {})


def gen_resource(base_name, methods, templates):
    """Generate resource file from method specs."""
    create_spec = methods.get(f"{base_name}.create", {})
    update_spec = methods.get(f"{base_name}.update", {})
    delete_spec = methods.get(f"{base_name}.delete", {})

    method_spec = create_spec or update_spec
    if not method_spec or not method_spec.get("accepts"):
        return None

    schema = method_spec["accepts"][0] if isinstance(method_spec["accepts"], list) else method_spec["accepts"]
    properties, required = merge_anyof_schema(schema)
    if not properties:
        return None

    id_is_string = any(spec.get("accepts", [{}])[0].get("type") == "string" for spec in [update_spec, delete_spec] if spec)

    update_props = {}
    if update_spec and len(update_spec.get("accepts", [])) >= 2:
        up_schema = update_spec["accepts"][1]
        update_props, _ = merge_anyof_schema(up_schema) if "anyOf" in up_schema else (up_schema.get("properties", {}), [])
        update_props = {k: v for k, v in update_props.items() if k != "id"}
    create_only = set(properties.keys()) - set(update_props.keys()) if update_props else set(properties.keys())
    properties = {**properties, **update_props}

    has_start = f"{base_name}.start" in methods and base_name != "app"
    has_stop = f"{base_name}.stop" in methods
    create_is_job = create_spec.get("job", False)
    update_is_job = update_spec.get("job", False)
    delete_is_job = delete_spec.get("job", False)
    delete_needs_opts = len(delete_spec.get("accepts", [])) >= 2

    resource_name = base_name.replace(".", "_").title().replace("_", "")
    tf_name = base_name.replace(".", "_")
    api_name = base_name
    desc = (method_spec.get("description") or f"TrueNAS {tf_name} resource").split("\n")[0][:200].replace('"', '\\"')

    # ID handling
    if id_is_string:
        id_read = "\tid = data.ID.ValueString()"
        id_update = "\tid = state.ID.ValueString()"
        id_delete = "\tid = []interface{}{data.ID.ValueString(), map[string]interface{}{}}" if delete_needs_opts else "\tid = data.ID.ValueString()"
    else:
        id_read = '\tid, err = strconv.Atoi(data.ID.ValueString())\n\tif err != nil {\n\t\tresp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))\n\t\treturn\n\t}'
        id_update = '\tid, err = strconv.Atoi(state.ID.ValueString())\n\tif err != nil {\n\t\tresp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))\n\t\treturn\n\t}'
        id_delete = id_read if not delete_needs_opts else '\tid, err = strconv.Atoi(data.ID.ValueString())\n\tif err != nil {\n\t\tresp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))\n\t\treturn\n\t}\n\tid = []interface{}{id, map[string]interface{}{}}'

    # Lifecycle
    lifecycle = ""
    if has_start:
        start_call = f"data.ID.ValueString()" if id_is_string else f"func() int {{ id, _ := strconv.Atoi(data.ID.ValueString()); return id }}()"
        lifecycle = f'\n\tstartOnCreate := true\n\tif !data.StartOnCreate.IsNull() {{ startOnCreate = data.StartOnCreate.ValueBool() }}\n\tif startOnCreate {{\n\t\t_, err = r.client.Call("{api_name}.start", {start_call})\n\t\tif err != nil {{ resp.Diagnostics.AddWarning("Start Failed", fmt.Sprintf("Resource created but failed to start: %s", err.Error())) }}\n\t}}'

    predelete = ""
    if has_stop:
        stop_call = f"data.ID.ValueString()" if id_is_string else f"func() int {{ id, _ := strconv.Atoi(data.ID.ValueString()); return id }}()"
        predelete = f'\n\t_, _ = r.client.Call("{api_name}.stop", {stop_call})\n\ttime.Sleep(2 * time.Second)\n'

    # Imports
    needs_strconv = not id_is_string or (has_start and not id_is_string) or (has_stop and not id_is_string)
    has_list = any(get_tf_type(p) == "List" and (n not in create_only or n in ("name", "type")) and n in required for n, p in properties.items())
    has_json = has_complex_objects(properties)

    imports = []
    if needs_strconv: imports.append('"strconv"')
    if has_list and required: imports.append('"github.com/hashicorp/terraform-plugin-framework/attr"')
    if has_json: imports.append('"encoding/json"')
    if has_stop: imports.append('"time"')

    mods = {get_tf_type(properties[f]) for f in create_only if f != "name" and f in properties}
    if mods & {"String", "Int64", "Bool"}:
        imports.append('"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"')
    for t, m in [("String", "stringplanmodifier"), ("Int64", "int64planmodifier"), ("Bool", "boolplanmodifier")]:
        if t in mods: imports.append(f'"github.com/hashicorp/terraform-plugin-framework/resource/schema/{m}"')

    template = templates["resource_vm_device.go"] if api_name == "vm.device" else templates["resource.go"]

    return template.format(
        resource_name=resource_name, name=tf_name, api_name=api_name, description=desc,
        fields=gen_fields(properties, has_start),
        schema_attrs=gen_schema_attrs(properties, required, has_start, create_only),
        create_params=gen_create_params(properties),
        update_params=gen_create_params(update_props or properties),
        read_mapping=gen_read_mapping(properties, create_only=create_only, required=set(required)),
        lifecycle_code=lifecycle, predelete_code=predelete,
        id_read_code=id_read, id_update_code=id_update, id_delete_code=id_delete,
        extra_imports="\n\t".join(imports),
        create_call="CallWithJob" if create_is_job else "Call",
        update_call="CallWithJob" if update_is_job else "Call",
        delete_call="CallWithJob" if delete_is_job else "Call",
    )


def gen_action_resource(method_name, method_spec, templates):
    """Generate action resource."""
    parts = method_name.split(".")
    if len(parts) < 2:
        return None

    accepts = method_spec.get("accepts", [])
    properties = {p.get("_name_", ""): p for p in accepts if p.get("_name_")}

    resource_name = "Action" + "".join(p.title() for p in parts)
    resource_type = f"action_{method_name.replace('.', '_')}"
    desc = (method_spec.get("description") or f"Execute {method_name}").replace("\n", " ").replace('"', '\\"')[:200].strip()

    # Reserved Terraform attribute names
    reserved = {"count", "for_each", "depends_on", "provider", "lifecycle"}
    def safe_attr(n):
        attr = n.replace("-", "_")
        return f"{attr}_value" if attr in reserved else attr

    fields = "\n".join(
        f'\t{to_field_name(n)} types.{get_tf_type(p) if get_tf_type(p) not in ("List", "Object") else "String"} `tfsdk:"{safe_attr(n)}"`'
        for n, p in properties.items()
    ) if properties else ""

    schema_lines = []
    for n, p in properties.items():
        tf_type = get_tf_type(p) if get_tf_type(p) not in ("List", "Object") else "String"
        req = p.get("_required_", False)
        d = p.get("description", "").replace('"', '\\"').replace("\n", " ")[:200]
        req_opt = "Required" if req else "Optional"
        attr_name = safe_attr(n)
        schema_lines.append(f'\t\t\t"{attr_name}": schema.{tf_type}Attribute{{{req_opt}: true, MarkdownDescription: "{d}"}},')

    param_lines = ["\tparams := []interface{}{}"]
    needs_json = False
    for n, p in properties.items():
        field = to_field_name(n)
        var_name = n.replace("-", "_")
        tf_type = get_tf_type(p)
        req = p.get("_required_", False)

        # Check if this is a complex object that needs JSON parsing
        if is_complex_object(p) or p.get("type") == "object":
            needs_json = True
            if req:
                param_lines.extend([
                    f"\tvar {var_name}Val interface{{}}",
                    f"\tif err := json.Unmarshal([]byte(data.{field}.ValueString()), &{var_name}Val); err == nil {{ params = append(params, {var_name}Val) }}",
                ])
            else:
                param_lines.extend([
                    f"\tif !data.{field}.IsNull() {{",
                    f"\t\tvar {var_name}Val interface{{}}",
                    f"\t\tif err := json.Unmarshal([]byte(data.{field}.ValueString()), &{var_name}Val); err == nil {{ params = append(params, {var_name}Val) }}",
                    "\t}",
                ])
        elif tf_type == "List":
            needs_json = True
            if req:
                param_lines.extend([
                    f"\tvar {var_name}Val interface{{}}",
                    f"\tif err := json.Unmarshal([]byte(data.{field}.ValueString()), &{var_name}Val); err == nil {{ params = append(params, {var_name}Val) }}",
                ])
            else:
                param_lines.extend([
                    f"\tif !data.{field}.IsNull() {{",
                    f"\t\tvar {var_name}Val interface{{}}",
                    f"\t\tif err := json.Unmarshal([]byte(data.{field}.ValueString()), &{var_name}Val); err == nil {{ params = append(params, {var_name}Val) }}",
                    "\t}",
                ])
        else:
            val_method = {"String": "ValueString", "Int64": "ValueInt64", "Bool": "ValueBool", "Float64": "ValueFloat64"}.get(tf_type, "ValueString")
            if req:
                param_lines.append(f"\tparams = append(params, data.{field}.{val_method}())")
            else:
                param_lines.append(f"\tif !data.{field}.IsNull() {{ params = append(params, data.{field}.{val_method}()) }}")

    code = templates["action_resource.go"]
    for k, v in {
        "{resource_name}": resource_name, "{resource_type_name}": resource_type,
        "{fields}": fields, "{schema_attrs}": "\n".join(schema_lines),
        "{param_building}": "\n".join(param_lines), "{method_name}": method_name,
        "{description}": desc, "{is_job}": "true" if method_spec.get("job") else "false",
        "{extra_imports}": '\n\t"encoding/json"' if needs_json else "",
    }.items():
        code = code.replace(k, v)
    return code


def gen_uploadable_resource(method_name, method_spec, is_action, templates):
    """Generate uploadable resource or action."""
    parts = method_name.split(".")
    if len(parts) < 2:
        return None

    template = templates["action_uploadable.go"] if is_action else templates["resource_uploadable.go"]

    if is_action:
        resource_name = "Action" + "".join(p.title() for p in parts)
        resource_type = f"action_{method_name.replace('.', '_')}"
    else:
        resource_name = "".join(p.title() for p in parts)
        resource_type = method_name.replace(".", "_")

    endpoint = method_name.replace(".", "/")
    desc = (method_spec.get("description") or f"{'Execute' if is_action else 'Upload via'} {method_name}").replace("\n", " ").replace('"', '\\"')[:200].strip()

    accepts = method_spec.get("accepts", [])
    properties = {}
    for p in accepts:
        name = p.get("_name_", "")
        if name == "id": name = "dataset_id"
        if name: properties[name] = p

    fields = "\n".join(
        f'\t{to_field_name(n)} types.{get_tf_type(p) if get_tf_type(p) not in ("List", "Object") else "String"} `tfsdk:"{n.replace("-", "_")}"`'
        for n, p in properties.items()
    )

    schema_lines = []
    for n, p in properties.items():
        tf_type = get_tf_type(p) if get_tf_type(p) not in ("List", "Object") else "String"
        req = p.get("_required_", False)
        d = p.get("description", "").replace('"', '\\"').replace("\n", " ")[:200]
        req_opt = "Required" if req else "Optional"
        attr_name = n.replace("-", "_")
        schema_lines.append(f'\t\t\t"{attr_name}": schema.{tf_type}Attribute{{{req_opt}: true, MarkdownDescription: "{d}"}},')

    param_lines = ["\tparams := make(map[string]interface{})"]
    needs_json = False
    for n, p in properties.items():
        field = to_field_name(n)
        var_name = n.replace("-", "_")
        tf_type = get_tf_type(p)
        req = p.get("_required_", False)
        api_name = "id" if n == "dataset_id" else n

        if tf_type in ("List", "Object"):
            needs_json = True
            if req:
                param_lines.extend([
                    f"\tvar {var_name}Val interface{{}}",
                    f'\tif err := json.Unmarshal([]byte(data.{field}.ValueString()), &{var_name}Val); err == nil {{ params["{api_name}"] = {var_name}Val }}',
                ])
            else:
                param_lines.extend([
                    f"\tif !data.{field}.IsNull() {{",
                    f"\t\tvar {var_name}Val interface{{}}",
                    f'\t\tif err := json.Unmarshal([]byte(data.{field}.ValueString()), &{var_name}Val); err == nil {{ params["{api_name}"] = {var_name}Val }}',
                    "\t}",
                ])
        else:
            val_method = {"String": "ValueString", "Int64": "ValueInt64", "Bool": "ValueBool", "Float64": "ValueFloat64"}.get(tf_type, "ValueString")
            if req:
                param_lines.append(f'\tparams["{api_name}"] = data.{field}.{val_method}()')
            else:
                param_lines.append(f'\tif !data.{field}.IsNull() {{ params["{api_name}"] = data.{field}.{val_method}() }}')

    if properties:
        first_name = list(properties.keys())[0]
        first_field = to_field_name(first_name)
        first_type = get_tf_type(properties[first_name])
        if first_type == "String":
            id_gen = f"\tdata.ID = data.{first_field}"
        elif first_type == "Int64":
            id_gen = f'\tdata.ID = types.StringValue(fmt.Sprintf("%d", data.{first_field}.ValueInt64()))'
        else:
            id_gen = f"\tdata.ID = types.StringValue(data.{first_field}.String())"
    else:
        id_gen = f'\tdata.ID = types.StringValue("{method_name}")'

    for k, v in {
        "{resource_name}": resource_name, "{resource_type_name}": resource_type,
        "{fields}": fields, "{schema_attrs}": "\n".join(schema_lines),
        "{param_building}": "\n".join(param_lines), "{method_name}": method_name,
        "{endpoint_path}": endpoint, "{description}": desc,
        "{is_job}": "true" if method_spec.get("job") else "false",
        "{id_generation}": id_gen,
        "{extra_imports}": '\n\t"encoding/json"' if needs_json else "",
    }.items():
        template = template.replace(k, v)
    return template


def gen_datasource(base_name, methods, templates):
    """Generate data source."""
    get_spec = methods.get(f"{base_name}.get_instance", {})
    returns = get_spec.get("returns", [])
    if not returns:
        return None

    schema = returns[0] if isinstance(returns, list) else returns
    properties = schema.get("properties", {})
    if not properties:
        return None

    resource_name = base_name.replace(".", "_").title().replace("_", "")
    tf_name = base_name.replace(".", "_")
    desc = (get_spec.get("description") or f"Retrieves TrueNAS {tf_name} data").split("\n")[0][:200].replace('"', '\\"')

    id_type = get_tf_type(properties.get("id", {"type": "string"}))
    if id_type == "Int64":
        id_param = "func() int { id, _ := strconv.Atoi(data.ID.ValueString()); return id }()"
        extra_imports = '\n\t"strconv"'
    else:
        id_param = "data.ID.ValueString()"
        extra_imports = ""

    if any(get_tf_type(p) == "List" for p in properties.values()):
        extra_imports += '\n\t"github.com/hashicorp/terraform-plugin-framework/attr"'

    return templates["datasource.go"].format(
        resource_name=resource_name, name=tf_name, api_name=base_name, description=desc,
        fields=gen_fields(properties, False),
        schema_attrs=gen_schema_attrs(properties, [], False),
        read_mapping=gen_read_mapping(properties, skip_id=True),
        extra_imports=extra_imports, id_param=id_param,
    )


def gen_query_datasource(base_name, methods, templates):
    """Generate query data source."""
    query_spec = methods.get(f"{base_name}.query", {})
    returns = query_spec.get("returns", [])
    if not returns:
        return None

    schema = returns[0] if isinstance(returns, list) else returns
    if "anyOf" in schema:
        for v in schema["anyOf"]:
            if v.get("type") == "array":
                schema = v
                break
    if schema.get("type") != "array":
        return None

    items = schema.get("items", {})
    items = items[0] if isinstance(items, list) else items
    properties = {k: v for k, v in items.get("properties", {}).items() if get_tf_type(v) != "List"}
    if not properties:
        return None

    resource_name = base_name.replace(".", "_").title().replace("_", "") + "s"
    tf_name = base_name.replace(".", "_") + "s"
    desc = (query_spec.get("description") or f"Query {tf_name}").split("\n")[0][:200].replace('"', '\\"')

    read_lines = []
    for n, p in properties.items():
        if n == "provider": continue
        field = to_field_name(n)
        tf_type = get_tf_type(p)
        read_lines.append(f'\t\tif v, ok := resultMap["{n}"]; ok && v != nil {{')
        if tf_type == "Bool":
            read_lines.append(f"\t\t\tif bv, ok := v.(bool); ok {{ itemModel.{field} = types.BoolValue(bv) }}")
        elif tf_type == "Int64" and n != "id":
            read_lines.append(f"\t\t\tif fv, ok := v.(float64); ok {{ itemModel.{field} = types.Int64Value(int64(fv)) }}")
        elif tf_type == "Float64":
            read_lines.append(f"\t\t\tif fv, ok := v.(float64); ok {{ itemModel.{field} = types.Float64Value(fv) }}")
        else:
            read_lines.append(f'\t\t\titemModel.{field} = types.StringValue(fmt.Sprintf("%v", v))')
        read_lines.append("\t\t}")

    attr_lines = []
    for n, p in sorted(properties.items()):
        if n == "provider": continue
        tf_type = "String" if n == "id" else get_tf_type(p)
        if tf_type == "List": continue
        attr_name = n.lower() if n != "CSR" else "csr"
        type_map = {"String": "StringType", "Int64": "Int64Type", "Bool": "BoolType", "Float64": "Float64Type"}
        if tf_type in type_map:
            attr_lines.append(f'\t\t\t"{attr_name}": types.{type_map[tf_type]},')

    return templates["datasource_query.go"].format(
        resource_name=resource_name, name=tf_name, api_name=base_name, description=desc,
        fields=gen_fields(properties, False),
        schema_attrs=gen_schema_attrs(properties, [], False),
        read_mapping="\n".join(read_lines), attr_types="\n".join(attr_lines),
    )


def gen_provider(resources, datasources, actions, uploadables, templates):
    """Generate provider.go."""
    template_path = Path(__file__).parent / "templates" / "provider.go.tmpl"
    template = template_path.read_text()

    resource_funcs = [f"New{r.replace('.', '_').title().replace('_', '')}Resource" for r in resources]
    uploadable_funcs = [f"New{''.join(p.title() for p in u.split('.'))}Resource" for u in uploadables]
    action_funcs = [f"NewAction{''.join(p.title() for p in a.split('.'))}Resource" for a in actions]

    all_funcs = resource_funcs + uploadable_funcs + action_funcs
    ds_funcs = [f"New{d.replace('.', '_').title().replace('_', '')}DataSource" for d in datasources]

    code = template.replace("{{resource_list}}", ",\n\t\t".join(all_funcs))
    code = code.replace("{{datasource_list}}", ",\n\t\t".join(ds_funcs) + ("," if ds_funcs else ""))

    Path("internal/provider/provider.go").write_text(code)
    print("✅ Generated provider.go", file=sys.stderr)


def main():
    print("=" * 60, file=sys.stderr)
    print("TrueNAS Provider Generator", file=sys.stderr)
    print("=" * 60, file=sys.stderr)

    methods, metadata = load_spec()
    config = load_config()
    templates = load_templates()
    
    print(f"Version: {metadata.get('truenas_version')}", file=sys.stderr)
    print(f"Methods: {len(methods)}", file=sys.stderr)

    output_dir = Path("internal/provider")
    skip = {"nvmet.port"}

    # Get config values
    ds_candidates = config.get("datasources", ["vm", "pool", "pool.dataset", "disk", "user", "group", "interface", "service"])
    query_ds_candidates = config.get("query_datasources", ds_candidates)
    action_keywords = config.get("actions", {}).get("keywords", ["start", "stop", "restart", "run", "sync", "scrub", "backup", "restore", "rollback", "redeploy"])
    explicit_actions = set(config.get("actions", {}).get("explicit", {}).keys())

    # Resources
    resources = [m[:-7] for m in methods if m.endswith(".create")]
    generated_resources = []
    for base in resources:
        if base in skip:
            continue
        code = gen_resource(base, methods, templates)
        if code:
            (output_dir / f"resource_{base.replace('.', '_')}_generated.go").write_text(code)
            generated_resources.append(base)

            spec = methods.get(f"{base}.create", {})
            if spec.get("accepts"):
                schema = spec["accepts"][0] if isinstance(spec["accepts"], list) else spec["accepts"]
                props, req = merge_anyof_schema(schema)
                desc = (spec.get("description") or f"Manages {base}").split("\n")[0][:200]
                gen_resource_docs(base, props, req, desc, methods, schema.get("anyOf"), templates)

    print(f"✅ Generated {len(generated_resources)} resources", file=sys.stderr)

    # Actions
    uploadable_actions = {"mail.send", "support.attach_ticket"}
    skip_uploadable = {"pool.dataset.encryption_summary"}

    generated_actions, generated_uploadables = [], []

    for method, spec in methods.items():
        if any(method.endswith(s) for s in [".create", ".update", ".delete", ".query", ".get_instance"]):
            continue

        is_uploadable = spec.get("uploadable", False)

        if is_uploadable:
            if method in skip_uploadable:
                continue
            if method in uploadable_actions:
                code = gen_uploadable_resource(method, spec, is_action=True, templates=templates)
                if code:
                    (output_dir / f"action_{method.replace('.', '_')}_generated.go").write_text(code)
                    generated_actions.append(method)
            else:
                code = gen_uploadable_resource(method, spec, is_action=False, templates=templates)
                if code:
                    (output_dir / f"resource_{method.replace('.', '_')}_generated.go").write_text(code)
                    generated_uploadables.append(method)
            continue

        # Check if method should be an action (job, keyword match, or explicit config)
        is_action = spec.get("job") or any(k in method.split(".")[-1] for k in action_keywords) or method in explicit_actions
        
        if is_action:
            code = gen_action_resource(method, spec, templates)
            if code:
                (output_dir / f"action_{method.replace('.', '_')}_generated.go").write_text(code)
                generated_actions.append(method)

                props = {p.get("_name_", ""): p for p in spec.get("accepts", []) if p.get("_name_")}
                desc = (spec.get("description") or f"Execute {method}").replace("\n", " ").strip()
                # Get config metadata for this action if available
                action_config = config.get("actions", {}).get("explicit", {}).get(method, {})
                gen_action_docs(method, props, desc, action_config)

    print(f"✅ Generated {len(generated_actions)} actions, {len(generated_uploadables)} uploadables", file=sys.stderr)

    # Data sources
    generated_ds, generated_query = [], []

    for base in ds_candidates:
        if f"{base}.get_instance" in methods:
            code = gen_datasource(base, methods, templates)
            if code:
                (output_dir / f"datasource_{base.replace('.', '_')}_generated.go").write_text(code)
                generated_ds.append(base)

                spec = methods[f"{base}.get_instance"]
                returns = spec.get("returns", [])
                if returns:
                    schema = returns[0] if isinstance(returns, list) else returns
                    gen_datasource_docs(base, schema.get("properties", {}),
                        (spec.get("description") or f"Get {base}").split("\n")[0][:200], templates)

    for base in query_ds_candidates:
        if f"{base}.query" in methods:
            code = gen_query_datasource(base, methods, templates)
            if code:
                (output_dir / f"datasource_{base.replace('.', '_')}s_generated.go").write_text(code)
                generated_query.append(base + "s")

    print(f"✅ Generated {len(generated_ds)} datasources, {len(generated_query)} query datasources", file=sys.stderr)

    gen_provider(generated_resources, generated_ds + generated_query, generated_actions, generated_uploadables, templates)

    # Coverage tracking
    generated = set()
    for base in generated_resources:
        generated.update([f"{base}.create", f"{base}.update", f"{base}.delete", f"{base}.get_instance"])
    generated.update(generated_actions)
    generated.update(generated_uploadables)
    for base in ds_candidates:
        if f"{base}.get_instance" in methods: generated.add(f"{base}.get_instance")
        if f"{base}.query" in methods: generated.add(f"{base}.query")

    not_generated = sorted(set(methods.keys()) - generated)
    queries = [m for m in not_generated if m.endswith('.query')]
    configs = [m for m in not_generated if '.config' in m or m.endswith('.update')]
    choices = [m for m in not_generated if 'choice' in m.lower()]
    auth = [m for m in not_generated if m.startswith('auth.')]

    print(f"\n📊 Coverage: {len(generated)}/{len(methods)} methods ({100*len(generated)//len(methods)}%)", file=sys.stderr)
    print(f"   Not generated: {len(not_generated)} (queries:{len(queries)}, configs:{len(configs)}, choices:{len(choices)}, auth:{len(auth)})", file=sys.stderr)


if __name__ == "__main__":
    main()
