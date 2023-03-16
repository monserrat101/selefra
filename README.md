<!-- Your Title -->
<p align="left">
<img src="https://user-images.githubusercontent.com/124020340/224677116-44ae9c6c-a543-4813-9ef3-c7cbcacd2fbe.png" width="400">
</p>

<!-- Badges -->
<p align="left">   
<a href="https://www.selefra.io/community/join"><img alt="Slack" src="https://img.shields.io/badge/Slack-666?logo=slack"></a>
<a href="https://github.com/selefra/selefra"><img alt="Total" src="https://img.shields.io/github/downloads/selefra/selefra/total?logo=github"></a>
<a href="https://twitter.com/SelefraCorp"><img alt="Twitter" src="https://img.shields.io/badge/Twitter-666?logo=Twitter"></a>
<a href="https://github.com/selefra/selefra/blob/master/LICENSE"><img alt="GitHub license" src="https://img.shields.io/github/license/selefra/selefra?style=social"></a>
  </p>

<!-- Description -->
  <p align="left">
    <i>Selefra is an open-source policy-as-code software that provides analytics for multi-cloud and SaaS.</i>
  </p>

<br/>

<!-- About Selefra -->

## About Selefra

Selefra is an open-source data integration and analysis tool for developers. You can use Selefra to extract, load, and analyze infrastructure data anywhere from Public Cloud, SaaS platform, development platform, and more.

See [Docs](https://selefra.io/docs/introduction) for best practices and detailed instructions. In docs, you will find info on installation, CLI usage, project workflow and more guides on how to accomplish cloud inspection tasks.

<img align="right" width="400" src="https://user-images.githubusercontent.com/124020340/224889579-556ee877-28e0-4638-b88f-ee9a4564c33a.png" />

#### 🔥 Policy As Code

Custom analysis policies (security, compliance, cost) can be written through a combination of SQL and YAML.

#### 💥 Configuration of Multi-Cloud, Multi-SaaS

Unified multi-cloud configuration data integration capabilities that can support analysis of configuration data from any cloud service via SQL.

#### 🌟 Version Control

Analysis policies can be managed through VCS such as GitHub/Gitlab.

#### 🥤 Automation

Policies can be automated to enforce compliance, security, and cost optimization rules through Scheduled tasks and cloud automation tools.

## Getting started

Read detailed documentation for how to [get started](https://selefra.io/docs/get-started/) with Selefra.

For quick start, run this demo, it should take less than a few miniutes:

1. **Install Selefra**

    For non macOS users, [download packages](https://github.com/selefra/selefra/releases) to install Selefra.

    On macOS, tap Selefra with Homebrew:

    ```bash
    brew tap selefra/tap
    ```

    Next, install Selefra:

    ```bash
    brew install selefra/tap/selefra
    ```

2. **Initialization project**

    ```bash
    selefra init selefra-demo && cd selefra-demo
    ```

3. **Build code**

    ```bash
    selefra apply 
    ```
    
## Selefra Community Ecosystem









|    | Language | Status |
| -- | -------- | ------ |
| <img width="38" alt="aws logo" src="https://user-images.githubusercontent.com/124020340/225558573-35579326-0fc8-4100-8c30-7aad82788d61.png">     | [AWS](https://www.selefra.io/docs/providers-connector/aws) | Stable |
| <img width="38" alt="Google logo" src="https://user-images.githubusercontent.com/124020340/225558584-6309e72b-b92c-405c-90dd-64516f6965ef.png">    | [GCP](https://www.selefra.io/docs/providers-connector/gcp) | Stable |
| <img width="38" alt="k8s" src="https://user-images.githubusercontent.com/124020340/225558598-09e03a70-b4ea-47ec-890c-d110d2eb5b5d.png">    | [K8S](https://www.selefra.io/docs/providers-connector/k8s) | Stable |
| <img width="38" alt="Microsoft" src="https://user-images.githubusercontent.com/124020340/225558609-4aac1a66-92b7-4c9b-9ccb-75948f86b61c.png">      | [Microsoft365](https://www.selefra.io/docs/providers-connector/microsoft365)     | Stable |
| <img width="38" alt="slack" src="https://user-images.githubusercontent.com/124020340/225558623-50850a40-7505-44dc-b255-a2574ae4216f.png">     | [Slack](https://www.selefra.io/docs/providers-connector/slack)     | Stable |
| <img width="38" alt="snowflake" src="https://user-images.githubusercontent.com/124020340/225558631-c7b26728-bc7b-495a-8b48-efd846e703c8.png"> | [Snowflake](https://www.selefra.io/docs/providers-connector/snowflake)     | Stable |

## Community

Selefra is a community-driven project, we welcome you to open a [GitHub Issue](https://github.com/selefra/selefra/issues/new/choose) to report a bug, suggest an improvement, or request new feature.

-  Join [Selefra Community](https://selefra.io/community/join) on Slack. We host `Community Hour` for tutorials and Q&As on regular basis.
-  Follow us on [Twitter](https://twitter.com/SelefraCorp) and share your thoughts！
-  Email us at support@selefra.io

## Contributing

For developers interested in building Selefra codebase, read through [Contributing.md](https://github.com/selefra/selefra/blob/main/CONTRIBUTING.md) and [Selefra Roadmap](https://github.com/orgs/selefra/projects/1).
Let us know what you would like to work on!

## License

[Mozilla Public License v2.0](https://github.com/selefra/selefra/blob/main/LICENSE)
