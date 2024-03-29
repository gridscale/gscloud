# Release Checklist

We use GoReleaser to create release builds and GH release drafts.

On master branch:

- [ ] Test everything is working as expected
- [ ] Make sure README is up to date
- [ ] Changelog: add additions and fixes
- [ ] Changelog: set release date
- [ ] Create a new annotated tag (e.g. `git tag -a -m 'New release' v0.12.0`)
- [ ] Push the tag (e.g. `git push origin --follow-tags`)

That's it. Everything else should be done by the pipeline.

- [ ] Check [github.com/gridscale/gscloud/releases](https://github.com/gridscale/gscloud/releases) for a new draft release. Copy changelog entries there and submit the release
- [ ] Finally, create new changelog stub for next release
