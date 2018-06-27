# Eos PWA Specification

## Contents

1. Project overview
2. Server API [server.md]
3. Client API [client.md]
4. Technical details [tech.md]
5. Research, notable decisions [research.md]

## Project overview

Eos is a progressive web app (and, potentially in the future, a native application for Android) offering alternative support pathways for mental health difficulties such as depression or Autism spectrum disorder (ASD). Specifically, Eos features a 3-part support framework:

1. Mood tracking, identifying trends (such as those prevalent in seasonally affective disorder or identifying specific days which may be worse than others - potentially an indication for bullying)
2. Adaptive responses determined based on reported mood - such as providing resources and tools to aid with recovery, or praising a positive mood.
3. One-on-one support with fellow service users, following the concept that helping others is likely to trigger releases of dopamine, leading to an increase in mood - whilst also aiding the other individual.

Acknowledging the sensitive nature of the application given its close relation to mental health concerns, the following requirements of the Eos project work to ensure healthy, research-driven aid to those in need of the service:

- Forever free (cost) finance model, to avoid locking individuals behind a paywall;
- Free software development ideology (libre), to ensure that users have the right to not use a specific version of the service, and are entitled to verify how the service is using their data manually via the source code;
- AI-driven automated moderation of the 'one-to-one chat system', flagging up potentially harmful discourse without the self-censorship caused by active human-driven moderation.

The official Eos branch also adheres to a custom confidentiality agreement, which is limited via the reporting system (chats which are reported by a participant are automatically declassified as confidential).