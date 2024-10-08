## lagoon update organization

Update an organization

```
lagoon update organization [flags]
```

### Options

```
      --description string         Description of the organization
      --environment-quota int      Environment quota for the organization
      --friendly-name string       Friendly name of the organization
      --group-quota int            Group quota for the organization
  -h, --help                       help for organization
      --notification-quota int     Notification quota for the organization
  -O, --organization-name string   Name of the organization to update
      --project-quota int          Project quota for the organization
      --route-quota int            Route quota for the organization
```

### Options inherited from parent commands

```
      --config-file string                Path to the config file to use (must be *.yml or *.yaml)
      --debug                             Enable debugging output (if supported)
  -e, --environment string                Specify an environment to use
      --force                             Force yes on prompts (if supported)
  -l, --lagoon string                     The Lagoon instance to interact with
      --no-header                         No header on table (if supported)
      --output-csv                        Output as CSV (if supported)
      --output-json                       Output as JSON (if supported)
      --pretty                            Make JSON pretty (if supported)
  -p, --project string                    Specify a project to use
      --skip-update-check                 Skip checking for updates
  -i, --ssh-key string                    Specify path to a specific SSH key to use for lagoon authentication
      --ssh-publickey string              Specify path to a specific SSH public key to use for lagoon authentication using ssh-agent.
                                          This will override any public key identities defined in configuration
      --strict-host-key-checking string   Similar to SSH StrictHostKeyChecking (accept-new, no, ignore) (default "accept-new")
  -v, --verbose                           Enable verbose output to stderr (if supported)
```

### SEE ALSO

* [lagoon update](lagoon_update.md)	 - Update a resource

