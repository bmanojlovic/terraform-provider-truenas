#!/usr/bin/env python3
"""Generate Terraform provider code and documentation from TrueNAS OpenAPI spec."""
import json
import sys
import re
from pathlib import Path

# OpenAPI to Terraform type mapping
TYPE_MAP = {
    'string': 'String', 'integer': 'Int64', 'number': 'Float64',
    'boolean': 'Bool', 'array': 'List', 'object': 'Object'
}

SKIP_PATTERNS = [
    '/auth/', '/system/', '/config/', '/debug/', '/test', '/validate',
    '/available', '/choices', '/schemas', '/query', '/stats', '/status',
    '/restart', '/reload', '/start', '/stop', '/sync', '/backup', '/restore',
    'get_', 'list_', 'check_', 'download', 'export', 'import_', 
    '_run', 'execute', 'verify', 'details', '_info', 'set_', 'attach_', 
    'onetime', '_ticket'
]

# Action endpoints that look like resources but aren't
SKIP_SUFFIXES = [
    '/scrub', '/run', '/sync', '/abort', '/dismiss', '/activate', '/clone',
    '/attach', '/detach', '/replace', '/upgrade', '/rollback', '/redeploy',
    '/pull', '/push', '/send', '/receive', '/promote', '/demote'
]


def get_schema(spec, schema_name):
    """Get schema by name from spec."""
    return spec.get('components', {}).get('schemas', {}).get(schema_name)


def get_create_schema_name(path_spec):
    """Extract schema name from POST operation."""
    schema = path_spec.get('post', {}).get('requestBody', {}).get('content', {}).get('application/json', {}).get('schema', {})
    return schema.get('$ref', '').split('/')[-1] or None


def extract_resource_name(path):
    """Extract resource name from API path."""
    name = path.strip('/').replace('/', '_')
    return re.sub(r'_(get_instance|id)$', '', name)


def find_crud_resources(spec):
    """Find true CRUD resources - have POST on base and DELETE on /id/{id_}."""
    resources = set()
    for path in spec['paths']:
        if path.endswith('/id/{id_}'):
            base = path.replace('/id/{id_}', '')
            if 'delete' in spec['paths'][path] and 'post' in spec['paths'].get(base, {}):
                resources.add(base)
    return resources


def find_operational_actions(spec, resource_path):
    """Find operational actions (sync, scrub, run, upgrade, etc.) for a resource."""
    operational_verbs = ['sync', 'scrub', 'run', 'upgrade', 'rollback', 'redeploy', 'restore', 'backup']
    actions = {}
    
    for path in spec['paths']:
        if path.startswith(resource_path + '/') and 'post' in spec['paths'][path]:
            action = path.replace(resource_path + '/', '').split('/')[0]
            # Handle /resource/id/{id_}/action pattern
            if action == 'id' and len(path.split('/')) > 3:
                action = path.split('/')[-1]
            
            if any(verb in action for verb in operational_verbs):
                # Get schema for action
                post_spec = spec['paths'][path].get('post', {})
                schema_ref = post_spec.get('requestBody', {}).get('content', {}).get('application/json', {}).get('schema', {}).get('$ref', '')
                schema_name = schema_ref.split('/')[-1] if schema_ref else None
                actions[action] = {'path': path, 'schema': schema_name}
    
    return actions


def find_lifecycle_actions(spec, resource_path):
    """Find lifecycle actions (start/stop/restart) for a resource."""
    actions = []
    # Check both /resource/action and /resource/id/{id_}/action patterns
    for path in spec['paths']:
        if path.startswith(resource_path + '/'):
            action = path.replace(resource_path + '/', '').split('/')[0]
            # Handle /resource/id/{id_}/action pattern
            if action == 'id' and len(path.split('/')) > 3:
                action = path.split('/')[-1]
            if action in ['start', 'stop', 'restart'] and 'post' in spec['paths'][path]:
                actions.append(action)
    return actions


