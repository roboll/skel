# skel

project templates from github releases and golang templates

## about

_skel_ uses golang templates to template skeleton projects. Projects can be sourced from github release archives in `.tar.gz` format, or from local directories.

_skel_ applies golang template engine (with modified delimiters `{{{` and `}}}` to avoid collission with templating templates), with default data from skel.yaml in the template, and overriding with values from the `data.yaml` (or a file named in the `--data` arg) in the current directory. _skel_ will open your `$EDITOR` with the resulting data for last minute editing before applying the template (disable with `--open-editor=false`).

## usage

When dealing with a private repository, `$GH_TOKEN` is required.

_skel_ has two modes of operation; pulling a github release artifact to use as a template, or using a source directory. When pulling from github, the following options are relevant.

```
--gh-owner "roboll"  github owner
--gh-repo "skel"     github repo
--gh-tag "latest"    release tag
--skel               template to pull - name of release artifact (no extension)
--name               name (replaces skel name in dir/file names, and in templates as {{{ .name }}})
--dest               template destination
```

When using a source directory, use `--src {path}`.

For full documentation, `skel --help`.

## templates

Templates use standard go templating, with modified delimiters to avoid collissions with existing templating. (`{{{` and `}}}`) In addition, the name of the skel in any file or directory names will be changed to the `--name` argument. (i.e. in a skel named template, using `--name mytempllate`, a directory named template will end up named mytemplate).

## simple example

See the [go](go) skel for an example.

```
go get github.com/roboll/skel
skel --skel go --name mynewproject
```
