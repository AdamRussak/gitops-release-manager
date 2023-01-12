# gitops-release-manager
## links used in the creation:
- [Azure DevOps boarrds work-items API](https://learn.microsoft.com/en-us/azure/devops/boards/work-items/work-item-url-hyperlink?view=azure-devops)
- [Stackoverflow git compare commits](https://stackoverflow.com/questions/18679870/list-commits-between-2-commit-hashes-in-git)
- example of work-item url structure: `https://dev.azure.com/OrganizationName/ProjectName/_workitems/edit/WorkItemNumber`
- go module used to manage git: `go-git` <br>

**the git command that is used to get the commits:**
```sh
# tag old commit
git tag -a 1.0.1 e351d1c704ca5a7754943a21a2d58ec75608e7f7

#get list of commits betwen 2 commits 
git rev-list --ancestry-path 7b4a07a..ecf5891
```
## curent flow:
1. provide the command with the optional flags: `gitops-version release --repo-path /home/adam/<repo> --org <AzureDevOps orgID> --project <AzureDevOps project> --pat <AzureDevOps PAT > --hash <Git commit Hash> <path to git ssh keys> -v`
1. The code looks for latest tag and gets all the commits between it and the provided commit Hash
1. if dry run was not provided, the output will be in the `CWD` with the name: `Report.md`
### supported paramters:
global paramters:
```sh
-v, --verbose   verbose logging
```
`gitops-release-manager -h`:
```sh
    --version   print out the current version
```
`gitops-release-manager release -h`:
```sh
      --auth        string  Set Auth type (ssh or https (default "https")
      --dry-run     bool    If true, only run a dry-run with cli output
      --filename    string  Costume file name
      --git-branch  string  Set Brnach to tag (default "main")
      --git-email   string  Set email to tag with (default ".")
      --git-keyPath string  Set email to tag with (default "~/.ssh/id_rsa")
      --git-push    bool    If true, only run a dry-run with cli output
      --git-user    string  Set userName to tag with (default ".")
      --hash        string  Set new TAG Hash
      --org         string  Set Azure DevOps orgenziation
      --output      string  Set path to report output (default "./")
      --pat         string  Set PAT for API calls
      --project     string  Set Azure DevOps project
      --repo-path   string  Set Path to Git repo root (default ".")
```