def is_crud_endpoint(path, methods, crud_resources):
    """Check if endpoint is a true CRUD resource."""
    return path in crud_resources


def get_example_value(prop_name, prop_spec):
    """Get example value from schema: enum > default > type-based."""
    prop_type = prop_spec.get('type', 'string')
    
    if prop_spec.get('enum'):
        val = prop_spec['enum'][0]
        return f'"{val}"' if isinstance(val, str) else str(val)
    
    if 'default' in prop_spec and prop_spec['default'] is not None:
        val = prop_spec['default']
        if isinstance(val, str):
            return f'"{val}"'
        if isinstance(val, bool):
            return str(val).lower()
        return str(val)
    
    return {
        'string': f'"example-{prop_name}"',
        'integer': '1',
        'boolean': 'false',
        'array': '[]',
        'object': '{}'
    }.get(prop_type, '"example-value"')


def generate_action_resource(base_resource_name, action_name, action_info, spec):
    """Generate Go code for an action resource."""
    action_path = action_info['path']
    schema_name = action_info['schema']
    
    # Get schema if available
    properties = {}
    if schema_name:
        schema = get_schema(spec, schema_name)
        if schema:
            properties = schema.get('properties', {})
    
    # Build fields and params
    fields = ['\tID types.String `tfsdk:"id"`']
    schema_attrs = ['\t\t\t"id": schema.StringAttribute{\n\t\t\t\tComputed: true,\n\t\t\t},']
    action_params = []
    
    # Add resource_id field to reference the base resource
    fields.append(f'\tResourceID types.String `tfsdk:"resource_id"`')
    schema_attrs.append('\t\t\t"resource_id": schema.StringAttribute{\n\t\t\t\tRequired: true,\n\t\t\t\tDescription: "ID of the resource to perform action on",\n\t\t\t},')
    
    # Add fields from action schema
    for prop_name, prop_spec in properties.items():
        if prop_name in ['uuid', 'id']:
            continue
        go_name = ''.join(w.capitalize() for w in prop_name.replace('-', '_').split('_'))
        go_type = TYPE_MAP.get(prop_spec.get('type', 'string'), 'String')
        
        fields.append(f'\t{go_name} types.{go_type} `tfsdk:"{prop_name}"`')
        schema_attrs.append(f'\t\t\t"{prop_name}": schema.{go_type}Attribute{{\n\t\t\t\tOptional: true,\n\t\t\t}},')
        
        if go_type in ['String', 'Int64', 'Bool']:
            method = {'String': 'ValueString', 'Int64': 'ValueInt64', 'Bool': 'ValueBool'}[go_type]
            action_params.append(f'\t\tif !data.{go_name}.IsNull() {{\n\t\t\tparams["{prop_name}"] = data.{go_name}.{method}()\n\t\t}}')
    
    resource_name = f'{base_resource_name.title().replace("_", "")}{action_name.title().replace("_", "")}Action'
    api_call = action_path.strip('/')
    
    # Format action_params with params declaration if needed
    if action_params:
        params_code = '\tparams := map[string]interface{}{}\n' + chr(10).join(action_params)
    else:
        params_code = '\t// No additional parameters'
    
    return ACTION_RESOURCE_TEMPLATE.format(
        resource_name=resource_name,
        base_name=base_resource_name,
        action_name=action_name,
        api_call=api_call,
        fields=chr(10).join(fields),
        schema_attrs=chr(10).join(schema_attrs),
        action_params=params_code
    )


