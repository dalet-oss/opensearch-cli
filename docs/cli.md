# Short opensearch-cli guide
<!-- TOC -->
* [Short opensearch-cli guide](#short-opensearch-cli-guide)
  * [Commands](#commands)
  * [Config](#config)
  * [Specification](#specification)
    * [Cluster spec](#cluster-spec)
    * [User spec](#user-spec)
      * [Credential backends](#credential-backends)
        * [Keyring](#keyring)
          * [OS X](#os-x)
          * [Linux](#linux)
        * [Vault](#vault)
          * [Vault file](#vault-file)
    * [Context spec](#context-spec)
<!-- TOC -->
## Commands
The List of available CLI commands is listed in the [auto-generated guide](opensearch-cli.md)

## Config
Config file of the opensearch-cli is inspired by `kubectl`. 
Default location of the config file is `${HOME}/.dalet/oscli/config`.

## Specification
```yaml
apiVersion: v1
clusters: []
users: []
contexts: []
current: ""
```

### Cluster spec
The cluster specification:
```yaml
- name: local-opensearch
  params:
    server: http://localhost:9200
    skipTlsVerify: false
```
Brief explanation of config parameters:
- `name` - the name of the cluster entry (must be **unique** across the cluster array)
- `params` - connection parameters:
  -  `server` - url of the OpenSearch cluster
  - `skipTlsVerify` - skip TLS verification when interacting with the cluster

### User spec
The user configuration:
```yaml
- name: example-user
  user:
    token: ""
    vault: {}
```
Brief explanation of config parameters:
- `name` - the name of the user entry (must be **unique** across the user array)
- `token` - id of the token in the keyring (see the supported [credential backends](#keyring))
- `vault` - vault configuration (see the supported [credential backends](#vault))

#### Credential backends
At the moment we are supporting two credential backends:

| Backend | Needs extra config | Operation systems | Machine independent |
|---------|--------------------|-------------------|---------------------|
| Keyring | +                  | MacOs/Linux*      | -                   |
| Vault   | -                  | +                 | +                   |

##### Keyring
Keyring support is implemented with [zalando/go-keyring](github.com/zalando/go-keyring). 
An additional machine configuration might be required to use this backend.
Config example:
```yaml
#- name: example-user
#  user:
    token: 01996de3-d5c0-7b9c-94a6-29cdace2d8eb
```
Brief explanation of config parameters:
  - `token` - the id of the entry in the user's keyring

###### OS X
On the OS X keyring relies on the `/usr/bin/security` binary available by default in most cases.
Data stored and retrieved from the OS X Keychain.
###### Linux
On the Linux keyring is depending on [Secret Service](https://specifications.freedesktop.org/secret-service-spec/latest/) and
expectation that the default collection `login` **exists in the keyring**. 
If it doesn't exist, you have to create it through the keyring frontend program/cli.

##### Vault
Vault authentication backend is implemented with help [ansible-vault-go](https://github.com/sosedoff/ansible-vault-go) library.
In the config file it could be used in two different ways:
- by **embedding** vault string into the config
- by reading existing file on the FS
This backend doesn't require any additional configuration from the user and `ansible-vault` to be installed on the user system.

Config example:
```yaml
 # name: example-user
 # user:
    vault:
      vaultString: |-
        $ANSIBLE_VAULT;1.1;AES256
        38663936316661386566623734326166643262643165343665353962306538653561346438623161
        3133626365383134326638663665656566363732616139630a363765356632646331363030666339
        35386234666561313238316132373139313531313332656137323330643066373961643664656430
        3037333633343830360a353961393132663638643038353363616136303466303039363539616134
        32663861663230663834313532363231313066656531656365333935303063303237393164373937
        65353432393533323632636332613266373165303664363161313262633538303930656365343666
        63663766393631666361323138356435356330353365636133346235373465323135643364643661
        35313132383265633831
      file: ~/my-super-secret.vault.yml
      userKey: opensearch_admin_user
      passKey: opensearch_admin_password
```
Brief explanation of config parameters:
- `vaultString` - meaning read the vault content from the string
- `file` - location if the vault file
- `userKey` - the key to gather username from ***
- `passKey`- the key to gather password from ***
  *** - only YAML files supported at the moment

###### Vault file
At the moment we are supporting only the yaml files. Example:
```yaml
opensearch_admin_user: "superadmin"
opensearch_admin_password: "superadmin"
```

### Context spec
The context spec is used to connect cluster A to the user B during the commands execution:
```yaml
- name: example@local-opensearch
  cluster: local-opensearch
  user: example-user
```
Brief explanation of config parameters:
- `name` - the name of the context entry (must be **unique** across the context array)
- `cluster` - the name of the cluster entry
- `user` - the name of the user entry

The **active context** is defined in the `current` field of the config file.
