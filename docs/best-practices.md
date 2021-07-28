# Development Best Practices

## Stories should have at most 1 day of development work

Often when working on stories it may be forgotten that writing tests and documentation can take at least as long as writing the code itself. Therefore, we have a rule of thumb that stories should have at most 1 day of development work. In addition to that it will probably take 1-2 days to write tests and update the documentation, resulting in a total of 2-3 days from status Open to In Review.

In addition to the faster development cycle, it is much easier to do estimations and to reason about a change with 1 day of development work, than a change of 1 week of development work.

When a story is encountered that has a too large scope to be finished within 1 day, the story may be rewritten to a reduced scope, so it can be finished within 1 day. Follow up stories can then be logged to extend the functionality to the scope of the initial story.

## Test categories

https://martinfowler.com/articles/practical-test-pyramid.html

Within WoningFinder we use different test categories for different use cases. We distinguish the following test categories:

- Unit tests
- Component tests
- End-to-end tests
- Regression tests
- Performance tests

### Unit tests

https://martinfowler.com/bliki/UnitTest.html

Unit tests are used to test single unit of code. Within WoningFinder we strive for covering as many possible code paths using unit tests. Therefore, we also generate code coverage reports of our unit tests, to help PR reviewers validate that the expected code paths are covered.

Unit tests are run on every commit to the GIT repository and they must succeed before a PR can be merged.

### Component tests

https://martinfowler.com/bliki/ComponentTest.html

Component tests are used within WoningFinder to test a single API endpoint. We use the same bootstrap process as the actual code, but some of the dependencies are replaced with stubs, to prevent calls to external services. The component tests do not require the deployment of the code, so are relatively fast.

Component tests are run on every commit to the GIT repository and they must succeed before a PR can be merged.
