# Starchive

<!-- About The Project -->
## About The Project

There are lots of awesome projects out there, and its common to star your
favorite projects/repos.  Sometimes those projects become archived, or disappear
from the internet for various reasons.  In order to archive most of these
projects and ensure you still have a copy of them, this project aims to enable
a quick binary to use with cron, and enable scheduled backups of all those
repos.


<!-- Getting Started -->
## Getting Started




<!-- Prerequisites -->
### Prerequisites

Prerequisities of installing this application include the potential compilation,
and retrieval of a github user access token.

#### Compiling

Some pre-requisites for compiling this binary are the following packages:
`golang`, and `make`.

##### macOS

```bash
brew install golang make
```

##### debian

```bash
sudo apt install golang make
```

##### Clone the repository
    
```bash
git clone git@github.com:ChrisSwanson/Starchive.git && cd Starchive
```

##### Compile the binary

```bash
make build
```

Once the binary is done being compiled, the build can be found in build
directory.

#### Create Github Personal Access Token

Instructions can be found [here](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token) on how to create a github personal access
token.
This application will just require basic access to `public_repo` (browse public
repositories).


#### Configuration

It is recommended that this application be setup with cron to enable regularly
scheduled backups of github starred projects.

This was built with simplicity in mind, all flags have a 1:1 relation that can
be used in a configuration file.

##### Flags

`--dir $TARGET_DIRECTORY`
 - the target directory to populate with git clones and pulls.

`--token $GITHUB_PERSONAL_ACCESS_TOKEN`
 - The token flag is used to provide the github personal access token.  This is
 used to provide context to which it will follow the starred repo

`--debug`
 - This flag is utilized to enable debug logging.

##### config.yml
```yaml
dir: /tmp/starchive
debug: true
token: mypersonalgithubaccesstoken1234567890abc
```

<!-- Installation -->
### Installation

#### Copy the binary

##### To make this binary more accessible, copy to your local bin directory

```bash
cp build/starchive /usr/local/bin/
```

##### Create a cronjob, with the path for your config in line:

I.E. if the config.yml file is provided in /etc/starchive/config.yml,
amd we want to run hourly at the top of the hour, we could enter something like
this:

```bash
crontab -l | { cat; echo "0 * * * * cd /etc/starchive && /usr/local/bin/starchive"; } | crontab -
```


<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.


## FAQ

Why does this project use a github personal access token, when this can just screen scrape?

[Github rate limits](https://docs.github.com/en/rest/overview/resources-in-the-rest-api#rate-limiting)
can be more restrictive for unauthenticated API based requests, and APIs are
going to be more consistent than screen scraping and breaking in future builds. 
Additionally the github personal access token is used to determine your
username.
In the link provided for Github Rate Limits (at the time of writing this),
github allows for the following rate limits:

| Auth | Rate Limit |
| :--- | :---: |
| Unauthenticated | 60/hour |
| Personal Access Token | 1000/hour |
| Basic / OAuth | 5000/hour |