def generate_resource(name, path, schema, spec):
    """Generate Go resource code."""
    properties = schema.get('properties', {})
    required_fields = schema.get('required', [])
    
    # Check for lifecycle actions
    lifecycle_actions = find_lifecycle_actions(spec, path)
    has_start = 'start' in lifecycle_actions
    
    fields, schema_attrs, create_params = ['\tID types.String `tfsdk:"id"`'], [], []
    schema_attrs.append('\t\t\t"id": schema.StringAttribute{\n\t\t\t\tComputed: true,\n\t\t\t},')
    
    # Add start_on_create field if resource has start action
    if has_start:
        fields.append('\tStartOnCreate types.Bool `tfsdk:"start_on_create"`')
        schema_attrs.append('\t\t\t"start_on_create": schema.BoolAttribute{\n\t\t\t\tOptional: true,\n\t\t\t\tDescription: "Start the resource immediately after creation (default: true if not specified)",\n\t\t\t},')
    
    for prop_name, prop_spec in properties.items():
        if prop_name in ['uuid', 'id']:
            continue
        
        go_name = ''.join(w.capitalize() for w in prop_name.replace('-', '_').split('_'))
        go_type = TYPE_MAP.get(prop_spec.get('type', 'string'), 'String')
        required = prop_name in required_fields
        
        fields.append(f'\t{go_name} types.{go_type} `tfsdk:"{prop_name}"`')
        schema_attrs.append(f'\t\t\t"{prop_name}": schema.{go_type}Attribute{{\n\t\t\t\tRequired: {str(required).lower()},\n\t\t\t\tOptional: {str(not required).lower()},\n\t\t\t}},')
        
        # Build params conditionally
        if go_type in ['String', 'Int64', 'Bool']:
            method = {'String': 'ValueString', 'Int64': 'ValueInt64', 'Bool': 'ValueBool'}[go_type]
            if required:
                create_params.append(f'\tparams["{prop_name}"] = data.{go_name}.{method}()')
            else:
                # Optional fields - only include if not null
                create_params.append(f'\tif !data.{go_name}.IsNull() {{\n\t\tparams["{prop_name}"] = data.{go_name}.{method}()\n\t}}')
    
    resource_name = name.replace('/', '_').replace('-', '_').title().replace('_', '')
    api_name = path.strip('/')
    
    # Generate lifecycle action code
    lifecycle_code = ''
    if has_start:
        lifecycle_code = f'''
\t// Handle lifecycle action - start on create if requested
\tstartOnCreate := true  // default when not specified
\tif !data.StartOnCreate.IsNull() {{
\t\tstartOnCreate = data.StartOnCreate.ValueBool()
\t}}
\tif startOnCreate {{
\t\t_, err = r.client.Call("{api_name}.start", data.ID.ValueString())
\t\tif err != nil {{
\t\t\tresp.Diagnostics.AddWarning("Start Failed", fmt.Sprintf("Resource created but failed to start: %s", err.Error()))
\t\t}}
\t}}'''
    
    return GO_RESOURCE_TEMPLATE.format(
        resource_name=resource_name,
        name=name.replace('/', '_').replace('-', '_'),
        api_name=api_name,
        fields=chr(10).join(fields),
        schema_attrs=chr(10).join(schema_attrs),
        create_params=chr(10).join(create_params),
        lifecycle_code=lifecycle_code
    )


