# GitHub Actions Badge

There is a public version of this deployed and free to use at [https://actions-badge.atrox.dev](https://actions-badge.atrox.dev).

## Your own

You can build your own badge in the playground [available here](https://actions-badge.atrox.dev).

## Routes

- `/`: playground
- `/<user>/<repo>/badge`: returns the [endpoint](https://shields.io/endpoint) for shields.io
- `/<user>/<repo>/goto`: redirects to the action

## Example
[![GitHub Actions](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Fatrox%2Fsync-dotenv%2Fbadge)](https://actions-badge.atrox.dev/atrox/sync-dotenv/goto)

```
[![GitHub Actions](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Fatrox%2Fsync-dotenv%2Fbadge)](https://actions-badge.atrox.dev/atrox/sync-dotenv/goto)
```

[![GitHub Actions](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Fatrox%2Fsync-dotenv%2Fbadge&style=flat-square)](https://actions-badge.atrox.dev/atrox/sync-dotenv/goto)

```
[![GitHub Actions](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Fatrox%2Fsync-dotenv%2Fbadge&style=flat-square)](https://actions-badge.atrox.dev/atrox/sync-dotenv/goto)
```

For example, you can see this badge in action at [atrox/sync-dotenv](https://github.com/atrox/sync-dotenv).

## Contributing
Everyone is encouraged to help improve this project. Here are a few ways you can help:

- [Report bugs](https://github.com/atrox/github-actions-badge/issues)
- Fix bugs and [submit pull requests](https://github.com/atrox/github-actions-badge/pulls)
- Write, clarify, or fix documentation
- Suggest or add new features
