<p align="center">
    <img src="https://assets.sverdejot.dev/geemail/icon.png"/>
</p>

# `geemail`

"...You have 6001 unread messages..." has come to an end.

A terminal-based utility written in Go to help you manage and reduce unread emails in your Gmail inbox programmatically.

---

## Table of Contents

- [About](#about)
- [Features](#features)
- [Installation](#installation)

---

## About

`geemail` is a lightweight TUI application built in Go. It connects with your Gmail account and provides tools to efficiently reduce the number of unread messages — ideal for power users who want to clean up their inbox without using the Gmail web UI.

Use cases include:
- Unsubscribing from mailing list with just one key-stroke
- Bulk operations like arhiving, deleting, and mark-as-spam
- Inspecting mailing list history
- Open mails from the terminal

---

## Features

- Terminal user interface for quick navigation
- Efficient batch operations on Gmail messages
- Configurable filter and action workflows
- Designed for large volumes of unread emails

---

## Installation

First of all, I cannot distribute an application embedding OAauth secrets into it, its not safe (although Google says so). Moreover, there's a per-project-limit on how many requests can be done to the Gmail API, so it's not realistic to have a single project on a public, OSS project. Thus, I'm leaving you [here](https://developers.google.com/workspace/gmail/api/quickstart/go#set_up_your_environment) the documentation about how to generate your own Google OAuth application so it's also safe for you to play around with this project. Once you have the application and `client_secret.json` generated:

1. Export the credentials as an env var

```bash
export GEEMAIL_API_CREDENTIALS='<content of the json file>'
```

2. Install `geemail`

```bash
curl -fSL https://assets.sverdejot.dev/geemail/install.sh | sh -
```

3. Run `geemail` in your terminal
