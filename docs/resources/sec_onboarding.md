---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sccfm_sec_onboarding Resource - sccfm"
subcategory: ""
description: |-
  Use this resource to wait for an SEC to finish onboarding. When an SEC is onboarded, either manually or using the SCC Firewall Manager Terraform Modules for AWS https://github.com/CiscoDevNet/terraform-aws-cdo-sec, it can take a few minutes before the SEC is active and capable of proxying communications between SCC Firewall Manager and the device. This resource allows you to wait until this is done.
---

# sccfm_sec_onboarding (Resource)

Use this resource to wait for an SEC to finish onboarding. When an SEC is onboarded, either manually or using the SCC Firewall Manager Terraform Modules for [AWS](https://github.com/CiscoDevNet/terraform-aws-cdo-sec), it can take a few minutes before the SEC is active and capable of proxying communications between SCC Firewall Manager and the device. This resource allows you to wait until this is done.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Specify the name of the SEC.

### Read-Only

- `id` (String) The unique identifier of this SEC onboarding resource.
