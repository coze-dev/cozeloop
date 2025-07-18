<div align="center">
<h1>CozeLoop Open-source Edition</h1>
<p><strong> Platform-level solution for developing and operating AI agent</strong></p>
<p>
<a href="#What-can-CozeLoop-do">CozeLoop</a> •
<a href="#Feature-list">Feature list</a> •
<a href="#Quickstart">Quick start</a> •
<a href="#Development-guide">Development guide</a>
</p>
<p>
  <img alt="License" src="https://img.shields.io/badge/license-apache2.0-blue.svg">
  <img alt="Go Version" src="https://img.shields.io/badge/go-%3E%3D%201.23.4-blue">
</p>

English | [中文](README.cn.md)

</div>

## What is CozeLoop

[CozeLoop ](https://www.coze.cn/loop) (CozeLoop) is a developer-oriented, platform-level solution focused on the development and operation of AI agents. It addresses various challenges faced during the AI agent development process, providing full lifecycle management capabilities from development, debugging, evaluation, to monitoring.

Based on the commercial version, CozeLoop introduces a open-source edition that offers developers free access to core foundational feature modules. By sharing its core technology framework in an open-source model, developers can customize and extend according to business needs, facilitating community co-construction, sharing, and exchange, helping developers participate in AI agent exploration and practice with zero barriers.

## What can CozeLoop do?
CozeLoop helps developers efficiently develop and operate AI agents by providing full-lifecycle management capabilities. Whether it's prompt engineering, AI agent evaluation, or monitoring and optimization after deployment, CozeLoop offers powerful tools and intelligent support, significantly simplifying the AI agent development process and improving the performance and stability of AI agents.

* **Prompt development**: CozeLoop's prompt development module provides developers with full-process support from writing, debugging, and optimization to version management. With a visual Playground, developers can conduct real-time interactive testing of prompts, enabling intuitive comparisons of the output effects of different LLMs.
* **Evaluation**: CozeLoop's evaluation module provides developers with systematic evaluation capabilities, enabling multi-dimensional automated testing of the output effects of prompts and Coze agents, such as accuracy, conciseness, and compliance.
* **Observation**: CozeLoop provides developers with visual observation capabilities for the full-chain execution process, fully recording each processing step from user input to AI output. This includes key nodes such as prompt parsing, model invocation, and tool execution, while automatically capturing intermediate results and abnormal states.

## Feature list
| **Feature** | **Feature points** | **Commercial version** | **Open-source Edition** |
| --- | --- | --- | --- |
| Prompt debugging | Playground debugging, comparison, and version management | ✔️ | ✔️ |
|  | Prompt optimization | ✔️ | - |
| Evaluation | Evaluation set | ✔️ | ✔️ |
|  | Evaluator | ✔️ | ✔️ |
|  | Experiment | ✔️ | ✔️ |
| Observation | Trace | ✔️ | ✔️ |
|  | Metric | ✔️ | - |
|  | Automated tasks | ✔️ | - |
| Model | Model management | ✔️ | ✔️ |
| Security | SSO login | ✔️ | - |
|  | Data security (VPC private network link, content security policy) | ✔️ | - |
| Team and enterprise management | Workspace and member management | ✔️ | - |
|  | Collaboration | ✔️ | - |
|  | Enterprise teams and permissions | ✔️ | - |
## Quickstart
Refer to the [Quickstart](https://github.com/coze-dev/cozeloop/wiki/2.-Quickstart) to learn how to install and deploy the latest version of CozeLoop.
## Using CozeLoop Open-source Edition

* [Prompt development and debugging](https://loop.coze.cn/open/docs/cozeloop/create-prompt): CozeLoop provides a complete prompt development workflow.
* [Evaluation](https://loop.coze.cn/open/docs/cozeloop/create-prompt): CozeLoop's evaluation functionality offers standardized evaluation data management, automated assessment engines, and comprehensive experimental result statistics.
* [Trace reporting and querying](https://loop.coze.cn/open/docs/cozeloop/trace-integrate): CozeLoop supports automatic Trace reporting for prompt debugging conducted on the platform, enabling real-time tracking of each Trace data.
* [Open-source Edition usage of the CozeLoop SDK](https://github.com/coze-dev/cozeloop/wiki/8.-Open-source-edition-uses-CozeLoop-SDK): The CozeLoop SDK in three languages is suitable for both commercial and open-source editions. For the Open-source Edition, developers only need to modify some parameter configurations during initialization.

## Development guide

* [System architecture](https://github.com/coze-dev/cozeloop/wiki/3.-Architecture): Learn about the technical architecture and core components of CozeLoop Open-source Edition.
* [Startup mode](https://github.com/coze-dev/cozeloop/wiki/4.-Service-startup-modes): When installing and deploying CozeLoop Open-source Edition, the default development mode allows backend file modifications without requiring service redeployment.
* [Model configuration](https://github.com/coze-dev/cozeloop/wiki/5.-Model-configuration): CozeLoop Open-source Edition supports various LLM models through the Eino framework. Refer to this document to view the supported model list and learn how to configure models.
* [Code development and testing](https://github.com/coze-dev/cozeloop/wiki/6.-Code-development-and-testing): Learn how to perform secondary development and testing based on CozeLoop Open-source Edition.
* [Fault troubleshooting](https://github.com/coze-dev/cozeloop/wiki/7.-Troubleshooting): Learn how to check container status and system logs.

## License
This project uses the Apache 2.0 license. For more details, please refer to the [LICENSE](LICENSE) file.
## Community Contributions
We welcome community contributions. For contribution guidelines, please refer to [CONTRIBUTING](CONTRIBUTING.md) and [Code of conduct](CODE_OF_CONDUCT.md). We look forward to your contributions!
## Security and Privacy
If you identify potential security issues in this project or believe you may have found one, please notify Bytedance's security team via our [Security Center](https://security.bytedance.com/src) or [Vulnerability Report Email](sec@bytedance.com).
Please **do not** create public GitHub Issues.
## Join the Community
Scan the QR code below on the Lark mobile app to join the CozeLoop technical discussion group

![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/8fae8f0e7b124831b1dd94aa9a5c60c1~tplv-goo7wpa0wc-image.image)
## Acknowledgments
Thanks to all developers and community members who contributed to the CozeLoop project Special thanks:

* LLM integration support provided by the Eino framework team
* High-performance frameworks developed by the Cloudwego team
* All users who participated in testing and feedback
