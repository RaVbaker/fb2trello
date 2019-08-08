# fb2trello

Small tool to setup, plan and track social media posts.

Can setup the whole Trello board for planning social media posts and import existing posts from Facebook page. 

## Arguments:

```
Usage of fb2trello:
  -board string
        Trello board name to which posts should be archived
  -lists Calendar,Ideas,Planned,Published
        Trello list names, last after comma is the one for archive, default/e.g. Calendar,Ideas,Planned,Published (default "Calendar,Ideas,Planned,Published")
  -page string
        Facebook pageName/ID  which should get archived
  -setup
        If specified it will create whole structure in trello for board and lists
  -until string
        Archive until date - oldest post publication date, e.g. 2019-07-30
``` 

Additionally if `-page` argument is provided API  `FACEBOOK_ACCESS_TOKEN` is needed.

If the `-board` is provided it requires `TRELLO_API_KEY` and `TRELLO_TOKEN` for accessing Trello API.

The `-setup` is optional and needed only on first setup of Trello. On second run the `-lists` can have only single last node `Published`. 

## Contribution

Please make a PR/issue if anything needed.

--- 
(c) Rafal "RaVbaker" Piekarski 2019 

License? as specified in LICENSE file.
