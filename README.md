# Walkline [![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

### Simplistic sql database migration tool.

![Walkline](https://github.com/giwiro/walkline/raw/develop/resources/logo.png "Walkline")

#### Works on
Windows, macOS, Linux

#### Supported databases
PostgreSQL

## Usage
```
               _ _    _            
              | | |  | (_)           
__      __ _ _| | | _| |_ _ __   ___ 
\ \ /\ / / _` | | |/ | | | '_ \ / _ \
 \ V  V | (_| | |   <| | | | | |  __/
  \_/\_/ \__,_|_|_|\_|_|_|_| |_|\___|
        Simplistic sql database migration tool

Usage:
  walkline [command]

Available Commands:
  downgrade   Downgrades database n times
  generate    Generates sql revision based on the version ranged provided
  help        Help about any command
  history     A brief description of your command
  init        Initializes the version table in the default schema
  upgrade     Upgrades database to the target version

Flags:
  -h, --help          help for walkline
  -p, --path string   path of the migration files
  -u, --url string    sql database connection url
  -v, --verbose       add verbosity
      --version       version for walkline

```

## Configuration

You may optionally use a config file in order to automatically configure the flags. 
The file must be named `walkline.yaml` and has to be in the working directory.
For example:

```yaml
url: postgres://user:password@localhost/database?sslmode=disable
# The path flag can be relative (to the working directory) or absolute
path: ./relative/path/to/migration
schema: user_schema
verbose: false
```


## License
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.