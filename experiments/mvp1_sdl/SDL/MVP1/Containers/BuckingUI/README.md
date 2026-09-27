# BuckingUI

BuckingUI is a separately declared Container in [System.design](../../System.design).
It owns operator projections, application context and typed operator intent routing.

- [Domain.design](Domain.design): owned Units and Functionalities.
- [Views.design](Views.design): application composition and presentation definitions.
- [SDUI screens](../../SDUI/BuckingUI/MainPage.sdui): application presentation assets.
- [Shared components](../../Shared/README.md): mechanisms reused without merging application state.
- [Binding guide](../../SDUI/README.md): intended integration and current prototype limits.

Concrete rendering remains behind the Presentation contract. These files describe
application behavior; they contain no Fyne widget constructors or React code.
The legacy experimental profile remains explicit; relocation does not make the
Container executable or settle its host/process placement.
