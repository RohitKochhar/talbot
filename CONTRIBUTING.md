# Contributing

This document explains the process of contributing to the Talbot project. The contributing guidelines outlined here are inspired by those outlined in the [Thanos Project](https://github.com/thanos-io/thanos/blob/main/CONTRIBUTING.md)

## Talbot Philosophy

The philosophy of Talbot borrow heavily from UNIX philosophy and the Golang programming language

- Each function should do one thing and do it well
- Functions should be highly testable for a wide variety of test cases
- Code should be easy to read, improve and write

## Feedback / Issues

If you encounter any issue or you have an idea to improve, please:

- Check existing [GitHub Issues](https://github.com/RohitKochhar/talbot/issues). 
  - If you find a relevant topic, please comment on the issue
  - If you do not find a relevant issue, create a new GitHub issue providing any relevant information as suggested by the issue template.

## Adding New Features / Components

When contributing a complex change to the Talbot repository, please discuss the change you wish to make within a GitHub issue. 

## General Naming

In the code and documentation, prefer non-offensive terminology, for example:

- `allowlist`/`denylist` (instead of `whitelist`/`blacklist`)
- `primary`/`replica` (instead of `master`/`slave`)
- `openbox`/`closedbox` (instead of `whitebox`/`blackbox`)

## Commit Naming

Commits must follow the guidelines specified by the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)
