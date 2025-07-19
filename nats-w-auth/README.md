### Goals

Explore NATS capabilities.

- [x] Setup embedded Nats
- [x] Basic web
- [ ] Setup Auth Callout w. OIDC (probably Google or GitHub)
   - Works by subscribing to JWT events from the auth service
- [x] Sessions w. NATS
- [ ] Users workspace in KV
- [ ] Store users in NATS KV w. tokens
- [ ] Verify GitHub token
- [ ] Try Datastar
- [ ] Try clustering
- [ ] Do we need websockets?

## Persistance

Use NATS KV

### Bucket for Users

User details, roles, etc

### Bucket for OAUTH

### References

- https://github.com/synadia-io/rethink_connectivity/blob/main/19-auth-callout
- https://github.com/alexedwards/scs
