# Semantic changes

[Overview](index.md)

- added SDPTool/unit/GovernanceApplication
- added SDPTool/unit/McpResponseAdapter
- added McpResponseAdapter owns FrameSessionMcpResponse. [TARGET Responsibilities.design line 69](sources/TARGET/Responsibilities.design)
- removed JsonInvocationAdapter owns AdaptSessionInvocation. [NOW Responsibilities.design line 61](sources/NOW/Responsibilities.design)
- added OperationRegistry owns CheckActionAttribution. [TARGET Responsibilities.design line 71](sources/TARGET/Responsibilities.design)
- added SdpToolHost contains GovernanceApplication. [TARGET Responsibilities.design line 81](sources/TARGET/Responsibilities.design)
- added GovernanceApplication owns DispatchGovernanceOperation. [TARGET Responsibilities.design line 56](sources/TARGET/Responsibilities.design)
- added McpResponseAdapter owns FrameGeneralMcpResponse. [TARGET Responsibilities.design line 68](sources/TARGET/Responsibilities.design)
- removed CommandLineAdapter owns DispatchGovernanceOperation. [NOW Responsibilities.design line 50](sources/NOW/Responsibilities.design)
- removed JsonInvocationAdapter owns CheckActionAttribution. [NOW Responsibilities.design line 62](sources/NOW/Responsibilities.design)
- added GovernanceApplication owns AdaptSessionInvocation. [TARGET Responsibilities.design line 55](sources/TARGET/Responsibilities.design)
- added McpProcess contains McpResponseAdapter. [TARGET Responsibilities.design line 67](sources/TARGET/Responsibilities.design)
- removed McpAdapter owns FrameGeneralMcpResponse. [NOW Responsibilities.design line 65](sources/NOW/Responsibilities.design)
- removed McpAdapter owns FrameSessionMcpResponse. [NOW Responsibilities.design line 66](sources/NOW/Responsibilities.design)

```mermaid
flowchart LR
    f11["GovernanceApplication owns DispatchGovernanceOperation. (CHANGE added)"]
    f26["McpResponseAdapter owns FrameGeneralMcpResponse. (CHANGE added)"]
    f27["CommandLineAdapter owns DispatchGovernanceOperation. (CHANGE removed)"]
    f28["JsonInvocationAdapter owns CheckActionAttribution. (CHANGE removed)"]
    f3["McpResponseAdapter owns FrameSessionMcpResponse. (CHANGE added)"]
    f35["GovernanceApplication owns AdaptSessionInvocation. (CHANGE added)"]
    f38["McpProcess contains McpResponseAdapter. (CHANGE added)"]
    f4["JsonInvocationAdapter owns AdaptSessionInvocation. (CHANGE removed)"]
    f40["McpAdapter owns FrameGeneralMcpResponse. (CHANGE removed)"]
    f44["McpAdapter owns FrameSessionMcpResponse. (CHANGE removed)"]
    f5["OperationRegistry owns CheckActionAttribution. (CHANGE added)"]
    f8["SdpToolHost contains GovernanceApplication. (CHANGE added)"]
    n0["McpProcess (container CONTEXT)"]
    n1["SdpToolHost (container CONTEXT)"]
    n11["DispatchGovernanceOperation (functionality CONTEXT)"]
    n13["FrameGeneralMcpResponse (functionality CONTEXT)"]
    n14["FrameSessionMcpResponse (functionality CONTEXT)"]
    n2["AdaptSessionInvocation (functionality CONTEXT)"]
    n28["CommandLineAdapter (unit CONTEXT)"]
    n29["GovernanceApplication (unit CHANGE added)"]
    n32["JsonInvocationAdapter (unit CONTEXT)"]
    n33["McpAdapter (unit CONTEXT)"]
    n34["McpResponseAdapter (unit CHANGE added)"]
    n36["OperationRegistry (unit CONTEXT)"]
    n6["CheckActionAttribution (functionality CONTEXT)"]
    n14 -->|operand| f3
    n34 -->|operand| f3
    n2 -->|operand| f4
    n32 -->|operand| f4
    n6 -->|operand| f5
    n36 -->|operand| f5
    n1 -->|operand| f8
    n29 -->|operand| f8
    n11 -->|operand| f11
    n29 -->|operand| f11
    n13 -->|operand| f26
    n34 -->|operand| f26
    n11 -->|operand| f27
    n28 -->|operand| f27
    n6 -->|operand| f28
    n32 -->|operand| f28
    n2 -->|operand| f35
    n29 -->|operand| f35
    n0 -->|operand| f38
    n34 -->|operand| f38
    n13 -->|operand| f40
    n33 -->|operand| f40
    n14 -->|operand| f44
    n33 -->|operand| f44
```
