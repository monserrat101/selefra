<!-- Your Title -->
<p align="left">
<img src="https://user-images.githubusercontent.com/124020340/224677116-44ae9c6c-a543-4813-9ef3-c7cbcacd2fbe.png" width="400">
</p>

<!-- Badges -->
<p align="left">   
<a href="https://github.com/selefra/selefra/blob/master/LICENSE"><img alt="GitHub license" src="https://img.shields.io/badge/Slack-666?logo=slack"></a>
<a href="https://github.com/selefra/selefra/blob/master/LICENSE"><img alt="GitHub license" src="https://img.shields.io/github/downloads/selefra/selefra/total?logo=github"></a>
<a href="https://github.com/selefra/selefra/blob/master/LICENSE"><img alt="Twitter" src="https://img.shields.io/badge/Twitter-666?logo=Twitter"></a>
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

See [Docs](https://selefra.io/docs) for best practices and detailed instructions. In docs, you will find info on installation, CLI usage, project workflow and more guides on how to accomplish cloud inspection tasks.

<img align="right" width="400" src="https://user-images.githubusercontent.com/124020340/224889579-556ee877-28e0-4638-b88f-ee9a4564c33a.png" />

#### 🔥 Olicy As Code

Custom analysis policies (security, compliance, cost) can be written through a combination of SQL and YAML.

#### 💥 Configuration of Multi-Cloud, Multi-SaaS

Unified multi-cloud configuration data integration capabilities that can support analysis of configuration data from any cloud service via SQL.

#### 🌟 Version Control

Analysis policies can be managed through VCS such as GitHub/Gitlab.

#### 🥤 Automation

Policies can be automated to enforce compliance, security, and cost optimization rules through Scheduled tasks and cloud automation tools.

## Welcome

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
    
## Feature

|    | Language | Status |
| -- | -------- | ------ |
| <img src="https://user-images.githubusercontent.com/124020340/224913715-846ee552-1ecf-4ad2-ae62-b59f35a44a75.png" height=38 />     | [AWS](https://www.selefra.io/docs/providers-connector/aws) | Stable |
| <img src="https://user-images.githubusercontent.com/124020340/224914312-4889ecc5-7389-46c6-b702-5d23e3e1be16.png" height=38 />     | [GCP](https://www.selefra.io/docs/providers-connector/gcp) | Stable |
| <img src="https://user-images.githubusercontent.com/124020340/224914454-dac803a6-7f1e-4b98-869a-7b72e329f312.png" height=38 />     | [K8S](https://www.selefra.io/docs/providers-connector/k8s) | Stable |
| <img src="https://user-images.githubusercontent.com/124020340/224914705-ee2f1d63-c4e2-4bce-aea3-72851d65c135.png" height=38 />      | [Microsoft365](https://www.selefra.io/docs/providers-connector/microsoft365)     | Stable |
| <img src="https://user-images.githubusercontent.com/124020340/224914806-8d6d9f91-e332-47b9-9003-f877081383c0.png" height=38 />      | [Slack](https://www.selefra.io/docs/providers-connector/slack)     | Stable |
| <img src="https://user-images.githubusercontent.com/124020340/224914970-404a97c9-40eb-432a-b01f-d54f11fdc4c3.png" height=38 />      | [Snowflake](https://www.selefra.io/docs/providers-connector/snowflake)     | Stable |

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
