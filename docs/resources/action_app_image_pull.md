---
page_title: "truenas_action_app_image_pull Resource - terraform-provider-truenas"
subcategory: "Actions - Applications"
description: |-
  `image` is the name of the image to pull. Format for the name is "registry/repo/image:v1.2.3" where registry may be omitted and it will default to docker registry in this case. It can or cannot contain the tag - this will be passed as is to docker so this should be analogous to what `docker pull` expects.  `auth_config` should be specified if image to be retrieved is under a private repository.
---

# truenas_action_app_image_pull (Resource)

`image` is the name of the image to pull. Format for the name is "registry/repo/image:v1.2.3" where registry may be omitted and it will default to docker registry in this case. It can or cannot contain the tag - this will be passed as is to docker so this should be analogous to what `docker pull` expects.  `auth_config` should be specified if image to be retrieved is under a private repository.

This is an action resource that executes the `app.image.pull` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_app_image_pull" "example" {
  image_pull = "value"
}
```

## Schema

### Input Parameters

- `image_pull` (String, Required) AppImagePullArgs parameters.

### Computed Outputs

- `action_id` (String) Unique identifier for this action execution
- `job_id` (Int64) Background job ID (if applicable)
- `state` (String) Job state: SUCCESS, FAILED, or RUNNING
- `progress` (Float64) Job progress percentage (0-100)
- `result` (String) Action result data
- `error` (String) Error message if action failed

## Notes

- Actions execute immediately when the resource is created
- Background jobs are monitored until completion
- Progress updates are logged during execution
- The resource cannot be updated - changes force recreation
- Destroying the resource does not undo the action
