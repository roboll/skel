# skel

project templates from github releases and golang templates

## about

_skel_ uses golang templates to template skeleton projects. Projects can be
sourced from github release archives in `.tar.gz` format, or from directories.

_skel_ applies golang template engine to the files in the template, using a
default skel.yaml in the template, and overriding with values from the
`data.yaml` (or `--data`) arg.

_skel_ does not address keeping projects in sync, or anything else beyond
initial bootstrap and templating.

## usage

When dealing with a private repository, `$GH_TOKEN` or `--gh-token` is required.

_skel_ has two modes of operation; pulling a github release artifact to use as
a template, or using a source directory.

When pulling from github, the following options are relevant.

```
--gh-owner "roboll"  github owner
--gh-repo "skel"     github repo
--gh-tag "latest"    release tag
--skel               template to pull - name of release artifact (no extension)
--name               name (replaces `skel` in dir names, and in templates as .name)
--dest               path to place project
```

When using a source directory, use `--src "{path}"`.

For full documentation, `skel --help`.

## templates

Templates use standard go templating, with values from a `map[string]string`
unmarshalled from yaml files. In addition, the word `skel` in any directory
name will be replaced by the value of the `--name` arg.

## simple example

See the [go](go) skel for an example.

```
go get github.com/roboll/skel
skel --skel go --name mynewproject --dest mynewproject
```
