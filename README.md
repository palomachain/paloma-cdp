# Paloma Chain Data Provider

The Paloma Chain Data Provider is a collection of services that ingest, transform and serve data from the [Paloma Chain](https://palomachain.com). It's main intention is the aggregation of swap transactions on the network, with a versatile web API to satisfy the basic requirements of charting solutions, but the project can be extended to support more data sources and use cases in the future.

## APIs

The default REST API is available on `http://localhost:8011/docs` per default. The GraphQL API implementation is still present in part, but not enabled and not maintained at the moment.

## How to contribute

### Issues

Issues should be used to report problems with the solution, request a new feature, or to discuss potential changes before a PR is created. When you create a new Issue, a template will be loaded that will guide you through collecting and providing the information we need to investigate.

If you find an Issue that addresses the problem you're having, please add your own reproduction information to the existing issue rather than creating a new one. Adding a [reaction](https://github.blog/2016-03-10-add-reactions-to-pull-requests-issues-and-comments/) can also help be indicating to our maintainers that a particular problem is affecting more than just the reporter.

### Pull Requests

PRs are always welcome and can be a quick way to get your fix or improvement slated for the next release. In general, PRs should:

- Consist of [conventional](https://www.conventionalcommits.org/en/v1.0.0/) and [signed](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits) commits.
- Only fix/add the functionality in question **OR** address wide-spread whitespace/style issues, not both.
- Add unit or integration tests for fixed or changed functionality (if a test suite already exists).
- Address a single concern in the least number of changed lines as possible.
- Include documentation in the repo
- Be accompanied by a complete Pull Request template (loaded automatically when a PR is created).

For changes that address core functionality or would require breaking changes (e.g. a major release), it's best to open an Issue to discuss your proposal first. This is not required but can save time creating and reviewing changes.

In general, we follow the ["fork-and-pull" Git workflow](https://github.com/susam/gitpr)

1. Fork the repository to your own Github account
2. Clone the project to your machine
3. Create a branch locally with a succinct but descriptive name
4. Commit changes to the branch
5. Following any formatting and testing guidelines specific to this repo
6. Push changes to your fork
7. Open a PR in our repository and follow the PR template so that we can efficiently review the changes.

## Getting Help

We have active, helpful communities on Twitter and Telegram.

- [Twitter](https://twitter.com/paloma_chain)
- [Telegram](https://t.me/palomachain)
- [Discord](https://discord.gg/HtUvgxvh5N)
- [Forum](https://forum.palomachain.com/)
