# Hashicorp Vault 1Password Secrets Backend

This is a backend plugin to be used with [Hashicorp Vault](https://www.github.com/hashicorp/vault). This plugin allows for the retrieval, creation, and deletion of items stored in a 1Password vault accessed by use of the 1Password Connect.

## Installation

Prerequisites:

- A basic understanding of [Hashicorp Vault](https://www.hashicorp.com/products/vault). Read this guide to [get started with Vault](https://learn.hashicorp.com/tutorials/vault/getting-started-install).
- [Hashicorp Vault Command Line](https://www.vaultproject.io/docs/install) Installed
- A running Vault server (However, for a Quickstart this guide does have instructions on running a `local` server with the 1Password secrets backend).
- 1Password Connect deployed in your infrastructure.

The first step to using this plugin is to build it. Copy the binary to the `plugin directory` of your choice. For example:

```bash
$ go build -o vault/plugins/op-connect .
```

This directory will be specified as the `plugin_directory` in the Vault config used to start the server.

```json
plugin_directory = "path/to/plugin/directory"
```

Start a Vault server with this config file:

```bash
$ vault server -config=path/to/config.json ...
```

Register the plugin in in the Vault Server's plugin catalog:

```bash
$ vault write sys/plugins/catalog/secret/op-connect \
sha_256="$(shasum -a 256 path/to/plugin/directory/op-connect | cut -d " " -f1)" \
command="op-connect"
```

Enable the plugin:

```bash
$ vault secrets enable --plugin-name='op-connect' --path="op" plugin
...

Successfully enabled 'op-connect' at 'op'!
```

To check if your plugin has been registered you should be able to see the plugin when listing all registered plugins:

```bash
$ vault secrets list
```

### Quickstart Local Installation

The first step to using this plugin is to build it. Copy the binary to your `plugin directory` of your choice. For example:

```bash
$ go build -o vault/plugins/op-connect .
```

Run the vault server locally with a 1Password Connect plugin:

```bash
$ vault server -dev -dev-root-token-id=root -dev-plugin-dir=./vault/plugins -log-level=debug
```

Enable the plugin:

```bash
$ vault secrets enable op-connect
```

To check if your plugin has been registered you should be able to see the plugin when listing all registered plugins:

```bash
$ vault secrets list
```

### Plugin Configuration

In order to configure your plugin to access the 1Password Connect API, create a configuration json file:

```json
{
    "op_connect_host": "<host_address_of_1Password_Connect_API>",
    "op_connect_token": "<API_token_for_1Password_Connect>"
}
```

Save the configuration file:

```bash
$ vault write op/config @op-connect-config.json
```

## Usage

### Environment Variables

- **OP_CONNECT_TOKEN** (required if `op_connect_token` is not set in configuration): The API token created to be used to connect with the 1Password Connect API.

### Commands

**Listing vaults** available to the 1Password API token:

```bash
$ vault list op/vaults
```

**Listing items** stored in the specified vault:

```bash
$ vault list op/vaults/<vault_id>
```

**Read item**:

```bash
$ vault read op/vaults/<vault_id>/items/<item_id>
```

**Create item** (Please see the Creating and Updating Items section for more details on the json file contents):

```bash
$ vault write op/vaults/<vault_id>/items @some_json_file.json
```

**Update item** (Please see the Creating and Updating Items section for more details on the json file contents):

```bash
$ vault write op/vaults/<vault_id>/items/<item_id> @some_json_file.json
```

**Delete item**:

```bash
$ vault write op/vaults/<vault_id>/items/<item_id>
```

### **Creating and Updating Items Details**

- **category**(required): Describes the category of the item to create. Currently supported are `database`, `login`, and `password`.
- **title**(required on create): Specifies what the item will be titled.
- **url**: Specifies the url where the item may be used
- **fields**: Describes the fields to create for the item. Each field can be described with the following
    - **id**: The id of the field to create.
    - **label**: How the field will be titled in the UI
    - **type**: The type of the field. `STRING`, `EMAIL`, `CONCEALED`, `URL`, `TOTP`, `DATE`, `MONTH_YEAR`, and `MENU` are currently supported.
    - **purpose**: The purpose of the field. `""`, `USERNAME`, `PASSWORD`, `NOTES` are currently supported.
    - **value**: The value of the field.
    - **generate**: Used for fields with a password type. Set as true to have 1Password generate the value.
    - **entropy**: Used for fields with a password type. Set as an integer value for passwords where you would like to specify the value.
    - **section**: Describes what section to place the field. If not specified will be placed in the default section. Sections can be described with the following:
    - **section_id**: The id of the section to create the item in
- sections: Describes what sections to create for the item
    - id: The id of the section
    - label: How the section will be titled in the UI

Example Login Item with custom section:

```json
{
  "category": "login",
  "title": "Example Login",
  "fields": [
  {
    "id": "username",
    "label": "username",
    "type": "STRING",
    "purpose": "USERNAME",
    "value": "my_user"
  },
  {
    "id": "password",
    "label": "password",
    "purpose": "PASSWORD",
    "type": "CONCEALED",
    "value": "",
    "generate": true
  },
  {
    "id": "custom_field_id",
    "type": "STRING",
    "label": "My Custom Field",
    "value": "my custom value",
    "section": {
      "id": "my_new_section"
    }
  },
  ],
  "sections": [
  {
    "id": "my_new_section",
    "label": "New Section"
  }
  ]
}
```

Example Password Item:

```json
{
  "category": "password",
  "title": "Example Password",
  "fields": [
    {
      "id": "password",
      "label": "password",
      "purpose": "PASSWORD",
      "type": "CONCEALED",
      "value": "",
      "generate": true
    }
  ]
}
```

Example Database Item:

```json
{
  "category": "database",
  "title": "Example Database",
  "fields": [
    {
      "id": "username",
      "label": "username",
      "type": "STRING",
      "purpose": "USERNAME",
      "value": "my_user"
    },
    {
      "id": "password",
      "label": "password",
      "purpose": "PASSWORD",
      "type": "CONCEALED",
      "value": "",
      "generate": true
    },
    {
      "id": "hostname",
      "label": "hostname",
      "type": "STRING",
      "value": "my_host"
    },
    {
      "id": "database",
      "label": "database",
      "type": "STRING",
      "value": "my_database"
    },
    {
      "id": "port",
      "label": "port",
      "type": "STRING",
      "value": "8080"
    }
  ]
}
```

## Development

### Running Tests

```bash
$ go test -v ./... -cover
```

## Security

1Password requests you practice responsible disclosure if you discover a vulnerability. 

Please file requests via [**BugCrowd**](https://bugcrowd.com/agilebits). 

For information about security practices, please visit our [Security homepage](https://bugcrowd.com/agilebits).