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
- [Usage](#usage)

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

Clone the repository and build the tool using Go.

```bash
# Clone the repository
git clone https://github.com/sverdejot/geemail.git
cd geemail

# Download dependencies and build
go mod download
go build ./cmd/geemail

## Usage

TBD
