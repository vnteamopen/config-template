# config-template

A tool for embedding a file's contents into the template file.

The tool reads your input template file, loads pattern `{{file "/file/want/embed"}}`, adds embedded content to the template file and writes to the output path.

![config-template](https://raw.githubusercontent.com/vnteamopen/config-template/main/config-template.png)

## Features

 - Load input file and replace `{{file "/path/to/another/file"}}` with content from `/path/to/another/file`
 - Support both abolute path `{{file "/path/to/another/file"}}` and relative path `{{file "./another/file"}}`
 - Support recursive embedded file. If embedded file content contains `{{file ""}}`, tool keep to load it.

## Installation

### From source

Download the source code and try with:

```
go build -o output/config-template
```

Use `config-template`

### Use from Docker

Pull the docker image from:

```
docker pull ghcr.io/vnteamopen/config-template:main
```

## Quickstart:

1. Create a file `person.yml` with following content:

person.yml
```yml
name: thuc
{{file "./bio.yml"}}
{{file "./secrets.yml"}}
```

2. Create 2 files `bio.yml` and `secrets.yml` same folder with `secrets.yml`

bio.yml
```yml
username: abc
password: xyz
```

secrets.yml
```yml
job: developer
interests: running
```

3.1. Run config-template

```bash
./config-template person.yml output.yml
```

or

```bash
config-template person.yml output.yml
```

3.2. Run config-template with docker

```bash
docker run --rm -it -v $(pwd):/files/ -w /files ghcr.io/vnteamopen/config-template:main /app/config-template ./person.yml ./output.yml
```

4. output.yml will be write with content

```yml
name: thuc
username: abc
password: xyz
job: developer
interests: running
```

## Features

1. Overwrite template file with `-w` flag
```bash
config-template -w person.yml
```

2. Provide both flag `-w` and `output.yml` will overwrite template file and write output file
```bash
config-template -w person.yml output.yml
```

3. Provide multiple outputs
```bash
config-template person.yml output1.yml output2.yml
```

4. Support output to stdout
```bash
config-template -out-screen person.yml
```

## Future

Check https://github.com/vnteamopen/config-template/issues
