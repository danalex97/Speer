<sup> For an example skip to the end. </sup>

#### Description

Fixes #[issue number here]

Summary: [summary of the change and which issue is fixed here]

Changes: [specify the structures changed]

#### Testing

Tests:
- [test 1 description here]
- [test 2 description here]

#### Checklist:

Before requesting the PR, check you did the following:
- [ ] Tested the code(if applicable)
- [ ] Commented my code
- [ ] Changed the documentation(if applicable)

#### Dummy example

Fixes #22

**Summary:** I have added a new routing algorithm for the network layer.

**Changes:** I have changed the Router structure as follows:
  - different packet format
  - added a new helper structure *dummyRouter*

**Tests:**
  - test_router_runs_faster: comapres the new router with runtimes form the old version
  - test_router_small_topology: tests the router on a small dummy topology
  - test_router_ring: tests the router on a big(200 nodes) ring topology

**Checklist:**
- [x] Tested the code
- [x] Commented my code
- [ ] Changed the documentation - N/A
