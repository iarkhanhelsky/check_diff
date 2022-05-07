# check_diff (wip)

[![Coverage Status](https://coveralls.io/repos/github/iarkhanhelsky/check_diff/badge.svg)](https://coveralls.io/github/iarkhanhelsky/check_diff)

`check_diff` is a command-line tool targeted to apply static checks on changed
files and lines. 

## Install

### go install

```
go install github.com/iarkhanhelsky/check_diff@latest
```

### manually

Download the pre-compiled binaries from the [releases](https://github.com/iarkhanhelsky/check_diff/releases) 
page and copy them to the desired location.

## Setup

1. Create empty `check_diff.yaml` file in your project root directory.
2. Specify your linters configuration in `check_diff.yaml`
3. Change any of your source files introducing lint errors
4. Run the following command to check your changes
   ```
   $ git diff | check_diff 
   ```

### git hooks



## Builtin linter bindings

| Language | Linter                                              | Bundled Version |
|----------|-----------------------------------------------------|----------------|
| Java     | [Checkstyle](pkg/checker/java/checkstyle/Readme.md) | 9.3            |
| K8S      | [kube-linter](pkg/checker/k8s/kubelinter/Readme.md) | 0.2.5          |
| Ruby     | [rubocop](pkg/checker/ruby/rubocop/Readme.md)       |                |
 
## Output formats

* STDOUT - print lint issues in human-readable format
* Phabricator
* [Codeclimate](https://github.com/codeclimate/platform/blob/master/spec/analyzers/SPEC.md) 