def generate_resource_docs(name, schema, spec, path):
    """Generate Terraform resource documentation."""
    properties = schema.get('properties', {})
    required_fields = schema.get('required', [])
    description = schema.get('description', f'Manages TrueNAS {name} resources')
    resource_type = name.replace('/', '_').replace('-', '_')
    
    # Check for lifecycle actions
    lifecycle_actions = find_lifecycle_actions(spec, path)
    has_start = 'start' in lifecycle_actions
    
    # Build example lines
    example_lines = []
    for prop_name, prop_spec in properties.items():
        if prop_name in ['uuid', 'id'] or prop_name not in required_fields:
            continue
        example_lines.append(f'  {prop_name} = {get_example_value(prop_name, prop_spec)}')
    
    # Add start_on_create to example if resource has start action
    if has_start and len(example_lines) < 8:
        example_lines.append('  start_on_create = true')
    
    for prop_name, prop_spec in properties.items():
        if prop_name in ['uuid', 'id'] or prop_name in required_fields or len(example_lines) >= 8:
            continue
        if prop_spec.get('enum') or (prop_spec.get('default') is not None):
            example_lines.append(f'  {prop_name} = {get_example_value(prop_name, prop_spec)}')
    
    # Build args
    args = []
    
    # Add start_on_create if resource has start action
    if has_start:
        args.append('- `start_on_create` (Optional) - Start the resource immediately after creation. Default behavior: starts if not specified. Type: `boolean`')
    
    for prop_name, prop_spec in properties.items():
        if prop_name in ['uuid', 'id']:
            continue
        prop_desc = prop_spec.get('description', f'{prop_name} configuration')
        prop_type = prop_spec.get('type', 'string')
        required = '(Required)' if prop_name in required_fields else '(Optional)'
        
        if prop_spec.get('enum'):
            prop_desc += f' Valid values: {", ".join(f"`{v}`" for v in prop_spec["enum"][:3])}'
        if prop_spec.get('default') is not None:
            prop_desc += f' Default: `{prop_spec["default"]}`'
        
        args.append(f'- `{prop_name}` {required} - {prop_desc}. Type: `{prop_type}`')
    
    required_args = [a for a in args if '(Required)' in a]
    optional_args = [a for a in args if '(Optional)' in a]
    
    return RESOURCE_DOC_TEMPLATE.format(
        resource_type=resource_type,
        description=description,
        example_block=chr(10).join(example_lines) or '  # Configuration here',
        required_args=chr(10).join(required_args) or '- None',
        optional_args=chr(10).join(optional_args) or '- None'
    )


def generate_action_resource_docs(base_resource_name, action_name, action_info, spec):
    """Generate documentation for an action resource."""
    schema_name = action_info['schema']
    description = f'Executes {action_name} action on {base_resource_name} resource'
    
    # Get schema properties if available
    properties = {}
    if schema_name:
        schema = get_schema(spec, schema_name)
        if schema:
            properties = schema.get('properties', {})
            description = schema.get('description', description)
    
    # Build args
    args = ['- `resource_id` (Required) - ID of the resource to perform action on. Type: `string`']
    for prop_name, prop_spec in properties.items():
        if prop_name in ['uuid', 'id']:
            continue
        prop_desc = prop_spec.get('description', f'{prop_name} parameter')
        prop_type = prop_spec.get('type', 'string')
        args.append(f'- `{prop_name}` (Optional) - {prop_desc}. Type: `{prop_type}`')
    
    resource_type = f'{base_resource_name}_{action_name}_action'
    
    return f'''---
page_title: "truenas_{resource_type} Resource - terraform-provider-truenas"
subcategory: "Actions"
description: |-
  {description}
---

# truenas_{resource_type} (Resource)

{description}

This is an action resource that executes an operation when created or updated. It cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_{base_resource_name}" "example" {{
  # ... resource configuration
}}

resource "truenas_{resource_type}" "example" {{
  resource_id = truenas_{base_resource_name}.example.id
}}
```

## Schema

### Required

- `resource_id` (String) ID of the {base_resource_name} resource to perform action on

### Optional

{chr(10).join(args[1:]) or '- None'}

### Read-Only

- `id` (String) Action execution ID (timestamp-based)

## Notes

- This resource executes the {action_name} action when created
- Updates will re-execute the action
- Deletion removes from state but cannot undo the action
- Use with caution as actions are immediate and irreversible
'''


def generate_datasource_docs(name, schema):
    """Generate Terraform data source documentation."""
    properties = schema.get('properties', {})
    description = schema.get('description', f'Retrieves TrueNAS {name} data')
    resource_type = name.replace('/', '_').replace('-', '_')
    
    attrs = [f'- `{p}` ({s.get("type", "string")}) - {s.get("description", f"{p} value")}' 
             for p, s in properties.items()]
    
    return DATASOURCE_DOC_TEMPLATE.format(
        resource_type=resource_type,
        description=description,
        name=name,
        attrs=chr(10).join(attrs) or '- None'
    )


