package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceGoogleComputeInstanceGroupManager() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceComputeInstanceGroupManagerRead,
		Schema: map[string]*schema.Schema{
			"base_instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The base instance name to use for instances in this group. The value must be a valid RFC1035 name. Supported characters are lowercase letters, numbers, and hyphens (-). Instances are named by appending a hyphen and a random four-character string to the base instance name.`,
			},

			"version": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Application versions managed by this instance group. Each version deals with a specific instance template, allowing canary release scenarios.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Version name.`,
						},

						"instance_template": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: compareSelfLinkRelativePaths,
							Description:      `The full URL to an instance template from which all new instances of this version will be created.`,
						},

						"target_size": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: `The number of instances calculated as a fixed number or a percentage depending on the settings.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fixed": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: `The number of instances which are managed for this version. Conflicts with percent.`,
									},

									"percent": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(0, 100),
										Description:  `The number of instances (calculated as percentage) which are managed for this version. Conflicts with fixed. Note that when using percent, rounding will be in favor of explicitly set target_size values; a managed instance group with 2 instances and 2 versions, one of which has a target_size.percent of 60 will create 2 instances of that version.`,
									},
								},
							},
						},
					},
				},
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the instance group manager. Must be 1-63 characters long and comply with RFC1035. Supported characters include lowercase letters, numbers, and hyphens.`,
			},

			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: `The region where the managed instance group resides.`,
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `An optional textual description of the instance group manager.`,
			},

			"fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The fingerprint of the instance group manager.`,
			},

			"instance_group": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The full URL of the instance group created by the manager.`,
			},

			"named_port": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: `The named port configuration.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the port.`,
						},

						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The port number.`,
						},
					},
				},
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},

			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URL of the created resource.`,
			},

			"target_pools": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         selfLinkRelativePathHash,
				Description: `The full URL of all target pools to which new instances in the group are added. Updating the target pools attribute does not affect existing instances.`,
			},
			"target_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: `The target number of running instances for this managed instance group. This value should always be explicitly set unless this resource is attached to an autoscaler, in which case it should never be set. Defaults to 0.`,
			},

			// If true, the resource will report ready only after no instances are being created.
			// This will not block future reads if instances are being recreated, and it respects
			// the "createNoRetry" parameter that's available for this resource.
			"wait_for_instances": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: `Whether to wait for all instances to be created/updated before returning. Note that if this is set to true and the operation does not succeed, Terraform will continue trying until it times out.`,
			},
			"wait_for_instances_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "STABLE",
				ValidateFunc: validation.StringInSlice([]string{"STABLE", "UPDATED"}, false),
				Description:  `When used with wait_for_instances specifies the status to wait for. When STABLE is specified this resource will wait until the instances are stable before returning. When UPDATED is set, it will wait for the version target to be reached and any per instance configs to be effective as well as all instances to be stable before returning.`,
			},

			"auto_healing_policies": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: `The autohealing policies for this managed instance group. You can specify only one value.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_check": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: compareSelfLinkRelativePaths,
							Description:      `The health check resource that signals autohealing.`,
						},

						"initial_delay_sec": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 3600),
							Description:  `The number of seconds that the managed instance group waits before it applies autohealing policies to new instances or recently recreated instances. Between 0 and 3600.`,
						},
					},
				},
			},

			"distribution_policy_zones": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The distribution policy for this managed instance group. You can specify one or more values.`,
				Set:         hashZoneFromSelfLinkOrResourceName,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					DiffSuppressFunc: compareSelfLinkOrResourceName,
				},
			},

			"distribution_policy_target_shape": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The shape to which the group converges either proactively or on resize events (depending on the value set in updatePolicy.instanceRedistributionType).`,
			},

			"update_policy": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				MaxItems:    1,
				Description: `The update policy for this managed instance group.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"minimal_action": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"RESTART", "REPLACE"}, false),
							Description:  `Minimal action to be taken on an instance. You can specify either RESTART to restart existing instances or REPLACE to delete and create new instances from the target template. If you specify a RESTART, the Updater will attempt to perform that action only. However, if the Updater determines that the minimal action you specify is not enough to perform the update, it might perform a more disruptive action.`,
						},

						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"OPPORTUNISTIC", "PROACTIVE"}, false),
							Description:  `The type of update process. You can specify either PROACTIVE so that the instance group manager proactively executes actions in order to bring instances to their target versions or OPPORTUNISTIC so that no action is proactively executed but the update will be performed as part of other actions (for example, resizes or recreateInstances calls).`,
						},

						"max_surge_fixed": {
							Type:          schema.TypeInt,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"update_policy.0.max_surge_percent"},
							Description:   `The maximum number of instances that can be created above the specified targetSize during the update process. Conflicts with max_surge_percent. It has to be either 0 or at least equal to the number of zones. If fixed values are used, at least one of max_unavailable_fixed or max_surge_fixed must be greater than 0.`,
						},

						"max_surge_percent": {
							Type:          schema.TypeInt,
							Optional:      true,
							ConflictsWith: []string{"update_policy.0.max_surge_fixed"},
							Description:   `The maximum number of instances(calculated as percentage) that can be created above the specified targetSize during the update process. Conflicts with max_surge_fixed. Percent value is only allowed for regional managed instance groups with size at least 10.`,
							ValidateFunc:  validation.IntBetween(0, 100),
						},

						"max_unavailable_fixed": {
							Type:          schema.TypeInt,
							Optional:      true,
							Computed:      true,
							Description:   `The maximum number of instances that can be unavailable during the update process. Conflicts with max_unavailable_percent. It has to be either 0 or at least equal to the number of zones. If fixed values are used, at least one of max_unavailable_fixed or max_surge_fixed must be greater than 0.`,
							ConflictsWith: []string{"update_policy.0.max_unavailable_percent"},
						},

						"max_unavailable_percent": {
							Type:          schema.TypeInt,
							Optional:      true,
							ConflictsWith: []string{"update_policy.0.max_unavailable_fixed"},
							ValidateFunc:  validation.IntBetween(0, 100),
							Description:   `The maximum number of instances(calculated as percentage) that can be unavailable during the update process. Conflicts with max_unavailable_fixed. Percent value is only allowed for regional managed instance groups with size at least 10.`,
						},

						"instance_redistribution_type": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateFunc:     validation.StringInSlice([]string{"PROACTIVE", "NONE", ""}, false),
							DiffSuppressFunc: emptyOrDefaultStringSuppress("PROACTIVE"),
							Description:      `The instance redistribution policy for regional managed instance groups. Valid values are: "PROACTIVE", "NONE". If PROACTIVE (default), the group attempts to maintain an even distribution of VM instances across zones in the region. If NONE, proactive redistribution is disabled.`,
						},
						"replacement_method": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateFunc:     validation.StringInSlice([]string{"RECREATE", "SUBSTITUTE", ""}, false),
							DiffSuppressFunc: emptyOrDefaultStringSuppress("SUBSTITUTE"),
							Description:      `The instance replacement method for regional managed instance groups. Valid values are: "RECREATE", "SUBSTITUTE". If SUBSTITUTE (default), the group replaces VM instances with new instances that have randomly generated names. If RECREATE, instance names are preserved.  You must also set max_unavailable_fixed or max_unavailable_percent to be greater than 0.`,
						},
					},
				},
			},

			"stateful_disk": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: `Disks created on the instances that will be preserved on instance delete, update, etc. Structure is documented below. For more information see the official documentation. Proactive cross zone instance redistribution must be disabled before you can update stateful disks on existing instance group managers. This can be controlled via the update_policy.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The device name of the disk to be attached.`,
						},

						"delete_rule": {
							Type:         schema.TypeString,
							Default:      "NEVER",
							Optional:     true,
							Description:  `A value that prescribes what should happen to the stateful disk when the VM instance is deleted. The available options are NEVER and ON_PERMANENT_INSTANCE_DELETION. NEVER - detach the disk when the VM is deleted, but do not delete the disk. ON_PERMANENT_INSTANCE_DELETION will delete the stateful disk when the VM is permanently deleted from the instance group. The default is NEVER.`,
							ValidateFunc: validation.StringInSlice([]string{"NEVER", "ON_PERMANENT_INSTANCE_DELETION"}, true),
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The status of this managed instance group.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_stable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `A bit indicating whether the managed instance group is in a stable state. A stable state means that: none of the instances in the managed instance group is currently undergoing any type of change (for example, creation, restart, or deletion); no future changes are scheduled for instances in the managed instance group; and the managed instance group itself is not being modified.`,
						},

						"version_target": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `A status of consistency of Instances' versions with their target version specified by version field on Instance Group Manager.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_reached": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `A bit indicating whether version target has been reached in this managed instance group, i.e. all instances are in their target version. Instances' target version are specified by version field on Instance Group Manager.`,
									},
								},
							},
						},
						"stateful": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Stateful status of the given Instance Group Manager.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"has_stateful_config": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `A bit indicating whether the managed instance group has stateful configuration, that is, if you have configured any items in a stateful policy or in per-instance configs. The group might report that it has no stateful config even when there is still some preserved state on a managed instance, for example, if you have deleted all PICs but not yet applied those deletions.`,
									},
									"per_instance_configs": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `Status of per-instance configs on the instance.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"all_effective": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: `A bit indicating if all of the group's per-instance configs (listed in the output of a listPerInstanceConfigs API call) have status EFFECTIVE or there are no per-instance-configs.`,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"instance_template": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: `The instance template.`,
			},
		},
	}
}

func dataSourceComputeInstanceGroupManagerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, region, name, err := GetRegionalResourcePropertiesFromSelfLinkOrSchema(d, config)
	if err != nil {
		return err
	}

	instanceGroupManager, err := config.NewComputeClient(userAgent).RegionInstanceGroupManagers.Get(
		project, region, name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Region Instance Group Manager %q", name))
	}

	d.SetId(fmt.Sprintf("projects/%s/regions/%s/instanceGroupManagers/%s", project, region, name))
	if err := d.Set("self_link", instanceGroupManager.SelfLink); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}
	if err := d.Set("name", name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("region", region); err != nil {
		return fmt.Errorf("Error setting region: %s", err)
	}
	if err := d.Set("base_instance_name", instanceGroupManager.BaseInstanceName); err != nil {
		return fmt.Errorf("Error setting base_instance_name: %s", err)
	}
	if err := d.Set("instance_template", instanceGroupManager.InstanceTemplate); err != nil {
		return fmt.Errorf("Error setting base_instance_name: %s", err)
	}
	return nil

	// config := meta.(*Config)
	// if name, ok := d.GetOk("name"); ok {
	// 	zone, err := getZone(d, config)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	project, err := getProject(d, config)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	d.SetId(fmt.Sprintf("projects/%s/zones/%s/instanceGroupsManagers/%s", project, zone, name.(string)))
	// }

	// return resourceDataComputeInstanceGroupManagerRead(d, meta)
}
