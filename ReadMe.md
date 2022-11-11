# gitops-release-manager
# helping links:
https://learn.microsoft.com/en-us/azure/devops/boards/work-items/work-item-url-hyperlink?view=azure-devops
https://stackoverflow.com/questions/18679870/list-commits-between-2-commit-hashes-in-git
https://dev.azure.com/adamrussak/go-release-mamnger-test/_workitems/edit/1
example of work-item url: `https://dev.azure.com/OrganizationName/ProjectName/_workitems/edit/WorkItemNumber`
try use `go-git`
```sh
# tag old commit
git tag -a 1.0.1 e351d1c704ca5a7754943a21a2d58ec75608e7f7

#get list of commits betwen 2 commits 
git rev-list --ancestry-path 7b4a07a..ecf5891
```
## curent flow:
- each build bump version
    - each service has its own version (micro-service)
- product version is set

## requested:
- version will be in tags
- getting Work-items in release

### for commits 
-1 
-1 
-1 
- [ ] test