def main():
    spec_file = Path('truenas-openapi.json')
    if not spec_file.exists():
        print("truenas-openapi.json not found")
        sys.exit(1)
    
    spec = json.loads(spec_file.read_text())
    
    # Create directories
    for d in ['docs', 'docs/resources', 'docs/data-sources']:
        Path(d).mkdir(exist_ok=True)
    
    Path('docs/index.md').write_text(PROVIDER_DOC)
    
    # Find true CRUD resources
    crud_resources = find_crud_resources(spec)
    
    resources, datasources = [], []
    
    for path, path_spec in spec['paths'].items():
        schema_name = get_create_schema_name(path_spec)
        if not schema_name:
            continue
        
        schema = get_schema(spec, schema_name)
        if not schema:
            continue
        
        resource_name = extract_resource_name(path)
        
        # Data sources from get_instance
        if path.endswith('/get_instance'):
            print(f"Generating data source for {path} -> {resource_name}")
            doc = generate_datasource_docs(resource_name, schema)
            Path(f'docs/data-sources/{resource_name.lower()}.md').write_text(doc)
            datasources.append(resource_name)
            continue
        
        # Resources from CRUD endpoints
        if is_crud_endpoint(path, path_spec.keys(), crud_resources):
            print(f"Generating resource for {path} -> {resource_name}")
            
            code = generate_resource(resource_name, path, schema, spec)
            Path(f'internal/provider/resource_{resource_name.lower()}_generated.go').write_text(code)
            
            doc = generate_resource_docs(resource_name, schema, spec, path)
            Path(f'docs/resources/{resource_name.lower()}.md').write_text(doc)
            
            resources.append(resource_name)
            
            # Generate action resources for this CRUD resource
            operational_actions = find_operational_actions(spec, path)
            for action_name, action_info in operational_actions.items():
                print(f"  Generating action resource: {resource_name}_{action_name}")
                action_code = generate_action_resource(resource_name, action_name, action_info, spec)
                Path(f'internal/provider/action_{resource_name.lower()}_{action_name.lower()}_generated.go').write_text(action_code)
                
                action_doc = generate_action_resource_docs(resource_name, action_name, action_info, spec)
                Path(f'docs/resources/{resource_name.lower()}_{action_name.lower()}_action.md').write_text(action_doc)
                
                resources.append(f'{resource_name}_{action_name}_action')
    
    # Output summary
    print(f"\nGenerated {len(resources)} resources (including action resources)")
    print(f"Generated {len(datasources)} data sources")
    
    funcs = '\n'.join(f'\t\tNew{r.replace("_action", "Action").title().replace("_", "")}Resource,' for r in resources)
    print(f"\nAdd to provider.go Resources method:\nfunc (p *TrueNASProvider) Resources(ctx context.Context) []func() resource.Resource {{\n\treturn []func() resource.Resource{{\n{funcs}\n\t}}\n}}")


