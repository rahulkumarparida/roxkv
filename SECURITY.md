# Security Policy

Security is important for RoxKV because the project includes a network-accessible key-value database, persistence mechanisms, Pub/Sub functionality, monitoring capabilities, and RoxAI integrations.

If you discover a security vulnerability, please report it responsibly.

## Supported Versions

RoxKV is currently under active development.

Security fixes are generally applied to the latest version of the project.

| Version        | Supported   |
| -------------- | ----------- |
| Latest         | Yes         |
| Older versions | Best effort |

As the project matures and begins maintaining stable releases, this policy may be updated with explicit version support.

## Reporting a Vulnerability

Please **do not open a public GitHub issue for an unpatched security vulnerability**.

Public disclosure before a fix is available may put users or deployments at risk.

Instead, report the vulnerability privately through GitHub's private vulnerability reporting feature if it is enabled for this repository.

When submitting a report, include as much of the following information as possible:

* A description of the vulnerability
* The affected component
* Steps to reproduce the issue
* Expected behavior
* Actual behavior
* Potential security impact
* Relevant logs or error messages
* A minimal proof of concept, when appropriate
* Suggested mitigation or fix, if you have one

Please remove API keys, tokens, passwords, private user data, and other secrets from reports and logs.

## Security-Sensitive Areas

Security reports are particularly important for issues involving:

* TCP/network handling
* RESP parsing
* Authentication or authorization, if implemented
* Command execution
* Persistence and snapshot handling
* File access
* Path traversal
* Denial-of-service conditions
* Memory or resource exhaustion
* Concurrency issues with security implications
* Pub/Sub isolation
* Docker configuration
* Exposed services
* RoxAI tool execution
* LLM-controlled operations
* Prompt injection affecting privileged tools
* Ollama or external model communication
* Sensitive information appearing in logs
* Secrets or credentials being exposed

This list is not exhaustive.

## RoxAI Security

RoxAI introduces additional security considerations because LLM output is probabilistic and may be influenced by untrusted input.

RoxAI should not be treated as a security boundary.

Security-sensitive operations should be enforced by deterministic application code rather than relying on an LLM to decide whether an operation is permitted.

Contributors working on RoxAI should avoid designs where model-generated text can directly execute unrestricted:

* shell commands,
* filesystem operations,
* network requests,
* database operations,
* administrative operations,
* or other privileged actions.

Tool inputs should be validated before execution.

Where appropriate, tools should expose narrow, explicitly defined operations rather than unrestricted system access.

## Prompt Injection

Inputs processed by RoxAI may contain instructions designed to manipulate the model.

Treat:

* user input,
* database values,
* retrieved content,
* logs,
* external responses,
* and tool output

as potentially untrusted data.

An LLM instruction contained inside untrusted data should not automatically gain authority over RoxAI or its tools.

Authorization and security decisions must be enforced outside the language model.

## Network Exposure

RoxKV may expose services over TCP or HTTP.

Users are responsible for understanding which services are exposed outside their machine or Docker network.

Development configurations should not automatically be considered safe for public internet exposure.

Before exposing RoxKV publicly, consider:

* authentication,
* authorization,
* firewall rules,
* rate limiting,
* transport encryption,
* network isolation,
* resource limits,
* and monitoring.

If these protections are not implemented by RoxKV, they should be provided by the surrounding deployment infrastructure.

## Docker Security

Containerization improves deployment consistency but does not automatically make an application secure.

RoxKV containers should avoid:

* unnecessary privileged mode,
* unnecessary host networking,
* mounting sensitive host directories,
* exposing internal services publicly,
* embedding credentials inside images,
* running unnecessary services,
* and granting containers more permissions than required.

Secrets should be supplied through appropriate configuration mechanisms and must not be committed to the repository.

## Secrets

Never commit:

* passwords,
* API keys,
* access tokens,
* private keys,
* database credentials,
* cloud credentials,
* or other sensitive configuration.

Use environment variables or an appropriate secrets-management mechanism.

The repository should provide `.env.example` where configuration examples are required.

The real `.env` file should remain excluded from version control.

If a secret is accidentally committed, removing it from the latest commit is not sufficient.

The secret should be considered compromised and rotated.

## Dependencies

Dependencies can introduce vulnerabilities into the project.

Contributors should:

* avoid unnecessary dependencies,
* keep dependencies reasonably updated,
* review major dependency changes,
* and investigate known security advisories when relevant.

Automated dependency and security scanning may be introduced as the project matures.

## Denial of Service

Because RoxKV accepts network requests and maintains data in memory, malformed or excessive requests may create resource-exhaustion risks.

Security reports involving:

* unbounded allocations,
* extremely large RESP payloads,
* connection exhaustion,
* goroutine leaks,
* CPU exhaustion,
* uncontrolled Pub/Sub subscriptions,
* or other resource-exhaustion scenarios

are welcome.

## Responsible Disclosure

Please allow reasonable time for a vulnerability to be investigated and fixed before publishing technical details that could enable exploitation.

When appropriate, a security fix may include:

1. Reproducing and confirming the vulnerability.
2. Developing a fix.
3. Adding regression tests.
4. Preparing a security advisory.
5. Releasing the fix.
6. Publicly documenting the issue after users have had an opportunity to update.

## Security Research

Good-faith security research intended to improve RoxKV is welcome.

Please avoid:

* accessing data that does not belong to you,
* disrupting systems belonging to others,
* attempting attacks against deployments without permission,
* or publicly disclosing vulnerabilities before they can reasonably be addressed.

Test security issues against environments you own or have explicit permission to test.

## No Security Guarantee

RoxKV is an evolving open-source project and should not currently be assumed to provide the security guarantees of a mature production database.

Users evaluating RoxKV for sensitive or production workloads should perform their own security assessment and deploy appropriate external protections.

## Acknowledgements

Contributors who responsibly report security vulnerabilities may be acknowledged in release notes or security advisories, with their permission.

Thank you for helping make RoxKV safer.
