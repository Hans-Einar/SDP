# VP08 — Channel contracts and sequences

[Navigator](../../navigator.md)

Explicit scenario steps validated against permits, participation, mode and request/result correlation.

- [Scenario: RejectStaleModelChange — mode ModelEditing](VP08-RejectStaleModelChange.md)
- [Scenario: SaveModelCheckpoint — mode ModelEditing](VP08-SaveModelCheckpoint.md)

## Derived MessageSet

| Channel | Mode | Datagram | Sender | Receiver | Source IDs |
| --- | --- | --- | --- | --- | --- |
| ModelOperations | ModelEditing | ModelChangeAccepted | ModelCoordinator | ModelCommandAdapter | f2e1caa154fbdb18a3ddc5cbb6cf257b13c8bfd8c85bf831ecac718e8fde9195a, f5b64393f9f0a981496920beb60d85f42a6db8ae2bad2f4e97507343e5e49b90f, f7d40c3a3a8dcb5f3903e1ba111acdec6588118d8a175b64e4301d5293d0a14da, f9a370d8df6475de7eebcaa7ec646c1511cbafdcabd997770af8a982484f8626f, fe18ef94cac20e293dceaba162af72221217b1851075c68d2fa31c20e7350362a |
| ModelOperations | ModelEditing | ModelChangeRefused | ModelCoordinator | ModelCommandAdapter | f56b4fc46ebd6d142e1b05e046dbf5d791bcb419958b9305383869c99619632f6, f5c4c6e7fd517ea9a5da31901ad3dd48de628fdee0e929bf7fe74bb4da13e899b, f6ba015e57387acae03d6f08fe0c5d1678811936997fb454bddf4f90d2460835f, f7d40c3a3a8dcb5f3903e1ba111acdec6588118d8a175b64e4301d5293d0a14da, feebda698bf3fdce5ba263880da99dc691e3bdb4f1ef156177ea30e4227eb48ef |
| ModelOperations | ModelEditing | RequestModelChange | ModelCommandAdapter | ModelCoordinator | f180798f519b29a3beccc356e74304d2f64dadbea3634b1d80ac204db40bbc794, f359c15d28aa9dd771ec69ac9e8a3fd9317931d1c4a9708a7bc09085105c7cb91, f3b24a65b9a1aab492877b6c2d126e06ec4ba8334f40d805800213e75dd8afad1, f7d40c3a3a8dcb5f3903e1ba111acdec6588118d8a175b64e4301d5293d0a14da, fef5d29092e55a95ab9432118e149dafb08e62d75d409b298d4e02bc68538163d |

