# Release Checklist

This repo uses git flow, so better go and install those tools first.

On develop branch:

- [ ] `git flow release start 0.3.0-beta`

This will create and checkout a new release branch. On that release branch:

- [ ] test everything is working as expected
- [ ] make sure README is up to date
- [ ] update changelog: add additions and fixes
- [ ] update changelog: set release date
- [ ] `git flow release publish 0.3.0-beta`
- [ ] `git flow release finish 0.3.0-beta`
- [ ] create and push the tag `git push upstream --tags v0.3.0-beta` (Note: its upstream here but it needs to be your remote's name here)

Make sure the new release branch is pushed to the right remote. Then go to GitHub and

- [ ] create two (yes that's two) PRs on GitHub: one from release branch → develop
- [ ] another from release branch → master (make sure you do not accidentally remove the release branch when merging those PRs)

Finally, go to GitHub again and

- [ ] make a new GitHub release from the tag, copy the changelog text into that release
- [ ] do `make buildall` and drop the zip files into the new GitHub release
- [ ] have a beer, you just survived the most complicated way to release software

Back on develop branch:

- [ ] create new CHANGELOG.md stub for next release
- [ ] bump version in `VERSION` file
