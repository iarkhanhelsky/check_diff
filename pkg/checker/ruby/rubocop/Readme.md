# rubocop

Install rubocop:
```
touch Gemfile 
echo "gem 'rubocop'" >> Gemfile
bundle install
bundle binstsubs rubocop
```

Example config:
```
Rubocop:
  Enabled: true
  Command: bin/rubocop
  Config: .rubocop.yaml
```