# Templates
GO_RESOURCE_TEMPLATE = '''package provider

import (
\t"context"
\t"fmt"

\t"github.com/hashicorp/terraform-plugin-framework/resource"
\t"github.com/hashicorp/terraform-plugin-framework/resource/schema"
\t"github.com/hashicorp/terraform-plugin-framework/types"
\t"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type {resource_name}Resource struct {{
\tclient *client.Client
}}

type {resource_name}ResourceModel struct {{
{fields}
}}

func New{resource_name}Resource() resource.Resource {{
\treturn &{resource_name}Resource{{}}
}}

func (r *{resource_name}Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {{
\tresp.TypeName = req.ProviderTypeName + "_{name}"
}}

func (r *{resource_name}Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {{
\tresp.Schema = schema.Schema{{
\t\tMarkdownDescription: "TrueNAS {name} resource",
\t\tAttributes: map[string]schema.Attribute{{
{schema_attrs}
\t\t}},
\t}}
}}

func (r *{resource_name}Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {{
\tif req.ProviderData == nil {{
\t\treturn
\t}}
\tclient, ok := req.ProviderData.(*client.Client)
\tif !ok {{
\t\tresp.Diagnostics.AddError("Unexpected Resource Configure Type", "Expected *client.Client")
\t\treturn
\t}}
\tr.client = client
}}

func (r *{resource_name}Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {{
\tvar data {resource_name}ResourceModel
\tresp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
\tif resp.Diagnostics.HasError() {{
\t\treturn
\t}}

\tparams := map[string]interface{{}}{{}}
{create_params}

\tresult, err := r.client.Call("{api_name}.create", params)
\tif err != nil {{
\t\tresp.Diagnostics.AddError("Client Error", err.Error())
\t\treturn
\t}}

\tif resultMap, ok := result.(map[string]interface{{}}); ok {{
\t\tif id, exists := resultMap["id"]; exists {{
\t\t\tdata.ID = types.StringValue(fmt.Sprintf("%v", id))
\t\t}}
\t}}
{lifecycle_code}
\tresp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}}

func (r *{resource_name}Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {{
\tvar data {resource_name}ResourceModel
\tresp.Diagnostics.Append(req.State.Get(ctx, &data)...)
\tif resp.Diagnostics.HasError() {{
\t\treturn
\t}}

\t_, err := r.client.Call("{api_name}.get_instance", data.ID.ValueString())
\tif err != nil {{
\t\tresp.Diagnostics.AddError("Client Error", err.Error())
\t\treturn
\t}}
\tresp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}}

func (r *{resource_name}Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {{
\tvar data {resource_name}ResourceModel
\tresp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
\tif resp.Diagnostics.HasError() {{
\t\treturn
\t}}

\tparams := map[string]interface{{}}{{}}
{create_params}

\t_, err := r.client.Call("{api_name}.update", []interface{{}}{{data.ID.ValueString(), params}})
\tif err != nil {{
\t\tresp.Diagnostics.AddError("Client Error", err.Error())
\t\treturn
\t}}
\tresp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}}

func (r *{resource_name}Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {{
\tvar data {resource_name}ResourceModel
\tresp.Diagnostics.Append(req.State.Get(ctx, &data)...)
\tif resp.Diagnostics.HasError() {{
\t\treturn
\t}}

\t_, err := r.client.Call("{api_name}.delete", data.ID.ValueString())
\tif err != nil {{
\t\tresp.Diagnostics.AddError("Client Error", err.Error())
\t\treturn
\t}}
}}
'''

RESOURCE_DOC_TEMPLATE = '''---
page_title: "truenas_{resource_type} Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  {description}
---

# truenas_{resource_type} (Resource)

{description}

## Example Usage

```terraform
resource "truenas_{resource_type}" "example" {{
{example_block}
}}
```

## Schema

### Required

{required_args}

### Optional

{optional_args}

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_{resource_type}.example <id>
```
'''

DATASOURCE_DOC_TEMPLATE = '''---
page_title: "truenas_{resource_type} Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  {description}
---

# truenas_{resource_type} (Data Source)

{description}

## Example Usage

```terraform
data "truenas_{resource_type}" "example" {{
  id = "1"
}}
```

## Schema

### Required

- `id` (String) The ID of the {name} to retrieve.

### Read-Only

{attrs}
'''

