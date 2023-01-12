# gitops-release-manager
## Changes
- updated ReadMe.md
- added paramter for file output name
    - support for new tag as file name
    - custome file name
- validate outputa file has a suffix of `.md `
## Breaking changes
- output name is the new Tag by default and not Report.md
    - Old default: `Report.md`<br> New default: `1.0.0.md`