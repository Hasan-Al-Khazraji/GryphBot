# GryphBot – GDSC Hacks 2025 Discord Assistant

  

## Table of Contents

- [Project Overview](#project-overview)
- [Feature Matrix](#feature-matrix)
- [Tech Stack](#tech-stack)
- [Prerequisites](#prerequisites)
- [Configuration](#configuration)
- [Local Development](#local-development)
- [Contributors](#contributors)
- [Roadmap](#roadmap)

## Project Overview

**GryphBot** is a Go‑powered Discord bot that acts as the on‑call AI helper for **GDSC Hacks 2025**, a 30‑hour in‑person hackathon hosted by the Google Developer Student Club at the University of Guelph. Backed by Google Gemini, the bot ingests official event documentation and answers participant questions in real time.

## Feature Matrix

| Feature       | Trigger                              | Response / Action                                                                       |
| ------------- | ------------------------------------ | --------------------------------------------------------------------------------------- |
| **Smart Q&A** | Any message that contains a question | Uses Gemini to generate a concise, authoritative answer sourced from the event handbook |

## Tech Stack

- **Language:** Go 1.24.0
- **Libraries:**
  - [`bwmarrin/discordgo`](https://github.com/bwmarrin/discordgo) – Discord Gateway & REST client
  - [`google/generative-ai-go`](https://github.com/google/generative-ai-go) – Gemini client SDK
- **Hosting:** Digital Ocean

## Prerequisites

- Go ≥ 1.24.0 installed and on your `$PATH`
- A Discord application with a **bot token** and the **Message Content** intent enabled
- A **Google Gemini API key** with generative‑AI billing enabled
- Git

## Configuration

Create a `.env` file in the project root:

```env
DISCORD_TOKEN=your_bot_token_here
GEMINI_API_KEY=your_gemini_key_here
```

## Local Development

```bash
git clone https://github.com/YOUR_ORG/gryphbot.git
cd gryphbot
cp .env.example .env   # then edit the values

# Run locally
go run main.go
```

## Contributors

- Hasan Al‑Khazraji
- Timothy Khan

## Roadmap

- Add special commands
- Add [context-caching](https://ai.google.dev/gemini-api/docs/caching?lang=go) so AI can be fed PDFs
- AI can learn from previous questions and cache answers

