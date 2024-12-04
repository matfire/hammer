---
title: "Introduction"
---

## What is Hammer?

Hammer is a small tool that allows you to respond to Github webhooks to (usually) update an application.

Let's say you have a Laravel application (and are not using Docker, Forge, Vapor or anything similar); when you release an update, you might need to run:
- run migrations
- install php dependencies
- install & rebuild node dependencies
- rechache all the settings
- restart workers
- and more...

You could probably just have an endpoint in your application and a bash script you could run to run all these commands.

**but you could also use hammer**

Hammer allows you to handle the update process for any application via webhooks.

### The update process

:::info
Hammer is not yet stable, so features may be addded later on :)
:::

Hammer is (at the time of writing this document) meant to run on releases (meaning a Github release). When set up correctly, Github will send a webhook containing the tag name to hammer. Hammer can then:
- figure out which project needs an update
- temporarily store specified files in a separate folder
- pull and checkout the specified tag
- copy over the files
- run the commands specified in the configuration file

Sounds interesting?

## How Can I Use It?

Check out the [getting started section](/getting-started) to find out how to use `hammer`.

## Are There Any Detailed References?

Yep, just head on over to the [reference section](/reference)
