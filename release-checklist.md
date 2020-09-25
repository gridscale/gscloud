# Release Checklist

We use GoReleaser to create release builds and GH release drafts.

On master branch:

- [ ] Test everything is working as expected
- [ ] Make sure README is up to date
- [ ] Changelog: add additions and fixes
- [ ] Changelog: set release date
- [ ] Create a new tag (e.g. `git tag v0.5.0`)
- [ ] Push the tag (e.g. `git push origin --tags v0.5.0`)

That's it. Everything else should be done by the pipeline. Check [github.com/gridscale/gscloud/releases](https://github.com/gridscale/gscloud/releases).

- [ ] Finally, create new changelog stub for next release
