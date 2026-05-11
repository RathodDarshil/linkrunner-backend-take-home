# Linkrunner: Backend Take-Home

## About Linkrunner

Linkrunner is a mobile attribution and analytics platform. Customers ship our SDK, and we tell them which ad campaign installed the user, what that user did in-app, and what they paid. We process billions of events a month with tight latency budgets, because attribution decisions feed back into ad-network bidders within seconds.

We've raised **$580K (INR 5 Cr)** from top early-stage VC firms including **Titan Capital** and **2am VC**, along with several angel investors.

## The Role

**Distributed Systems Engineer**, 1 to 3 YoE, Go. You'll work on event ingestion, attribution computation, and the real-time analytics infrastructure that sits behind it. **Location:** in-person, Bangalore (HSR).

## About this assignment

Designed for **2 to 3 hours**. AI tools (Claude, Cursor, Copilot, etc.) are allowed and expected; we use them ourselves. We weight your decision-making and trade-off articulation as heavily as the code. Strong submissions move to a 45-minute technical conversation with the CTO within 7 days.

## The Problem

Build a Go service that ingests events over HTTP and forwards them to a flaky downstream, with at-least-once delivery and per-user ordering.

**Requirements**

- HTTP endpoint `POST /events` accepts JSON:
    ```json
    { "id": "string", "user_id": "string", "event_type": "string", "timestamp": "RFC3339", "payload": {} }
    ```
- Persist every accepted event durably before acknowledging the producer.
- Forward each event to the downstream at `http://localhost:9000/events`.
- The downstream is unreliable: random 500s, slow responses, and timeouts. We provide it (`starter/chaos/receiver.go`).
- Events must eventually reach the downstream, **ordered per `user_id`**, with no duplicates beyond the guarantee you state.
- Graceful shutdown on SIGTERM/SIGINT: in-flight events must not be lost.

You pick the persistence layer (SQLite, Postgres, BoltDB, files, whatever you can defend). You should not need Kafka or Redis for this scope.

## Constraints

- Go 1.22+
- External dependencies are fine; justify any non-trivial ones.
- Persistence is your call; justify it.
- **Time budget: 2 to 3 hours.** If you blow past that, stop and document what's missing.

## What we evaluate

1. **Consistency**: does the code actually do what your `Decisions` answers say it does?
2. **Trade-off articulation**: can you name what you gave up and why?
3. **Taste**: simplicity, and what you chose _not_ to build.

## Required: Decisions section

Add a `## Decisions` section to your submission's README answering the six questions below. **We weight these answers as heavily as the code itself.** Short, specific, opinionated beats long and hedged.

1. What end-to-end delivery guarantee does your service provide? Walk through the exact failure window where it could be violated.
2. How do you preserve per-`user_id` ordering? What's the throughput cost of that choice, and when does it become a bottleneck?
3. If the downstream is down for 10 minutes, what does your service do? What if it's down for 10 hours?
4. Where does state live, and what happens to in-flight events if the process is `kill -9`'d?
5. What did you choose **not** to build, and why? What would you build first if you had another day?
6. If we 100x the traffic, what breaks first?

## Submission

- Create a **private** GitHub repo.
- Record a **short video (under 30 seconds)** showing the service working end to end: producer posts an event, it persists, it lands at the downstream in order. Loom or unlisted YouTube is fine; **the link must be publicly accessible** (no sign-in or "request access" wall).
- **Embed the video link in your submission's README** (near the top, under a `## Demo` heading or similar).
- Submit by filling out this form: **https://forms.gle/8prWCSxubhd62tek8**

## What happens next

Strong submissions get a 45-minute call with Darshil (CTO) within 7 days of submission. The call walks through your assignment: your decisions, what you'd change, follow-up questions on the trade-offs.

## FAQ

**Can I use Claude / Cursor / Copilot?** Yes, expected. Just make sure you understand and can defend every decision.

**Do I need to build a UI?** No.

**Can I use Postgres / Redis / Kafka?** Postgres, sure. Redis or Kafka is overkill at this scope; if you reach for them, justify it.

**Will you respond if I'm not selected?** Yes, within 10 days of submission.
