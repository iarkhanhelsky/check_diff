# ScriptAndRegexp

Script and Regex linter allowing to invoke an existing lint engine that is not 
integrated with check_diff.

The Script and Regex linter is a simple glue linter which runs some script on 
each path, and then uses a regex to parse lint messages from the script's output. 
(This linter uses a script and a regex to interpret the results of some real 
linter, it does not itself lint both scripts and regexes).

```yaml
ScriptAndRegexp/test:
  Enabled: true
  # The script will be invoked once for each file that is to be linted, with the 
  # file passed as the first argument. 
  # The script should emit lint messages to stdout, which will be parsed with 
  # the provided regex.
  # The return code of the script must be 0, or an exception will be raised 
  # reporting that the linter failed.
  # Multiple instances of the script may be run in parallel if there are 
  # multiple files to be linted, so they should not use any unique resources. 
  # For instance, this configuration would not work properly, because several 
  # processes may attempt to write to the file at the same time
  Script: bin/lint.sh
  # The regex must be a valid Golang regex, including delimiters and flags.
  # The regex will be matched against each line of the script output.
  Regexp: '(?P<file>[^:]*):(?P<line>[0-9]*):(?P<message>.*)$'
  # At least Include should be specified implicitly to match needed files
  Include:
    - 'snr/*.txt'
```