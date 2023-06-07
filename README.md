# ges - Golang Event Store

Event Store with support for event sourcing in Go.

## Event Store Domain

* **Aggregate:** An entity constructed via event sourcing.
    * **Aggregate ID:** unique identifier for stream of all events related to a single aggregate. Used to append events to
      a stream, or retrieve events for that stream.
* **Append only:** property of storage, where new data can be appended to the storage, but existing data is immutable.
* **Event:** facts that have happened in the past, they are immutable.
* **Event ID:** TBD.
* **Event Sourcing:** The means by which the current state of an entity can be restored by replaying all its events.
* **Event Store:** is a type of database optimized for storage of events utilised for event sourcing.
* **Stream:** A sequence of events for the same aggregate.
* **Transaction ID:** TBD.
* **UUID:** universally unique identifier, used to provide strong identity where this is required.
* **Version:** TBD.

- [Event Sourcing - Key concepts](https://github.com/altairsix/eventsource#key-concepts)

## References

- [The Design of an Event Store](https://towardsdatascience.com/the-design-of-an-event-store-8c751c47db6f)
- [A Beginner’s Guide to Event Sourcing](https://www.eventstore.com/event-sourcing)
- [Essential features of an Event Store for Event Sourcing](https://itnext.io/essential-features-of-an-event-store-for-event-sourcing-13e61ca4d066)
- [Event Sourcing in Go](https://victoramartinez.com/posts/event-sourcing-in-go/)
- [Event Sourcing with PostgreSQL](https://github.com/evgeniy-khist/postgresql-event-sourcing)
- [Eventsource Go library](https://github.com/altairsix/eventsource)
