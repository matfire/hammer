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

> [!IMPORTANT]
> by default the programs are executed using the **sh** interpreter; be careful any program you need to run can be run using sh

## The config file

The config file is written in the TOML format and should be formatted as follows:

- each project should be under the group `apps`
- each project should have
  - name: the project name
  - path: the path the project resides in
  - commands: a list of commands to execute
  - secret; the github secret used to verify the webhook's authenticity

Here's an example:
```toml
[apps.example]
name="example"
path="/home/test/program"
commands=["pnpm ci", "pnpm build"]
```

## How to use this

- Create a config file (the program by default looks for a **config.toml** file in the same folder as the executable)
- Run the program (there is an example systemd service in the [deployment](#deployment) section)
- Enjoy :)

## Deployment

There are lots of different ways to run a program automatically; I personally prefer a simple systemd service, but it could be anything (like a tmux session). Here's an example of a systemd service you can use as a starting point for your own

```
[Unit]
Description=Hammer instance

[Service]
ExecStart=/path/to/hammer
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```
