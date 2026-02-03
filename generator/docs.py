"""Documentation generation for resources, actions, datasources."""

from pathlib import Path
from .schema import get_tf_type, is_complex_object


# Subcategory mapping for Terraform Registry sidebar grouping
SUBCATEGORY_MAP = {
    "pool": "Storage - Pools",
    "disk": "Storage - Disks",
    "vm": "Virtual Machines",
    "virt": "Virtualization",
    "sharing": "Sharing",
    "interface": "Network",
    "staticroute": "Network",
    "user": "Users & Groups",
    "group": "Users & Groups",
    "app": "Applications",
    "docker": "Applications",
    "certificate": "Certificates",
    "acme": "Certificates",
    "service": "Services",
    "cronjob": "Scheduled Tasks",
    "replication": "Replication",
    "cloudsync": "Cloud Sync",
    "cloud_backup": "Cloud Backup",
    "boot": "System - Boot",
    "system": "System",
    "config": "System - Config",
    "filesystem": "Filesystem",
    "alert": "Alerts",
    "reporting": "Reporting",
    "audit": "Audit",
    "mail": "Notifications",
    "smart": "S.M.A.R.T.",
    "ups": "UPS",
    "iscsi": "iSCSI",
    "nfs": "Sharing",
    "smb": "Sharing",
    "ftp": "Sharing",
}


def get_subcategory(name, is_action=False):
    """Get subcategory for a resource/datasource based on its name."""
    if is_action:
        prefix = "Actions - "
    else:
        prefix = ""
    
    # Check each prefix in order of specificity (longer first)
    parts = name.replace("_", ".").split(".")
    for i in range(len(parts), 0, -1):
        key = ".".join(parts[:i])
        if key in SUBCATEGORY_MAP:
            return prefix + SUBCATEGORY_MAP[key]
    
    # Fallback to first part
    first = parts[0]
    if first in SUBCATEGORY_MAP:
        return prefix + SUBCATEGORY_MAP[first]
    
    return prefix + "Other" if is_action else "Other"


def gen_resource_docs(base_name, properties, required, description, methods, anyof_variants, templates):
    """Generate resource documentation."""
    tf_name = base_name.replace(".", "_")
    has_start = f"{base_name}.start" in methods

    # Example
    example_lines = []
    for n in sorted(required):
        if n in properties and n not in ("uuid", "id"):
            tf_type = get_tf_type(properties[n])
            val = {"String": '"example"', "Int64": "1", "Bool": "true", "Float64": "1.0", "List": '["item"]'}.get(tf_type, '"value"')
            example_lines.append(f"  {n} = {val}")
    if has_start and len(example_lines) < 8:
        example_lines.append("  start_on_create = true")

    # Args
    req_args, opt_args = [], []
    if has_start:
        opt_args.append("- `start_on_create` (Bool) - Start immediately after creation. Default: `true`")

    # Collect all enum values for discriminator field from anyOf variants
    all_enum_values = {}
    if anyof_variants:
        for fn in ["type", "kind", "variant"]:
            values = []
            for variant in anyof_variants:
                v_props = variant.get("properties", {})
                if fn in v_props:
                    dp = v_props[fn]
                    if isinstance(dp, dict) and "enum" in dp:
                        values.extend(dp["enum"])
            if values:
                all_enum_values[fn] = list(dict.fromkeys(values))

    for n, p in sorted(properties.items()):
        if n in ("provider", "uuid", "id"):
            continue
        tf_type = get_tf_type(p)
        desc = p.get("description", "").replace("\n", " ")[:200] if isinstance(p, dict) else ""
        
        if tf_type == "String" and isinstance(p, dict) and is_complex_object(p):
            desc += " **Note:** This is a JSON object. Use `jsonencode()` to pass structured data."
            obj_props = None
            if "properties" in p:
                obj_props = p["properties"]
            elif "anyOf" in p:
                for v in p["anyOf"]:
                    if isinstance(v, dict) and "properties" in v:
                        obj_props = v["properties"]
                        break
            elif "oneOf" in p:
                for v in p["oneOf"]:
                    if isinstance(v, dict) and "properties" in v:
                        obj_props = v["properties"]
                        break
            if obj_props:
                examples = []
                for pn, pv in list(obj_props.items())[:3]:
                    pt = pv.get("type", "string") if isinstance(pv, dict) else "string"
                    if pt == "string": examples.append(f'{pn} = "value"')
                    elif pt == "integer": examples.append(f"{pn} = 0")
                    elif pt == "boolean": examples.append(f"{pn} = true")
                if examples:
                    more = ", ..." if len(obj_props) > 3 else ""
                    desc += f" Example: `jsonencode({{{', '.join(examples)}{more}}})`"
        
        if isinstance(p, dict) and "default" in p:
            desc += f" Default: `{p['default']}`"
        if n in all_enum_values:
            desc += f" Valid values: {', '.join(f'`{v}`' for v in all_enum_values[n][:10])}"
        elif isinstance(p, dict) and "enum" in p:
            desc += f" Valid values: {', '.join(f'`{v}`' for v in p['enum'][:10])}"
        
        line = f"- `{n}` ({tf_type}) - {desc}"
        (req_args if n in required else opt_args).append(line)

    generic_example = f"""
## Example Usage

```terraform
resource "truenas_{tf_name}" "example" {{
{chr(10).join(example_lines) or "  # Configure required attributes"}
}}
```
""" if not anyof_variants else ""

    # Build variant examples for anyOf schemas
    variant_examples = ""
    if anyof_variants:
        disc_field = None
        for fn in ["type", "kind", "variant"]:
            if fn in properties and isinstance(properties[fn], dict) and "enum" in properties[fn]:
                disc_field = fn
                break
        
        if disc_field:
            variant_examples = f"\n## Variants\n\nThis resource has **{len(anyof_variants)} variants** controlled by the `{disc_field}` field.\n\n"
            
            for variant in anyof_variants:
                v_props = variant.get("properties", {})
                v_req = set(variant.get("required", []))
                
                v_name = None
                if disc_field in v_props:
                    dp = v_props[disc_field]
                    if isinstance(dp, dict):
                        if "enum" in dp and dp["enum"]:
                            v_name = dp["enum"][0]
                        elif "default" in dp:
                            v_name = dp["default"]
                
                if v_name:
                    variant_examples += f"### {v_name}\n\n```terraform\n"
                    variant_examples += f'resource "truenas_{tf_name}" "example" {{\n'
                    variant_examples += f'  {disc_field} = "{v_name}"\n'
                    for rn in sorted(v_req):
                        if rn != disc_field and rn in properties:
                            variant_examples += f'  {rn} = "value"\n'
                    variant_examples += "}\n```\n\n"
                    variant_examples += f"**Required fields:** {', '.join(f'`{r}`' for r in sorted(v_req))}\n\n"

    doc = templates["resource_doc.md"].format(
        resource_type=tf_name,
        subcategory=get_subcategory(base_name),
        description=description,
        required_args=chr(10).join(req_args) or "- None",
        optional_args=chr(10).join(opt_args) or "- None",
        variant_examples=variant_examples,
        generic_example=generic_example,
    )

    Path("docs/resources").mkdir(parents=True, exist_ok=True)
    Path(f"docs/resources/{tf_name}.md").write_text(doc)


