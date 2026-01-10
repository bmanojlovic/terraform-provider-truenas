---
page_title: "truenas_interface Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create virtual interfaces (Link Aggregation, VLAN)
---

# truenas_interface (Resource)

Create virtual interfaces (Link Aggregation, VLAN)

## Example Usage

```terraform
resource "truenas_interface" "example" {
  type = "example-value"
}
```

## Schema

### Required

- `type` (String) - Type of interface to create. Valid values: `BRIDGE`, `LINK_AGGREGATION`, `VLAN`

### Optional

- `aliases` (List) - List of IP address aliases to configure on the interface. Default: `[]`
- `bridge_members` (List) - List of interfaces to add as members of this bridge. Default: `[]`
- `description` (String) - Human-readable description of the interface. Default: ``
- `enable_learning` (Bool) - Enable MAC address learning for bridge interfaces. When enabled, the bridge learns MAC addresses     from incoming frames and builds a forwarding table to optimize traffic flow. Default: `True`
- `failover_aliases` (List) - List of IP aliases for failover configuration. These IPs are assigned to the interface during normal     operation and migrate during failover. Default: `[]`
- `failover_critical` (Bool) - Whether this interface is critical for failover functionality. Critical interfaces are monitored for     failover events and can trigger failover when they fail. Default: `False`
- `failover_group` (Int64) - Failover group identifier for clustering. Interfaces in the same group fail over together during     failover events.
- `failover_vhid` (Int64) - Virtual Host ID for VRRP failover configuration. Must be unique within the VRRP group and match     between failover nodes.
- `failover_virtual_aliases` (List) - List of virtual IP aliases for failover configuration. These are shared IPs that float between nodes     during failover events. Default: `[]`
- `ipv4_dhcp` (Bool) - Enable IPv4 DHCP for automatic IP address assignment. Default: `False`
- `ipv6_auto` (Bool) - Enable IPv6 autoconfiguration. Default: `False`
- `lacpdu_rate` (String) - LACP data unit transmission rate. SLOW sends LACPDUs every 30 seconds, FAST sends every 1 second for     quicker link failure detection. Default: `None` Valid values: `SLOW`, `FAST`, `None`
- `lag_ports` (List) - List of interface names to include in the link aggregation group. Default: `[]`
- `lag_protocol` (String) - Link aggregation protocol to use for bonding interfaces. LACP uses 802.3ad dynamic negotiation,     FAILOVER provides active-backup, LOADBALANCE and ROUNDROBIN distribute traffic across links. Valid values: `LACP`, `FAILOVER`, `LOADBALANCE`, `ROUNDROBIN`, `NONE`
- `mtu` (Int64) - Maximum transmission unit size for the interface (68-9216 bytes). Default: `None`
- `name` (String) - Generate a name if not provided based on `type`, e.g. "br0", "bond1", "vlan0".
- `stp` (Bool) - Enable Spanning Tree Protocol for bridge interfaces. STP prevents network loops by blocking redundant     paths and enables automatic failover when the primary path fails. Default: `True`
- `vlan_parent_interface` (String) - Parent interface for VLAN configuration.
- `vlan_pcp` (Int64) - Priority Code Point for VLAN traffic prioritization (0-7). Values 0-7 map to different QoS priority levels,     with 0 being lowest and 7 highest priority.
- `vlan_tag` (Int64) - VLAN tag number (1-4094).
- `xmit_hash_policy` (String) - Transmit hash policy for load balancing in link aggregation. LAYER2 uses MAC addresses, LAYER2+3 adds IP     addresses, and LAYER3+4 includes TCP/UDP ports for distribution. Default: `None` Valid values: `LAYER2`, `LAYER2+3`, `LAYER3+4`, `None`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_interface.example <id>
```
