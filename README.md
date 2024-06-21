# Hammer

> easily execute commands on release

## Description

This binary reads a toml configuration file and can execute commands in response to a Github release webhook.

## Background

I started learning Laravel a while back and, while not that complicated, automaticing a deployement without using Docker or Kubernetes (or [Forge] and [Vapor]) was quite tedious. You'd need to:

- execute migrations
- cache all configs (optional)
- build frontend assets
- etc

I wanted a way to simplify this. A service that could receive a release notification from Github, pull that version and perform all the operations needed for a new deployment

> [!IMPORTANT]
> This project does not perform zero-downtime deployments; if you need something like this, I suggest looking at a kubernetes or docker cluster


