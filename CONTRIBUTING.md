# Contributing

We love pull requests from everyone. By participating in this project, as an individual or an organization, you are encouraged to submit patches and other contributions to improve this project by following our guidelines:

- [Code of Conduct](#code-of-conduct)
- [Issues and Bugs](#issues-and-bugs)
- [Feature Requests](#feature-requests)
- [Submission Guidelines](#submission-guidelines)
- [Code Style](#code-style)
- [Commit Message Guidelines](#commit-message-guidelines)

## Code of Conduct

This project and everyone participating in it is governed by the [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Issues and Bugs

If you find a bug in the source code or a mistake in the documentation, you can help us by submitting an issue following the issue template. If it's a minor bug without breaking changes, you can submit a Pull Request with a quick fix.

## Feature Requests

You can request a new feature by submitting an issue following the issue template. If you would like to implement a new feature, please submit an issue with a proposal for your work first, to be sure that we can use it.

## Submission Guidelines

### Submitting an Issue

Before you submit your issue, search the archive to check if a similar issue has already been filed. If there is a similar issue, you can contribute to it by adding a comment with your use case, example or additional information.

### Submitting a Pull Request

Before you submit your pull request, check that it meets these guidelines:

1. Search the archive to check if a similar pull request has already been submitted.
2. Fork the repository.
3. Make your changes in a new git branch:

   ```shell
   git checkout -b my-fix-branch master
   ```

4. Create your patch, **including appropriate test cases**.
5. Follow our [Code Style](#code-style).
6. Run the full test suite, and ensure that all tests pass.
7. Commit your changes using a descriptive commit message that follows our [commit message conventions](#commit-message-guidelines). Adherence to these conventions is necessary because release notes are automatically generated from these messages.
8. Push your branch to GitHub:

   ```shell
   git push origin my-fix-branch
   ```
9. In GitHub, send a pull request to `main`.
    - If the pull request adds functionality, it should include tests and documentation.
    - If the pull request changes existing functionality, it should include tests and documentation.
    - If the pull request fixes an issue, it should include the issue number in the commit message.
    - If the pull request includes breaking changes, it should include a migration guide.

## Code Style

We adhere to the [Google Style Guides](https://google.github.io/styleguide/) for all our code. We strongly recommend that you use the [Prettier](https://prettier.io/) code formatter to ensure that all code is formatted according to the style guides.

Keep in mind that Prettier is not a silver bullet. You should also follow the style guides and use your best judgment.

Try to keep both cyclomatic complexity and cognitive complexity low. You can use [SonarCloud](https://sonarcloud.io/) or [CodeClimate](https://codeclimate.com/) to check your code.

## Commit Message Guidelines

We have very precise rules over how our git commit messages can be formatted. This leads to **more readable messages** that are easy to follow when looking through the **project history**. But also, we use the git commit messages to **generate the changelog**.

### Commit Message Format

Each commit message consists of a **header**, a **body** and a **footer**. The header has a special format that includes a **type**, a **scope** and a **subject**:

```text
<type>(<scope>): <subject>
<BLANK LINE>
<body>
<BLANK LINE>
<footer>
```

The **header** is mandatory and the **scope** of the header is optional.

Any line of the commit message cannot be longer than 100 characters. This allows the message to be easier to read on GitHub as well as in various git tools.

The **type** is contained within the title and can be one of these types:

- **build**: Changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm)
- **chore**: Changes to auxiliary tools and libraries such as documentation generation or linters (example scopes: eslint, prettier, husky)
- **ci**: Changes to our CI configuration files and scripts (example scopes: Travis, Circle, GitHub Actions)
- **docs**: Documentation only changes
- **feat**: A new feature
- **!feat**: A breaking change to a feature
- **fix**: A bug fix
- **perf**: A code change that improves performance
- **refactor**: A code change that neither fixes a bug nor adds a feature
- **style**: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
- **test**: Adding missing tests or correcting existing tests

The **scope** should be the name of the npm package affected (as perceived by the person reading the changelog generated from commit messages).

The **subject** contains a succinct description of the change:

- use the imperative, present tense: "change" not "changed" nor "changes"
- don't capitalize the first letter
- no dot (.) at the end

The **body** should include the motivation for the change and contrast this with previous behavior.

The **footer** should contain any information about **Breaking Changes** and is also the place to reference GitHub issues that this commit **Closes**.

The **BREAKING CHANGE** section should always be the first section in the footer and it should contain a description of the breaking changes and any migration notes required.

Commit messages should be in compliance with gitlint rules. If you are working on our provided Dev Container, you can run `gitlint` to check your commit messages.