PROVIDER_DOC = '''---
page_title: "Provider: TrueNAS"
description: |-
  The TrueNAS provider is used to interact with TrueNAS SCALE systems via JSON-RPC over WebSocket.
---

# TrueNAS Provider

The TrueNAS provider is used to interact with TrueNAS SCALE systems. It uses the native JSON-RPC 2.0 protocol over WebSocket for optimal performance.

## Example Usage

```terraform
terraform {
  required_providers {
    truenas = {
      source = "bmanojlovic/truenas"
    }
  }
}

provider "truenas" {
  host  = "192.168.1.100"
  token = "your-api-token"
}
```

## Schema

### Required

- `host` (String) TrueNAS host address (IP or hostname)
- `token` (String) API token for authentication

### Optional

- `port` (Number) WebSocket port (default: 80 for HTTP, 443 for HTTPS)
- `use_ssl` (Boolean) Use HTTPS/WSS (default: false)

## Authentication

The provider uses API tokens for authentication. To create an API token:

1. Log into your TrueNAS web interface
2. Go to Account -> API Keys
3. Click "Add" to create a new API key
4. Copy the generated token and use it in your provider configuration

## WebSocket Connection

This provider uses WebSocket connections with JSON-RPC 2.0 protocol, providing:

- Real-time communication
- Better performance than REST APIs
- Native TrueNAS protocol support
- Persistent connections for bulk operations
'''

ACTION_RESOURCE_TEMPLATE = '''package provider

import (
\t"context"
\t"fmt"
\t"time"

\t"github.com/hashicorp/terraform-plugin-framework/resource"
\t"github.com/hashicorp/terraform-plugin-framework/resource/schema"
\t"github.com/hashicorp/terraform-plugin-framework/types"
\t"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type {resource_name}Resource struct {{
\tclient *client.Client
}}

type {resource_name}ResourceModel struct {{
{fields}
}}

func New{resource_name}Resource() resource.Resource {{
\treturn &{resource_name}Resource{{}}
}}

func (r *{resource_name}Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {{
\tresp.TypeName = req.ProviderTypeName + "_{base_name}_{action_name}_action"
}}

func (r *{resource_name}Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {{
\tresp.Schema = schema.Schema{{
\t\tMarkdownDescription: "Executes {action_name} action on {base_name} resource",
\t\tAttributes: map[string]schema.Attribute{{
{schema_attrs}
\t\t}},
\t}}
}}

func (r *{resource_name}Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {{
\tif req.ProviderData == nil {{
\t\treturn
\t}}
\tclient, ok := req.ProviderData.(*client.Client)
\tif !ok {{
\t\tresp.Diagnostics.AddError("Unexpected Resource Configure Type", "Expected *client.Client")
\t\treturn
\t}}
\tr.client = client
}}

func (r *{resource_name}Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {{
\tvar data {resource_name}ResourceModel
\tresp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
\tif resp.Diagnostics.HasError() {{
\t\treturn
\t}}

{action_params}

\t_, err := r.client.Call("{api_call}", data.ResourceID.ValueString())
\tif err != nil {{
\t\tresp.Diagnostics.AddError("Action Failed", fmt.Sprintf("Failed to execute {action_name}: %s", err.Error()))
\t\treturn
\t}}

\t// Use timestamp as ID since actions are ephemeral
\tdata.ID = types.StringValue(fmt.Sprintf("%s-%d", data.ResourceID.ValueString(), time.Now().Unix()))
\tresp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}}

func (r *{resource_name}Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {{
\t// Actions are ephemeral - nothing to read
\tvar data {resource_name}ResourceModel
\tresp.Diagnostics.Append(req.State.Get(ctx, &data)...)
}}

func (r *{resource_name}Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {{
\t// Actions are immutable - re-execute on update
\tvar data {resource_name}ResourceModel
\tresp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
\tif resp.Diagnostics.HasError() {{
\t\treturn
\t}}

{action_params}

\t_, err := r.client.Call("{api_call}", data.ResourceID.ValueString())
\tif err != nil {{
\t\tresp.Diagnostics.AddError("Action Failed", fmt.Sprintf("Failed to execute {action_name}: %s", err.Error()))
\t\treturn
\t}}

\tdata.ID = types.StringValue(fmt.Sprintf("%s-%d", data.ResourceID.ValueString(), time.Now().Unix()))
\tresp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}}

func (r *{resource_name}Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {{
\t// Actions cannot be undone - just remove from state
}}
'''

if __name__ == '__main__':
    main()
