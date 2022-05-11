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

### Local
1. Create empty `check_diff.yaml` file in your project root directory.
2. Specify your linters configuration in `check_diff.yaml`
3. Change any of your source files introducing lint errors
4. Run the following command to check your changes
   ```
   $ git diff | check_diff 
   ```

### git hooks

### CI

#### Gitlab

See [Gitlab docs](https://docs.gitlab.com/ee/user/project/merge_requests/code_quality.html#implementing-a-custom-tool) for more information 

Example step configuration
```
check-diff:
  stage: test
  script:
    # Find merge base and make a diff. We don't need changes that appeared in
    # upstream after feature branch was created
    - git diff -r $(git merge-base $CI_MERGE_REQUEST_DIFF_BASE_SHA HEAD) | ./bin/check_diff --format gitlab -o .gitlab-lint
  artifacts:
    reports:
      codequality: .gitlab-lint
```

## Builtin linter bindings

| Language | Linter                                                      | Bundled Version | Tested With  |
|----------|-------------------------------------------------------------|-----------------|--------------|
| Go       | [golangci-lint](pkg/checker/golang/golangci-lint/Readme.md) | 1.46.0          | -//-         |
| Java     | [Checkstyle](pkg/checker/java/checkstyle/Readme.md)         | 9.3             | -//-         |
| K8S      | [kube-linter](pkg/checker/k8s/kubelinter/Readme.md)         | 0.2.5           | -//-         |
| Ruby     | [rubocop](pkg/checker/ruby/rubocop/Readme.md)               |                 | 1.25.1       |
 
## Output formats

* STDOUT - print lint issues in human-readable format
* Phabricator
* [Codeclimate](https://github.com/codeclimate/platform/blob/master/spec/analyzers/SPEC.md) 