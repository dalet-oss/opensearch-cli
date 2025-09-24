

# TODO N1

Create migration command (should be called implicitly) for the app config:
Case:
As an engineer I want to support multiple authentication backends as the creds source:
 - sops
 - age
 - ansible-vault

However, it's possible to do so with defining creds for the user with {cmd,args},
another convenient way is to retrieve them right from the application [no external command required].
To do so, we need to understand the initial token source, example:
 if token comes from keyring the format for the token field will be:
```yaml
token: keyring-<UUID>
```
for the ansible-vault:
```yaml
token: vault-<username|passw keys>
```