<!-- Your Title -->
<p align="center">
<a href="https://www.selefra.io/" target="_blank">
<picture><source media="(prefers-color-scheme: dark)" srcset="https://user-images.githubusercontent.com/124020340/225567784-61adb5e7-06ae-402a-9907-69c1e6f1aa9e.png"><source media="(prefers-color-scheme: light)" srcset="https://user-images.githubusercontent.com/124020340/224677116-44ae9c6c-a543-4813-9ef3-c7cbcacd2fbe.png"><img width="400px" alt="Steampipe Logo" src="https://user-images.githubusercontent.com/124020340/224677116-44ae9c6c-a543-4813-9ef3-c7cbcacd2fbe.png"></picture>
<a/>
</p>

<!-- Description -->
  <p align="center">
    <i>Selefra is an open-source policy-as-code software that provides analytics for multi-cloud and SaaS.</i>
  </p>
  
  <!-- Badges -->
<p align="center">   
<img alt="go" src="https://img.shields.io/badge/go-1.19-1E90FF"></a>
<a href="https://github.com/selefra/selefra"><img alt="Total" src="https://img.shields.io/github/downloads/selefra/selefra/total?logo=github"></a>
<a href="https://github.com/selefra/selefra/blob/master/LICENSE"><img alt="GitHub license" src="https://img.shields.io/github/license/selefra/selefra?style=social"></a>
  </p>
  
  <!-- Badges -->
  <p align="center">
<a href="https://selefra.io/community/join"><img src="https://img.shields.io/badge/-Slack-424549?style=social&logo=Slack" height=25></a>
    &nbsp;
    <a href="https://twitter.com/SelefraCorp"><img src="https://img.shields.io/badge/-Twitter-red?style=social&logo=twitter" height=25></a>
    &nbsp;
    <a href="https://www.reddit.com/r/Selefra"><img src="https://img.shields.io/badge/-Reddit-red?style=social&logo=reddit" height=25></a>
    &nbsp;
    <a href="https://selefra.medium.com/"><img src="https://img.shields.io/badge/-Medium-red?style=social&logo=medium" height=25></a>

  </p>
  
<p align="center">
  <img src="https://user-images.githubusercontent.com/124020340/225897757-188f1a50-2efa-4a9e-9199-7cb7f68485be.png">
</p>
<br/>

<!-- About Selefra -->

## About Selefra

Selefra is an open-source data integration and analysis tool for developers. You can use Selefra to extract, load, and analyze infrastructure data anywhere from Public Cloud, SaaS platform, development platform, and more.

See [Docs](https://selefra.io/docs/introduction) for best practices and detailed instructions. In docs, you will find info on installation, CLI usage, project workflow and more guides on how to accomplish cloud inspection tasks.


<img align="right" width="520" alt="img_code" src="https://user-images.githubusercontent.com/124020340/225893407-40e5df4a-7ddd-439c-b5fc-f6d1dfac68bf.png">

#### 🔥 Policy As Code

Custom analysis policies (security, compliance, cost) can be written through a combination of SQL and YAML.

#### 💥 Configuration of Multi-Cloud, Multi-SaaS

Unified multi-cloud configuration data integration capabilities that can support analysis of configuration data from any cloud service via SQL.

#### 🌟 Version Control

Analysis policies can be managed through VCS such as GitHub/Gitlab.

#### 🥤 Automation

Policies can be automated to enforce compliance, security, and cost optimization rules through Scheduled tasks and cloud automation tools.

## Getting started

Read detailed documentation for how to [Get Started](https://selefra.io/docs/get-started/) with Selefra.

For quick start, run this demo, it should take less than a few minutes:

1. **Install Selefra**

    For non-macOS users, [download packages](https://github.com/selefra/selefra/releases) to install Selefra.

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
    mkdir selefra-demo & cd selefra-demo & selefra init
    ```

3. **Build code**

    ```bash
    selefra apply 
    ```
    
## Selefra Community Ecosystem









 Provider | Introduce | Status |
 | --------| ----- | ------ |
 | [AWS](https://www.selefra.io/docs/providers-connector/aws)|The AWS Provider for Selefra can be used to extract data from many of the cloud services by AWS. The provider must be configured with credentials to extract and analyze infrastructure data from AWS. | Stable |
| [GCP](https://www.selefra.io/docs/providers-connector/gcp)|The GCP Provider for Selefra can be used to extract data from many of the cloud services by GCP. The provider must be configured with credentials to extract and analyze infrastructure data from GCP. | Stable |
| [K8S](https://www.selefra.io/docs/providers-connector/k8s)|The K8s Provider for Selefra can be used to extract data from many of the cloud services by K8s. The provider must be configured with credentials to extract and analyze infrastructure data from K8s. | Stable |
| [Microsoft365](https://www.selefra.io/docs/providers-connector/microsoft365)| The Microsoft365 Provider for Selefra can be used to extract data from many of the cloud services by Microsoft365. The provider must be configured with credentials to extract and analyze infrastructure data from Microsoft365.    | Stable |
| [Slack](https://www.selefra.io/docs/providers-connector/slack)|  The Slack Provider for Selefra can be used to extract data from many of the cloud services by Slack. The provider must be configured with credentials to extract and analyze infrastructure data from Slack.   | Stable |
| [Snowflake](https://www.selefra.io/docs/providers-connector/snowflake)|  The snowflake Provider for Selefra can be used to extract data from many of the cloud services by snowflake. The provider must be configured with credentials to extract and analyze infrastructure data from snowflake.   | Stable |

## Community

Selefra is a community-driven project, we welcome you to open a [GitHub Issue](https://github.com/selefra/selefra/issues/new/choose) to report a bug, suggest an improvement, or request new feature.

-  Join <a href="https://selefra.io/community/join"><img height="16" alt="humanitarian" src="https://user-images.githubusercontent.com/124020340/225563969-3f3d4c45-fb3f-4932-831d-01ab9e59c921.png"></a> [Selefra Community](https://selefra.io/community/join) on Slack. We host `Community Hour` for tutorials and Q&As on regular basis.
-  Follow us on <a href="https://twitter.com/SelefraCorp"><img height="16" alt="humanitarian" src="https://user-images.githubusercontent.com/124020340/225564426-82f5afbc-5638-4123-871d-fec6fdc6457f.png"></a> [Twitter](https://twitter.com/SelefraCorp) and share your thoughts！
-  Email us at <a href="support@selefra.io"><img height="16" alt="humanitarian" src="https://user-images.githubusercontent.com/124020340/225564710-741dc841-572f-4cde-853c-5ebaaf4d3d3c.png"></a>&nbsp;support@selefra.io

## Contributing

For developers interested in building Selefra codebase, read through [Contributing.md](https://github.com/selefra/selefra/blob/main/CONTRIBUTING.md) and [Selefra Roadmap](https://github.com/orgs/selefra/projects/1).
Let us know what you would like to work on!

## License

[Mozilla Public License v2.0](https://github.com/selefra/selefra/blob/main/LICENSE)
