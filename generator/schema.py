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


def is_flattenable(accepts):
    """Check if method accepts can be flattened to top-level TF fields.
    
    Flattenable cases:
    - No params
    - All params are primitives (string, int, bool, number)
    - Single object param with simple properties (no arrays of objects, max 1 level nesting)
    """
    # No params = flattenable
    if len(accepts) == 0:
        return True
    
    # All primitives = flattenable
    primitives = ("string", "integer", "boolean", "number")
    if all(p.get("type") in primitives for p in accepts):
        return True
    
    # Single object param - check complexity
    if len(accepts) != 1:
        return False
    param = accepts[0]
    if param.get("type") != "object":
        return False
    
    props = param.get("properties", {})
    for name, prop in props.items():
        # Arrays of objects = not flattenable
        if prop.get("type") == "array":
            items = prop.get("items", {})
            item = items[0] if isinstance(items, list) else items
            if isinstance(item, dict) and item.get("type") == "object":
                return False
        # anyOf/oneOf with arrays or objects = not flattenable
        for key in ["anyOf", "oneOf"]:
            for variant in prop.get(key, []):
                if variant.get("type") == "array":
                    items = variant.get("items", {})
                    item = items[0] if isinstance(items, list) else items
                    if isinstance(item, dict) and item.get("type") == "object":
                        return False
                if variant.get("type") == "object":
                    return False
        # Nested objects OK if they're simple (no further nesting)
        if prop.get("type") == "object":
            nested_props = prop.get("properties", {})
            for np in nested_props.values():
                if np.get("type") in ("object", "array"):
                    return False
    return True


def flatten_schema(accepts):
    """Flatten params into list of (name, prop) tuples."""
    # No params
    if len(accepts) == 0:
        return []
    
    # All primitives - each param becomes a field
    primitives = ("string", "integer", "boolean", "number")
    if all(p.get("type") in primitives for p in accepts):
        result = []
        for p in accepts:
            prop = dict(p)
            prop["_required_"] = p.get("_required_", True)
            result.append((p.get("_name_", "param"), prop))
        return result
    
    # Single object param - flatten properties
    param = accepts[0]
    props = param.get("properties", {})
    required = set(param.get("required", []))
    result = []
    
    for name, prop in props.items():
        if prop.get("type") == "object":
            # Flatten nested object properties with prefix
            nested_props = prop.get("properties", {})
            nested_req = set(prop.get("required", []))
            for nname, nprop in nested_props.items():
                flat_name = f"{name}_{nname}"
                nprop = dict(nprop)
                nprop["_required_"] = nname in nested_req and name in required
                nprop["_nested_parent_"] = name
                nprop["_nested_name_"] = nname
                result.append((flat_name, nprop))
        else:
            prop = dict(prop)
            prop["_required_"] = name in required
            result.append((name, prop))
    
    return result
