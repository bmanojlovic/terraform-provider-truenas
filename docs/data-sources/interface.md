---
page_title: "truenas_interface Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Returns instance matching `id`. If `id` is not found, Validation error is raised.
---

# truenas_interface (Data Source)

Returns instance matching `id`. If `id` is not found, Validation error is raised.

## Example Usage

```terraform
data "truenas_interface" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the interface to retrieve.

### Read-Only

- `aliases` (List) - List of IP address aliases configured on the interface.
- `bridge_members` (List) - List of interface names that are members of this bridge.
- `description` (String) - Human-readable description of the interface.
- `enable_learning` (Bool) - Whether MAC address learning is enabled for bridge interfaces.
- `fake` (Bool) - Whether this is a fake/simulated interface for testing purposes.
- `ipv4_dhcp` (Bool) - Whether IPv4 DHCP is enabled for automatic IP address assignment.
- `ipv6_auto` (Bool) - Whether IPv6 autoconfiguration is enabled.
- `lag_ports` (List) - List of interface names that are members of this link aggregation group.
- `lag_protocol` (String) - Link aggregation protocol (LACP, FAILOVER, LOADBALANCE, etc.).
- `mtu` (Int64) - Maximum transmission unit size for the interface.
- `name` (String) - Name of the network interface.
- `state` (String) - Current runtime state information for the interface.
- `type` (String) - Type of interface (PHYSICAL, BRIDGE, LINK_AGGREGATION, VLAN, etc.).
- `vlan_parent_interface` (String) - Parent interface for VLAN configuration.
- `vlan_pcp` (Int64) - Priority Code Point for VLAN traffic prioritization.
- `vlan_tag` (Int64) - VLAN tag number for VLAN interfaces.
