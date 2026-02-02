"""Schema utilities - type mapping and schema parsing."""

TYPE_MAP = {
    "string": "String",
    "integer": "Int64",
    "number": "Float64",
    "boolean": "Bool",
    "array": "List",
    "object": "String",
}


def get_tf_type(prop):
    """Convert JSON schema to Terraform type."""
    if isinstance(prop, list):
        prop = prop[0] if prop else {}
    if not isinstance(prop, dict):
        return "String"
    if "anyOf" in prop:
        for v in prop["anyOf"]:
            if v.get("type") in ("integer", "boolean", "array"):
                return {"integer": "Int64", "boolean": "Bool", "array": "List"}[v["type"]]
        return "String"
    if "oneOf" in prop or "discriminator" in prop:
        return "String"
    return TYPE_MAP.get(prop.get("type"), "String")


def to_field_name(name):
    """Convert property name to Go field name."""
    if name == "CSR":
        return "Csr"
    if name == "id":
        return "ID"
    return name.replace("-", "_").title().replace("_", "")


def is_complex_object(prop):
    """Check if property needs JSON parsing."""
    if not isinstance(prop, dict):
        return False
    if prop.get("type") == "object":
        return True
    for key in ["anyOf", "oneOf"]:
        if any(v.get("type") == "object" for v in prop.get(key, [])):
            return True
    if "discriminator" in prop:
        return True
    if prop.get("type") == "array":
        items = prop.get("items", {})
        item = items[0] if isinstance(items, list) else items
        return isinstance(item, dict) and item.get("type") == "object"
    return False


def has_complex_objects(properties):
    return any(is_complex_object(p) for n, p in properties.items() if n != "provider")


def get_array_item_schema(prop):
    """Get item schema for array properties."""
    items = prop.get("items", {})
    return items[0] if isinstance(items, list) else items


def merge_anyof_schema(schema):
    """Merge anyOf variants into single schema."""
    if "anyOf" not in schema:
        return schema.get("properties", {}), schema.get("required", [])
    props = {}
    for v in schema["anyOf"]:
        props.update(v.get("properties", {}))
    all_req = [set(v.get("required", [])) for v in schema["anyOf"]]
    req = list(set.intersection(*all_req)) if all_req else []
    return props, req
