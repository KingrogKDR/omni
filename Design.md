# Design document for Omni

Everything related to the design decisions made keeping in mind the tradeoffs and usecases.

## The network layer

- Wire format: Protobuf
  Doc explaining this choice: https://kingrogkdr.github.io/post.html?slug=omni-post-1

The short version: Protobuf enforces types for fields and removes ambiguity, which is clearly not wanted for a distribution environment
