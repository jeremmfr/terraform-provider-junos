# Contributing

## Reporting Bugs

* First, ensure that you're running the
[latest version](https://github.com/jeremmfr/terraform-provider-junos/releases)
of the provider. If you're using an older version, it's possible that the bug
has already been fixed.

* Next, check the GitHub
[issues list](https://github.com/jeremmfr/terraform-provider-junos/issues)
to see if the bug you've found has already been reported. If you think you may
be experiencing a reported issue that hasn't already been resolved, please
click "add a reaction" in the top right corner of the issue and add a thumbs
up (+1). You might also want to add a comment describing how it's affecting your
use.

* When submitting an issue, please be as descriptive as possible. Be sure to
provide all information necessary like :
  * The Terraform version
  * The Junos device type and version
  * Any error messages generated

## Feature Requests

* First, check the GitHub
[issues list](https://github.com/jeremmfr/terraform-provider-junos/issues)
to see if the feature you're requesting is already listed. (Be sure to search
closed issues as well, since some feature requests have been rejected.) If the
feature you'd like to see has already been requested and is open, click "add a
reaction" in the top right corner of the issue and add a thumbs up (+1). Feel
free to add a comment with any additional justification for the feature.

* When submitting a feature request on GitHub, be sure to include all
information necessary like :
  * A detailed description of the proposed feature.
  * An example Junos configuration to describe the parameters involved.
  * If specific Junos device type is necessary.

## Submitting Pull Requests

* Be sure to open an issue **before** starting work on a pull request, and
discuss your idea with the provider maintainer before beginning work. This will
help prevent wasting time on something that might we might not be able to
implement. When suggesting a new feature, also make sure it won't conflict with
any work that's already in progress.

* Once you've opened or identified an issue you'd like to work on, ask that it
be **assigned to you** so that others are aware it's being worked on. Please
provide a method of contacting you to supervise the progress of the work.

* All major new functionality must include relevant
[acceptance tests](#acceptance-tests-testing-interactions-with-junos-devices).

* All code submissions should meet the following criteria (CI will enforce
these checks):

  * Golang syntax is valid
  * Golang tests pass when run with `go test`
  * Golang linter tests pass when run `golangci-lint run -c .golangci.yml`
  * Documentations for new feature is present.
  * [Acceptance tests](#acceptance-tests-testing-interactions-with-junos-devices)
     pass for new or changed resources.

## Acceptance Tests: Testing interactions with Junos devices

Terraform includes a framework for constructing acceptance tests that imitate
the execution of one or more steps of applying one or more configuration files,
allowing multiple scenarios to be tested. Terraform acceptance tests use real
Terraform configurations to exercise the code in real plan, apply, refresh, and
destroy life cycles. See more on terraform
[sdkv2/testing/acceptance-tests](https://www.terraform.io/plugin/sdkv2/testing/acceptance-tests)
and [framework/acctests](https://developer.hashicorp.com/terraform/plugin/framework/acctests)

Terraform requires an environment variable `TF_ACC` be set in order to run
acceptance tests. More environment variables is also necessary for configure the
provider. See [docs](https://registry.terraform.io/providers/jeremmfr/junos/latest/docs#argument-reference)

```shell
TF_ACC=1 go test -v ./...
```

To run specifically tests with switch models, add environment variable
`TESTACC_SWITCH`.

To run tests for specific resource, use `-run` args

```shell
TF_ACC=1 go test -v ./... -run TestAccJunos<ResourceName>_basic
```

## Commenting

Only comment on an issue if you are sharing a relevant idea or constructive
feedback. Do not comment on an issue just to show your support (give the top
post a 👍 instead) or ask for an ETA.

## English language

If you find an English mistake in documentations, please forgive the Frenchman
that I am and offer me a correction in a pull request.
