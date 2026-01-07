---
page_title: "truenas_interface Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Configuration data for the new interface.
---

# truenas_interface (Resource)

Configuration data for the new interface.

## Example Usage

```terraform
resource "truenas_interface" "example" {
  type = "BRIDGE"
  description = ""
  ipv4_dhcp = false
  ipv6_auto = false
  aliases = []
  failover_critical = false
  failover_aliases = []
  failover_virtual_aliases = []
}
```

## Schema

### Required

- `type` (Required) - Type of interface to create. Valid values: `BRIDGE`, `LINK_AGGREGATION`, `VLAN`. Type: `string`

### Optional

- `name` (Optional) - Generate a name if not provided based on `type`, e.g. "br0", "bond1", "vlan0".. Type: `string`
- `description` (Optional) - Human-readable description of the interface. Default: ``. Type: `string`
- `ipv4_dhcp` (Optional) - Enable IPv4 DHCP for automatic IP address assignment. Default: `False`. Type: `boolean`
- `ipv6_auto` (Optional) - Enable IPv6 autoconfiguration. Default: `False`. Type: `boolean`
- `aliases` (Optional) - List of IP address aliases to configure on the interface. Default: `[]`. Type: `array`
- `failover_critical` (Optional) - Whether this interface is critical for failover functionality. Critical interfaces are monitored for     failover events and can trigger failover when they fail. Default: `False`. Type: `boolean`
- `failover_group` (Optional) - Failover group identifier for clustering. Interfaces in the same group fail over together during     failover events.. Type: `string`
- `failover_vhid` (Optional) - Virtual Host ID for VRRP failover configuration. Must be unique within the VRRP group and match     between failover nodes.. Type: `string`
- `failover_aliases` (Optional) - List of IP aliases for failover configuration. These IPs are assigned to the interface during normal     operation and migrate during failover. Default: `[]`. Type: `array`
- `failover_virtual_aliases` (Optional) - List of virtual IP aliases for failover configuration. These are shared IPs that float between nodes     during failover events. Default: `[]`. Type: `array`
- `bridge_members` (Optional) - List of interfaces to add as members of this bridge. Default: `[]`. Type: `array`
- `enable_learning` (Optional) - Enable MAC address learning for bridge interfaces. When enabled, the bridge learns MAC addresses     from incoming frames and builds a forwarding table to optimize traffic flow. Default: `True`. Type: `boolean`
- `stp` (Optional) - Enable Spanning Tree Protocol for bridge interfaces. STP prevents network loops by blocking redundant     paths and enables automatic failover when the primary path fails. Default: `True`. Type: `boolean`
- `lag_protocol` (Optional) - Link aggregation protocol to use for bonding interfaces. LACP uses 802.3ad dynamic negotiation,     FAILOVER provides active-backup, LOADBALANCE and ROUNDROBIN distribute traffic across links. Valid values: `LACP`, `FAILOVER`, `LOADBALANCE`. Type: `string`
- `xmit_hash_policy` (Optional) - Transmit hash policy for load balancing in link aggregation. LAYER2 uses MAC addresses, LAYER2+3 adds IP     addresses, and LAYER3+4 includes TCP/UDP ports for distribution. Valid values: `LAYER2`, `LAYER2+3`, `LAYER3+4`. Type: `string`
- `lacpdu_rate` (Optional) - LACP data unit transmission rate. SLOW sends LACPDUs every 30 seconds, FAST sends every 1 second for     quicker link failure detection. Valid values: `SLOW`, `FAST`, `None`. Type: `string`
- `lag_ports` (Optional) - List of interface names to include in the link aggregation group. Default: `[]`. Type: `array`
- `vlan_parent_interface` (Optional) - Parent interface for VLAN configuration.. Type: `string`
- `vlan_tag` (Optional) - VLAN tag number (1-4094).. Type: `integer`
- `vlan_pcp` (Optional) - Priority Code Point for VLAN traffic prioritization (0-7). Values 0-7 map to different QoS priority levels,     with 0 being lowest and 7 highest priority.. Type: `string`
- `mtu` (Optional) - Maximum transmission unit size for the interface (68-9216 bytes).. Type: `string`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_interface.example <id>
```
