# TinyURL

Hello

## Design

```mermaid
sequenceDiagram
    Client-->>Backend: create short URL
    Backend-->>CFKV: write url
    Backend-->>CFKV: poll url
    Backend-->>Client: OK
```