def gen_datasource_docs(base_name, properties, description, templates):
    """Generate data source documentation."""
    tf_name = base_name.replace(".", "_")
    attrs = [
        f"- `{n}` ({get_tf_type(p)}) - {p.get('description', '')[:200].replace(chr(10), ' ').strip()}"
        for n, p in sorted(properties.items())
        if n != "id" and isinstance(p, dict)
    ]

    doc = templates["datasource_doc.md"].format(
        resource_type=tf_name,
        subcategory=get_subcategory(base_name),
        description=description,
        name=tf_name,
        attrs=chr(10).join(attrs) or "- None",
    )
    Path("docs/data-sources").mkdir(parents=True, exist_ok=True)
    Path(f"docs/data-sources/{tf_name}.md").write_text(doc)


def gen_action_docs(method_name, properties, description, config_meta=None):
    """Generate action documentation with optional config metadata."""
    resource_name = f"action_{method_name.replace('.', '_')}"
    config_meta = config_meta or {}

    # Use config example if provided, otherwise generate generic
    if "example" in config_meta:
        example = config_meta["example"].strip()
    else:
        example = f'resource "truenas_{resource_name}" "example" {{\n'
        for n, p in properties.items():
            if p.get("_required_"):
                tf_type = get_tf_type(p)
                val = {"String": '"value"', "Int64": "1", "Bool": "true"}.get(tf_type, '"value"')
                example += f"  {n} = {val}\n"
        example += "}"

    # Use config description if richer
    if config_meta.get("description"):
        description = config_meta["description"]

    schema_lines = []
    for n, p in properties.items():
        tf_type = get_tf_type(p)
        req = "Required" if p.get("_required_") else "Optional"
        desc = p.get("description", "").replace("\n", " ")[:200]
        schema_lines.append(f"- `{n}` ({tf_type}, {req}) {desc}")

    # Build workflow section if provided
    workflow_section = ""
    if config_meta.get("workflow"):
        workflow_section = f"\n## Workflow\n\n{config_meta['workflow'].strip()}\n"

    # Build related actions section
    related_section = ""
    if config_meta.get("related"):
        related = config_meta["related"]
        related_links = [f"- `truenas_action_{r.replace('.', '_')}`" for r in related]
        related_section = f"\n## Related Actions\n\n{chr(10).join(related_links)}\n"

    subcategory = get_subcategory(method_name, is_action=True)

    doc = f"""---
page_title: "truenas_{resource_name} Resource - terraform-provider-truenas"
subcategory: "{subcategory}"
description: |-
  {description}
---

# truenas_{resource_name} (Resource)

{description}
{workflow_section}
This is an action resource that executes the `{method_name}` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
{example}
```

## Schema

### Input Parameters

{chr(10).join(schema_lines) or "None"}

### Computed Outputs

- `action_id` (String) Unique identifier for this action execution
- `job_id` (Int64) Background job ID (if applicable)
- `state` (String) Job state: SUCCESS, FAILED, or RUNNING
- `progress` (Float64) Job progress percentage (0-100)
- `result` (String) Action result data
- `error` (String) Error message if action failed
{related_section}
## Notes

- Actions execute immediately when the resource is created
- Background jobs are monitored until completion
- Progress updates are logged during execution
- The resource cannot be updated - changes force recreation
- Destroying the resource does not undo the action
"""
    Path("docs/resources").mkdir(parents=True, exist_ok=True)
    Path(f"docs/resources/{resource_name}.md").write_text(doc)
