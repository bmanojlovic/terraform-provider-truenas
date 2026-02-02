"""TrueNAS Terraform Provider Generator."""

from .schema import get_tf_type, to_field_name, is_complex_object, has_complex_objects, get_array_item_schema, merge_anyof_schema
from .codegen import gen_schema_attrs, gen_fields, gen_create_params, gen_read_mapping
from .docs import gen_resource_docs, gen_datasource_docs, gen_action_docs

__all__ = [
    'get_tf_type', 'to_field_name', 'is_complex_object', 'has_complex_objects',
    'get_array_item_schema', 'merge_anyof_schema',
    'gen_schema_attrs', 'gen_fields', 'gen_create_params', 'gen_read_mapping',
    'gen_resource_docs', 'gen_datasource_docs', 'gen_action_docs',
]
